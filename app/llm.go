package app

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/HazelnutParadise/insyra"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"google.golang.org/genai"
)

const ollamaServerURL = "http://ollama:11434"

var (
	llmLocalRequestInterval  time.Duration = 0 * time.Millisecond
	llmGeminiRequestInterval time.Duration = 2100 * time.Millisecond
)

var (
	llmGeminiReqBuf = insyra.NewDataList()
	llmLocalReqBuf  = insyra.NewDataList()
)
var llmReqBufs = []*insyra.DataList{llmGeminiReqBuf, llmLocalReqBuf}
var llmRespBuf = sync.Map{}
var llmErrBuf = sync.Map{}

type llmReq struct {
	uuid        string
	prompt      string
	temperature float64
	model       int
}

const (
	LLM_Gemini = iota
	LLM_Local
)

func init() {
	go func() {
		for {
			for _, llmReqBuf := range llmReqBufs {
				var reqUUID string
				var req llmReq
				if llmReqBuf.Len() > 0 {
					req = llmReqBuf.Pop().(llmReq)
					reqUUID = req.uuid
				}
				if reqUUID == "" {
					continue
				}
				resp, err := requestLLM(req)
				if err != nil {
					llmErrBuf.Store(reqUUID, err)
				} else {
					llmRespBuf.Store(reqUUID, resp)
				}
			}
		}
	}()
}

func CallLLM(prompt string, temperature float64, model int) (string, error) {
	reqUUID := uuid.New().String()
	req := llmReq{
		uuid:        reqUUID,
		prompt:      prompt,
		temperature: temperature,
		model:       model,
	}
	switch model {
	case LLM_Gemini:
		llmGeminiReqBuf.InsertAt(0, req)
	case LLM_Local:
		llmLocalReqBuf.InsertAt(0, req)
	}
	for {
		if llmResp, ok := llmRespBuf.LoadAndDelete(reqUUID); ok {
			return llmResp.(string), nil
		} else if llmErr, ok := llmErrBuf.LoadAndDelete(reqUUID); ok {
			return "", llmErr.(error)
		}
	}
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
		response, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash-lite", genai.Text(req.prompt), &genai.GenerateContentConfig{
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
