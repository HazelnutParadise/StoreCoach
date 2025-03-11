package app

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/HazelnutParadise/insyra"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/ollama"
)

const ollamaServerURL = "http://ollama:11434"

var (
	llmLocalRequestInterval  time.Duration = 0 * time.Millisecond
	llmGeminiRequestInterval time.Duration = 2100 * time.Millisecond
)

var (
	llmGeminiReqBuf          = insyra.NewDataList()
	llmLocalGemma2_9b_ReqBuf = insyra.NewDataList()
)
var llmReqBufs = []*insyra.DataList{llmGeminiReqBuf, llmLocalGemma2_9b_ReqBuf}
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
	LLM_LocalGemma2_9b
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
	}
	switch model {
	case LLM_Gemini:
		llmGeminiReqBuf.InsertAt(0, req)
	case LLM_LocalGemma2_9b:
		llmLocalGemma2_9b_ReqBuf.InsertAt(0, req)
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
		llm, err := googleai.New(ctx, googleai.WithAPIKey(
			os.Getenv("GEMINI_API_KEY")),
			googleai.WithDefaultModel("gemini-2.0-flash-lite"),
			googleai.WithDefaultTemperature(req.temperature),
		)
		completion, err = llms.GenerateFromSinglePrompt(ctx, llm, req.prompt)
		if err != nil {
			return "", err
		}
		time.Sleep(llmGeminiRequestInterval / 2)
	case LLM_LocalGemma2_9b:
		time.Sleep(llmLocalRequestInterval / 2)
		ctx := context.Background()
		llm, err := ollama.New(ollama.WithModel("gemma2:9b"), ollama.WithRunnerNumCtx(10240), ollama.WithServerURL(ollamaServerURL))
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
