package clients_test

import (
	"strings"
	"testing"

	"github.com/sibukixxx/travelist/api/internal/infra/clients"
)

func TestNewLLMClientShouldReturnStubWhenProviderIsEmpty(t *testing.T) {
	client, err := clients.NewLLMClient("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewLLMClientShouldReturnStubWhenProviderIsStub(t *testing.T) {
	client, err := clients.NewLLMClient("stub", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewLLMClientShouldReturnClientWhenAnthropicWithKey(t *testing.T) {
	client, err := clients.NewLLMClient("anthropic", "sk-ant-test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewLLMClientShouldReturnClientWhenGeminiWithKey(t *testing.T) {
	client, err := clients.NewLLMClient("gemini", "AIzaSy-test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewLLMClientShouldReturnErrorWhenAnthropicWithoutKey(t *testing.T) {
	_, err := clients.NewLLMClient("anthropic", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "LLM_API_KEY is required") {
		t.Errorf("error = %q, want containing %q", err.Error(), "LLM_API_KEY is required")
	}
}

func TestNewLLMClientShouldReturnErrorWhenGeminiWithoutKey(t *testing.T) {
	_, err := clients.NewLLMClient("gemini", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "LLM_API_KEY is required") {
		t.Errorf("error = %q, want containing %q", err.Error(), "LLM_API_KEY is required")
	}
}

func TestNewLLMClientShouldReturnErrorWhenProviderIsUnsupported(t *testing.T) {
	_, err := clients.NewLLMClient("openai", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported LLM provider") {
		t.Errorf("error = %q, want containing %q", err.Error(), "unsupported LLM provider")
	}
}
