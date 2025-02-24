package app

import (
	"context"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

func CallLLM(prompt string) (string, error) {
	ctx := context.Background()
	llm, err := googleai.New(ctx, googleai.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return "", err
	}

	llms.GenerateFromSinglePrompt(ctx, llm, prompt)

	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		return "", err
	}
	return completion, nil
}
