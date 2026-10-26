package app

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/HazelnutParadise/insyra/engine/ring"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"google.golang.org/genai"
)

const ollamaServerURL = "http://ollama:11434"

var (
	llmLocalRequestInterval  time.Duration = 0 * time.Millisecond
	llmGeminiRequestInterval time.Duration = 2100 * time.Millisecond
)

type llmResponse struct {
	text string
	err  error
}

type llmReq struct {
	prompt      string
	temperature float64
	model       int
	reply       chan llmResponse
}

// llmQueue keeps each model FIFO and alternates models when both have work.
// All ring operations and the next-model cursor share the same lock.
type llmQueue struct {
	mu     sync.Mutex
	ready  *sync.Cond
	queues [2]*ring.Ring[llmReq]
	next   int
}

func newLLMQueue() *llmQueue {
	q := &llmQueue{queues: [2]*ring.Ring[llmReq]{
		ring.NewRing[llmReq](8),
		ring.NewRing[llmReq](8),
	}}
	q.ready = sync.NewCond(&q.mu)
	return q
}

func (q *llmQueue) push(req llmReq) {
	q.mu.Lock()
	q.queues[req.model].Push(req)
	q.ready.Signal()
	q.mu.Unlock()
}

func (q *llmQueue) pop() llmReq {
	q.mu.Lock()
	defer q.mu.Unlock()
	for {
		for offset := range len(q.queues) {
			model := (q.next + offset) % len(q.queues)
			queue := q.queues[model]
			if req, ok := queue.PopFront(); ok {
				q.next = (model + 1) % len(q.queues)
				if queue.Len() == 0 {
					// Ring.PopFront retains its backing slot; release drained prompts.
					queue.Clear()
				}
				return req
			}
		}
		q.ready.Wait()
	}
}

var llmRequests = newLLMQueue()

const (
	LLM_Gemini = iota
	LLM_Local
)

func init() {
	go func() {
		for {
			req := llmRequests.pop()
			text, err := requestLLM(req)
			req.reply <- llmResponse{text: text, err: err}
		}
	}()
}

func CallLLM(prompt string, temperature float64, model int) (string, error) {
	if model < 0 || model >= len(llmRequests.queues) {
		return "", fmt.Errorf("unsupported LLM model: %d", model)
	}
	reply := make(chan llmResponse, 1)
	llmRequests.push(llmReq{
		prompt:      prompt,
		temperature: temperature,
		model:       model,
		reply:       reply,
	})
	result := <-reply
	return result.text, result.err
}

func Unmarshal_LLM_JSON_Response[T any](resp string, v *T) error {
	resp = strings.Replace(resp, "```json", "", -1)
	resp = strings.Replace(resp, "```", "", -1)
	return json.Unmarshal([]byte(resp), v)
}

func requestLLM(req llmReq) (string, error) {
	var completion string
	switch req.model {
	case LLM_Gemini:
		time.Sleep(llmGeminiRequestInterval / 2)
		ctx := context.Background()
		key := os.Getenv("GEMINI_API_KEY")
		if key == "" {
			return "", fmt.Errorf("GEMINI_API_KEY is required")
		}
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  key,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			return "", err
		}
		response, err := client.Models.GenerateContent(ctx, "gemini-3.5-flash-lite", genai.Text(req.prompt), &genai.GenerateContentConfig{
			Temperature:           genai.Ptr(float32(req.temperature)),
			AutomaticContinuation: genai.Ptr(false),
		})
		if err != nil {
			return "", err
		}
		if response == nil {
			return "", fmt.Errorf("Gemini returned an empty response")
		}
		if feedback := response.PromptFeedback; feedback != nil && feedback.BlockReason != "" && feedback.BlockReason != genai.BlockedReasonUnspecified {
			return "", fmt.Errorf("Gemini blocked the prompt: %s", feedback.BlockReason)
		}
		for _, candidate := range response.Candidates {
			if candidate != nil && (candidate.FinishReason == genai.FinishReasonSafety || candidate.FinishReason == genai.FinishReasonRecitation) {
				return "", fmt.Errorf("Gemini blocked the response: %s", candidate.FinishReason)
			}
		}
		completion = response.Text()
		if completion == "" {
			return "", fmt.Errorf("Gemini returned no text")
		}
		time.Sleep(llmGeminiRequestInterval / 2)
	case LLM_Local:
		time.Sleep(llmLocalRequestInterval / 2)
		ctx := context.Background()
		llm, err := ollama.New(ollama.WithModel("gemma3:12b"), ollama.WithRunnerNumCtx(10240), ollama.WithServerURL(ollamaServerURL))
		if err != nil {
			return "", err
		}
		completion, err = llms.GenerateFromSinglePrompt(ctx, llm, req.prompt)
		if err != nil {
			return "", err
		}
		time.Sleep(llmLocalRequestInterval / 2)
	}

	return completion, nil
}
