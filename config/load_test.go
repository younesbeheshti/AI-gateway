package config

import (
	"os"
	"testing"
)

func TestLoad_OpenRouterAPIKey(t *testing.T) {
	os.Setenv("AI_GATEWAY_OPENROUTER_API_KEY", "test-key-12345")
	defer os.Unsetenv("AI_GATEWAY_OPENROUTER_API_KEY")

	cfg, err := Load("../config.yml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.OpenRouter.APIKey != "test-key-12345" {
		t.Errorf("expected APIKey to be 'test-key-12345', got '%s'", cfg.OpenRouter.APIKey)
	}
}

func TestLoad_OpenRouterDoubleUnderscore(t *testing.T) {
	os.Setenv("AI_GATEWAY_OPENROUTER__API_KEY", "test-key-double-underscore")
	defer os.Unsetenv("AI_GATEWAY_OPENROUTER__API_KEY")

	cfg, err := Load("../config.yml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.OpenRouter.APIKey != "test-key-double-underscore" {
		t.Errorf("expected APIKey to be 'test-key-double-underscore', got '%s'", cfg.OpenRouter.APIKey)
	}
}
