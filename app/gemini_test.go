package app

import (
	stdjson "encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeminiRequest(t *testing.T) {
	oldInterval := llmGeminiRequestInterval
	llmGeminiRequestInterval = 0
	t.Cleanup(func() { llmGeminiRequestInterval = oldInterval })
	cases := []struct {
		name, body, want string
		status           int
		wantErr          bool
	}{
		{"text", `{"candidates":[{"content":{"parts":[{"text":"hello "},{"text":"world"}],"role":"model"},"finishReason":"STOP"}]}`, "hello world", 200, false},
		{"continuation stays single request", `{"candidates":[{"content":{"parts":[{"text":"partial"}]},"finishReason":"MAX_TOKENS"}],"continuationToken":"next-page"}`, "partial", 200, false},
		{"blocked prompt", `{"promptFeedback":{"blockReason":"SAFETY"}}`, "", 200, true},
		{"blocked output", `{"candidates":[{"content":{"parts":[{"text":"partial"}]},"finishReason":"SAFETY"}]}`, "", 200, true},
		{"empty", `{}`, "", 200, true},
		{"retired model", `{"error":{"code":404,"message":"This model is no longer available","status":"NOT_FOUND"}}`, "", 404, true},
		{"upstream error", `{"error":{"code":503,"message":"unavailable","status":"UNAVAILABLE"}}`, "", 503, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if !strings.Contains(r.URL.Path, "gemini-3.5-flash-lite:generateContent") {
					t.Errorf("model path: %s", r.URL.Path)
				}
				if r.Header.Get("x-goog-api-key") != "test-key" {
					t.Error("missing configured API key")
				}
				var payload struct {
					Contents []struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					} `json:"contents"`
					GenerationConfig struct {
						Temperature float64 `json:"temperature"`
					} `json:"generationConfig"`
				}
				if err := stdjson.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Error(err)
				}
				if len(payload.Contents) != 1 || len(payload.Contents[0].Parts) != 1 || payload.Contents[0].Parts[0].Text != "prompt" {
					t.Errorf("unexpected prompt: %+v", payload.Contents)
				}
				if payload.GenerationConfig.Temperature != 0.5 {
					t.Errorf("temperature: %v", payload.GenerationConfig.Temperature)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			t.Setenv("GEMINI_API_KEY", "test-key")
			t.Setenv("GOOGLE_API_KEY", "other-key")
			t.Setenv("GOOGLE_GEMINI_BASE_URL", server.URL)
			got, err := requestLLM(llmReq{prompt: "prompt", temperature: 0.5, model: LLM_Gemini})
			if (err != nil) != tc.wantErr {
				t.Fatalf("error=%v, wantErr=%v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("got %q want %q", got, tc.want)
			}
			if calls != 1 {
				t.Errorf("requests=%d want 1", calls)
			}
		})
	}
}

func TestGeminiMissingKey(t *testing.T) {
	oldInterval := llmGeminiRequestInterval
	llmGeminiRequestInterval = 0
	t.Cleanup(func() { llmGeminiRequestInterval = oldInterval })
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "must-not-use")
	got, err := requestLLM(llmReq{model: LLM_Gemini})
	if err == nil || got != "" {
		t.Fatalf("got %q, error=%v", got, err)
	}
}
