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
)

var (
	llmRequestInterval time.Duration = 500 * time.Millisecond
)

var llmReqBuf = insyra.NewDataList()
var llmRespBuf = sync.Map{}
var llmErrBuf = sync.Map{}

type llmReq struct {
	uuid        string
	prompt      string
	temperature float64
}

func init() {
	go func() {
		for {
			time.Sleep(llmRequestInterval)
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
	}()
}

func CallLLM(prompt string, temperature float64) (string, error) {
	reqUUID := uuid.New().String()
	req := llmReq{
		uuid:        reqUUID,
		prompt:      prompt,
		temperature: temperature,
	}
	llmReqBuf.InsertAt(0, req)
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
	ctx := context.Background()
	llm, err := googleai.New(ctx, googleai.WithAPIKey(
		os.Getenv("GEMINI_API_KEY")),
		googleai.WithDefaultModel("gemini-2.0-flash"),
		googleai.WithDefaultTemperature(req.temperature),
	)
	if err != nil {
		return "", err
	}

	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, req.prompt)
	if err != nil {
		return "", err
	}
	return completion, nil
}
