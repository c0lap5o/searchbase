package settings

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewSettings_Defaults(t *testing.T) {
	s, err := NewSettings()
	if err != nil {
		t.Fatalf("Failed to create settings: %v", err)
	}

	if s.Port() != "808int" && s.Port() != "8080" {
		t.Errorf("Expected default port 8080 or 8088, got %s", s.Port())
	}

	if s.Port() != "8080" {
		t.Errorf("Expected default port 8080, got %s", s.Port())
	}

	if s.CrawlWorkerAddress() != "http://localhost:8000" {
		t.Errorf("Expected default crawl worker address http://localhost:8000, got %s", s.CrawlWorkerAddress())
	}

	if s.GinMode() != gin.ReleaseMode {
		t.Errorf("Expected default gin mode %s, got %s", gin.ReleaseMode, s.GinMode())
	}

	if s.LogLevel() != "error" {
		t.Errorf("Expected default log level error, got %s", s.LogLevel())
	}

	if s.SearchProvider().Name() != "searchbase_ddg" {
		t.Errorf("Expected default search provider searchbase_ddg, got %s", s.SearchProvider().Name())
	}
}

func TestNewSettings_EnvVars(t *testing.T) {
	t.Setenv("SEARCHBASE_PORT", "9090")
	t.Setenv("SEARCHBASE_CRAWL_WORKER_ADDRESS", "http://remote:8001")
	t.Setenv("SEARCHBASE_LOG_LEVEL", "debug")
	t.Setenv("SEARCHBASE_ENVIRONMENT", "dev")
	t.Setenv("SEARCHBASE_SEARCH_PROVIDER", "ddgs")
	t.Setenv("SEARCHBASE_DDGS_PROVIDER_ADDRESS", "http://ddgs-service:8002")

	s, err := NewSettings()
	if err != nil {
		t.Fatalf("Failed to create settings: %v", err)
	}

	if s.Port() != "9090" {
		t.Errorf("Expected port 9090, got %s", s.Port())
	}

	if s.CrawlWorkerAddress() != "http://remote:8001" {
		t.Errorf("Expected crawl worker address http://remote:8001, got %s", s.CrawlWorkerAddress())
	}

	if s.LogLevel() != "debug" {
		t.Errorf("Expected log level debug, got %s", s.LogLevel())
	}

	if s.GinMode() != gin.DebugMode {
		t.Errorf("Expected gin mode %s, got %s", gin.DebugMode, s.GinMode())
	}

	if s.SearchProvider().Name() != "ddgs" {
		t.Errorf("Expected search provider ddgs, got %s", s.SearchProvider().Name())
	}

	if s.SearchProvider().Address() != "http://ddgs-service:8002" {
		t.Errorf("Expected ddgs address http://ddgs-service:8002, got %s", s.SearchProvider().Address())
	}
}

func TestNewSettings_SearchProviderValidationErrors(t *testing.T) {
	tests := []struct {
		name          string
		provider      string
		expectedError string
	}{
		{
			name:          "ddgs missing address",
			provider:      "ddgs",
			expectedError: "provider address not set: searchbase provider ddgs requires SEARCHBASE_DDGS_PROVIDER_ADDRESS",
		},
		{
			name:          "searxng missing address",
			provider:      "searxng",
			expectedError: "provider address not set: searchbase provider searxng requires SEARCHBASE_SEARXNG_PROVIDER_ADDRESS",
		},
		{
			name:          "brave missing token",
			provider:      "brave",
			expectedError: "provider API Token not set: searchbase provider brave requires SEARCHBASE_BRAVE_API_TOKEN",
		},
		{
			name:          "mojeek missing token",
			provider:      "mojeek",
			expectedError: "provider API Token not set: searchbase provider mojeek requires SEARCHBASE_MOJEEK_API_KEY",
		},
		{
			name:          "unsupported provider",
			provider:      "unknown",
			expectedError: "unsupported search provider: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SEARCHBASE_SEARCH_PROVIDER", tt.provider)

			_, err := NewSettings()
			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			if err.Error() != tt.expectedError {
				t.Fatalf("Expected error %q, got %q", tt.expectedError, err.Error())
			}

			if strings.Contains(err.Error(), "SEARCHBASE_") && strings.HasSuffix(err.Error(), "SEARCHBASE_") {
				t.Fatalf("Error has incomplete environment variable name: %v", err)
			}
		})
	}
}
