package app

import (
	stdjson "encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/tmc/langchaingo/httputil"
)

type routingTransport func(*http.Request) (*http.Response, error)

func (f routingTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestCallLLMLocalRouting(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	oldClient := httputil.DefaultClient
	var requests atomic.Int32
	httputil.DefaultClient = &http.Client{Transport: routingTransport(func(r *http.Request) (*http.Response, error) {
		requests.Add(1)
		if r.URL.Host != "ollama:11434" || r.URL.Path != "/api/chat" || r.Method != http.MethodPost {
			return nil, fmt.Errorf("unexpected destination: %s %s", r.Method, r.URL)
		}
		var body struct {
			Model    string `json:"model"`
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
			Options struct {
				NumCtx int `json:"num_ctx"`
			} `json:"options"`
		}
		if err := stdjson.NewDecoder(r.Body).Decode(&body); err != nil {
			return nil, err
		}
		if body.Model != "gemma3:12b" || len(body.Messages) != 1 || body.Messages[0].Content != "local-routing-prompt" || body.Options.NumCtx != 10240 {
			return nil, fmt.Errorf("unexpected request: %+v", body)
		}
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/x-ndjson"}}, Body: io.NopCloser(strings.NewReader("{\"message\":{\"role\":\"assistant\",\"content\":\"local-only-response\"},\"done\":true}\n")), Request: r}, nil
	})}
	t.Cleanup(func() { httputil.DefaultClient = oldClient })
	got, err := CallLLM("local-routing-prompt", 0.1, LLM_Local)
	if err != nil {
		t.Fatal(err)
	}
	if got != "local-only-response" {
		t.Fatalf("got %q", got)
	}
	if requests.Load() != 1 {
		t.Fatalf("Ollama requests=%d", requests.Load())
	}
}

func TestCallLLMGeminiRouting(t *testing.T) {
	oldInterval := llmGeminiRequestInterval
	llmGeminiRequestInterval = 0
	t.Cleanup(func() { llmGeminiRequestInterval = oldInterval })
	var geminiRequests, ollamaRequests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		geminiRequests.Add(1)
		if !strings.Contains(r.URL.Path, "gemini-3.5-flash-lite:generateContent") {
			t.Errorf("unexpected model endpoint: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"candidates":[{"content":{"parts":[{"text":"gemini-only-response"}]},"finishReason":"STOP"}]}`)
	}))
	defer server.Close()
	t.Setenv("GEMINI_API_KEY", "test-key")
	t.Setenv("GOOGLE_GEMINI_BASE_URL", server.URL)
	oldClient := httputil.DefaultClient
	httputil.DefaultClient = &http.Client{Transport: routingTransport(func(r *http.Request) (*http.Response, error) {
		ollamaRequests.Add(1)
		return nil, fmt.Errorf("Gemini request unexpectedly reached Ollama")
	})}
	t.Cleanup(func() { httputil.DefaultClient = oldClient })
	got, err := CallLLM("gemini-routing-prompt", 0.1, LLM_Gemini)
	if err != nil {
		t.Fatal(err)
	}
	if got != "gemini-only-response" || geminiRequests.Load() != 1 || ollamaRequests.Load() != 0 {
		t.Fatalf("response=%q Gemini requests=%d Ollama requests=%d", got, geminiRequests.Load(), ollamaRequests.Load())
	}
}

func TestCallLLMLocalError(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	oldClient := httputil.DefaultClient
	httputil.DefaultClient = &http.Client{Transport: routingTransport(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"error":"local-unavailable"}`)), Request: r}, nil
	})}
	t.Cleanup(func() { httputil.DefaultClient = oldClient })
	got, err := CallLLM("local-error-prompt", 0.1, LLM_Local)
	if err == nil || got != "" || !strings.Contains(err.Error(), "local-unavailable") {
		t.Fatalf("response=%q error=%v", got, err)
	}
}
