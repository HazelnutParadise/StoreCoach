package app

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

var (
	llmRequestLimitPer10Seconds int32 = 10
)

var llmReqBuf = sync.Map{}
var llmRespBuf = sync.Map{}
var llmErrBuf = sync.Map{}

type llmReq struct {
	prompt      string
	temperature float64
}

func init() {
	go func() {
		var reqLimitLeft int32 = llmRequestLimitPer10Seconds
		var lastReqTime int64
		for {
			if reqLimitLeft == 0 {
				if time.Now().Unix()-lastReqTime > 10 {
					reqLimitLeft = llmRequestLimitPer10Seconds
				} else {
					time.Sleep(500 * time.Millisecond)
					continue
				}
			}
			var reqUUID string
			var req llmReq
			llmReqBuf.Range(func(key, val any) bool {
				if key.(string) != "" {
					reqUUID = key.(string)
					req = val.(llmReq)
					return false
				}
				return true
			})
			llmReqBuf.Delete(reqUUID)
			if reqUUID == "" {
				continue
			}
			lastReqTime = time.Now().Unix()
			reqLimitLeft--
			resp, err := requestLLM(req)
			if err != nil {
				llmErrBuf.Store(reqUUID, err)
			} else {
				llmRespBuf.Store(reqUUID, resp)
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()
}

func CallLLM(prompt string, temperature float64) (string, error) {
	reqUUID := uuid.New().String()
	req := llmReq{
		prompt:      prompt,
		temperature: temperature,
	}
	llmReqBuf.Store(reqUUID, req)
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
