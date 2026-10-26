package app

import (
	stdjson "encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/tmc/langchaingo/httputil"
)

func TestCallLLMMixedModels(t *testing.T) {
	oldInterval := llmGeminiRequestInterval
	llmGeminiRequestInterval = 0
	t.Cleanup(func() { llmGeminiRequestInterval = oldInterval })
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Contents []struct{ Parts []struct{ Text string } }
		}
		if err := stdjson.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Contents) != 1 || len(body.Contents[0].Parts) != 1 {
			t.Error("invalid Gemini request")
			http.Error(w, "invalid request", 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = stdjson.NewEncoder(w).Encode(map[string]any{"candidates": []any{map[string]any{
			"content": map[string]any{"parts": []any{map[string]any{"text": "gemini:" + body.Contents[0].Parts[0].Text}}}, "finishReason": "STOP",
		}}})
	}))
	defer server.Close()
	t.Setenv("GEMINI_API_KEY", "test-key")
	t.Setenv("GOOGLE_GEMINI_BASE_URL", server.URL)
	oldClient := httputil.DefaultClient
	httputil.DefaultClient = &http.Client{Transport: routingTransport(func(r *http.Request) (*http.Response, error) {
		var body struct{ Messages []struct{ Content string } }
		if err := stdjson.NewDecoder(r.Body).Decode(&body); err != nil {
			return nil, err
		}
		if len(body.Messages) != 1 {
			return nil, fmt.Errorf("invalid Ollama request")
		}
		payload, _ := stdjson.Marshal(map[string]any{"message": map[string]any{"role": "assistant", "content": "local:" + body.Messages[0].Content}, "done": true})
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/x-ndjson"}}, Body: io.NopCloser(strings.NewReader(string(payload) + "\n")), Request: r}, nil
	})}
	defer func() { httputil.DefaultClient = oldClient }()
	results := make(chan error, 32)
	for i := range 32 {
		go func(i int) {
			model := i % 2
			prompt := fmt.Sprint(i)
			text, err := CallLLM(prompt, 0.1, model)
			prefix := "gemini:"
			if model == LLM_Local {
				prefix = "local:"
			}
			if err == nil && text != prefix+prompt {
				err = fmt.Errorf("model %d prompt %s got %q", model, prompt, text)
			}
			results <- err
		}(i)
	}
	timeout := time.After(10 * time.Second)
	for range 32 {
		select {
		case err := <-results:
			if err != nil {
				t.Error(err)
			}
		case <-timeout:
			t.Fatal("mixed model requests did not complete")
		}
	}
}
