package app

import (
	"fmt"
	"testing"
	"time"
)

func TestLLMQueueFIFOAndRoundRobin(t *testing.T) {
	queue := newLLMQueue()
	for i := range 40 {
		queue.push(llmReq{model: LLM_Gemini, prompt: fmt.Sprintf("gemini-%d", i)})
		queue.push(llmReq{model: LLM_Local, prompt: fmt.Sprintf("local-%d", i)})
	}
	for i := range 40 {
		for _, prefix := range []string{"gemini", "local"} {
			if got := queue.pop().prompt; got != fmt.Sprintf("%s-%d", prefix, i) {
				t.Fatalf("unexpected queue order: %q at round %d", got, i)
			}
		}
	}
}

func TestLLMQueueResumesOtherModel(t *testing.T) {
	queue := newLLMQueue()
	queue.push(llmReq{model: LLM_Local, prompt: "local-1"})
	queue.push(llmReq{model: LLM_Local, prompt: "local-2"})
	if got := queue.pop().prompt; got != "local-1" {
		t.Fatal(got)
	}
	queue.push(llmReq{model: LLM_Gemini, prompt: "gemini-1"})
	for _, want := range []string{"gemini-1", "local-2"} {
		if got := queue.pop().prompt; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	}
}

func TestLLMQueueRepeatedWakeup(t *testing.T) {
	queue := newLLMQueue()
	received := make(chan llmReq, 1)
	go func() {
		for range 64 {
			received <- queue.pop()
		}
	}()
	for i := range 64 {
		want := fmt.Sprint(i)
		queue.push(llmReq{model: i % 2, prompt: want})
		select {
		case got := <-received:
			if got.prompt != want {
				t.Fatalf("got %q, want %q", got.prompt, want)
			}
		case <-time.After(time.Second):
			t.Fatal("queue did not wake after submission")
		}
	}
}

func TestCallLLMRejectsInvalidModel(t *testing.T) {
	for _, model := range []int{-1, 2, 100} {
		if text, err := CallLLM("invalid", 0, model); err == nil || text != "" {
			t.Fatalf("model %d: text=%q, error=%v", model, text, err)
		}
	}
}
