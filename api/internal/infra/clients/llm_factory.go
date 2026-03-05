package clients

import "fmt"

// NewLLMClient creates an LLMClient based on the provider name.
// Supported providers: "stub", "anthropic", "gemini".
// An empty provider defaults to "stub".
func NewLLMClient(provider, apiKey string) (LLMClient, error) {
	switch provider {
	case "", "stub":
		return NewStubLLMClient(), nil
	case "anthropic":
		if apiKey == "" {
			return nil, fmt.Errorf("LLM_API_KEY is required for provider %q", provider)
		}
		// TODO: return real Anthropic client
		return NewStubLLMClient(), nil
	case "gemini":
		if apiKey == "" {
			return nil, fmt.Errorf("LLM_API_KEY is required for provider %q", provider)
		}
		// TODO: return real Gemini client
		return NewStubLLMClient(), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %q (supported: stub, anthropic, gemini)", provider)
	}
}
