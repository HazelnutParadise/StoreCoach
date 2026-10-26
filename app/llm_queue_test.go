package app

import (
	"bytes"
	"context"
	stdjson "encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/tmc/langchaingo/httputil"
)

const llmCPUProbeEnv = "STORECOACH_LLM_CPU_PROBE"

func TestLLMCPUProbeHelper(t *testing.T) {
	mode := os.Getenv(llmCPUProbeEnv)
	if mode == "" {
		t.Skip("subprocess-only CPU probe")
	}
	switch mode {
	case "idle":
		time.Sleep(750 * time.Millisecond)
	case "waiting":
		t.Setenv("GEMINI_API_KEY", "")
		oldClient := httputil.DefaultClient
		httputil.DefaultClient = &http.Client{Transport: routingTransport(func(r *http.Request) (*http.Response, error) {
			time.Sleep(750 * time.Millisecond)
			return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/x-ndjson"}}, Body: io.NopCloser(strings.NewReader("{\"message\":{\"role\":\"assistant\",\"content\":\"waiting-probe-response\"},\"done\":true}\n")), Request: r}, nil
		})}
		defer func() { httputil.DefaultClient = oldClient }()
		got, err := CallLLM("waiting-probe-prompt", 0.1, LLM_Local)
		if err != nil || got != "waiting-probe-response" {
			t.Fatalf("CallLLM response=%q error=%v", got, err)
		}
	default:
		t.Fatalf("unknown CPU probe mode %q", mode)
	}
}

func assertLLMSubprocessCPU(t *testing.T, mode string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestLLMCPUProbeHelper$", "-test.count=1")
	cmd.Env = append(os.Environ(), llmCPUProbeEnv+"="+mode)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	started := time.Now()
	err := cmd.Run()
	wall := time.Since(started)
	if ctx.Err() != nil {
		t.Fatalf("probe mode %s timed out: %v\nstdout/stderr:\n%s", mode, ctx.Err(), output.String())
	}
	if err != nil {
		t.Fatalf("probe mode %s failed: %v\nstdout/stderr:\n%s", mode, err, output.String())
	}
	state := cmd.ProcessState
	cpu := state.UserTime() + state.SystemTime()
	t.Logf("probe mode %s CPU=%v wall=%v", mode, cpu, wall)
	if cpu >= 300*time.Millisecond {
		t.Fatalf("probe mode %s used CPU %v (limit 300ms)\nstdout/stderr:\n%s", mode, cpu, output.String())
	}
}

func TestLLMIdleCPU(t *testing.T) { assertLLMSubprocessCPU(t, "idle") }

func TestCallLLMWaitingCPU(t *testing.T) { assertLLMSubprocessCPU(t, "waiting") }

func TestCallLLMConcurrentReplies(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	oldClient := httputil.DefaultClient
	var active, maxActive atomic.Int32
	httputil.DefaultClient = &http.Client{Transport: routingTransport(func(r *http.Request) (*http.Response, error) {
		now := active.Add(1)
		defer active.Add(-1)
		for prior := maxActive.Load(); now > prior; prior = maxActive.Load() {
			if maxActive.CompareAndSwap(prior, now) {
				break
			}
		}
		var body struct {
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := stdjson.NewDecoder(r.Body).Decode(&body); err != nil {
			return nil, err
		}
		if len(body.Messages) != 1 {
			return nil, fmt.Errorf("request has %d messages", len(body.Messages))
		}
		id := body.Messages[0].Content
		if strings.HasPrefix(id, "request-") {
			if n, err := strconv.Atoi(strings.TrimPrefix(id, "request-")); err == nil && n%7 == 0 {
				return &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(fmt.Sprintf(`{"error":"failure-%s"}`, id))), Request: r}, nil
			}
		}
		return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/x-ndjson"}}, Body: io.NopCloser(strings.NewReader(fmt.Sprintf("{\"message\":{\"role\":\"assistant\",\"content\":%q},\"done\":true}\n", "reply-"+id))), Request: r}, nil
	})}
	defer func() { httputil.DefaultClient = oldClient }()

	type result struct {
		id   int
		text string
		err  error
	}
	results := make(chan result, 64)
	var callers sync.WaitGroup
	for id := 0; id < 64; id++ {
		callers.Add(1)
		go func(id int) {
			defer callers.Done()
			text, err := CallLLM(fmt.Sprintf("request-%d", id), 0.1, LLM_Local)
			results <- result{id: id, text: text, err: err}
		}(id)
	}
	done := make(chan struct{})
	go func() { callers.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(20 * time.Second):
		t.Fatal("64 concurrent CallLLM requests did not finish within 20 seconds")
	}
	close(results)
	for got := range results {
		wantFailure := got.id%7 == 0
		if wantFailure {
			if got.err == nil || got.text != "" || !strings.Contains(got.err.Error(), fmt.Sprintf("request-%d", got.id)) {
				t.Errorf("request-%d response=%q error=%v, want its own error", got.id, got.text, got.err)
			}
		} else if got.err != nil || got.text != fmt.Sprintf("reply-request-%d", got.id) {
			t.Errorf("request-%d response=%q error=%v", got.id, got.text, got.err)
		}
	}
	if maxActive.Load() > 1 {
		t.Errorf("concurrent HTTP calls=%d, want at most one", maxActive.Load())
	}
}
