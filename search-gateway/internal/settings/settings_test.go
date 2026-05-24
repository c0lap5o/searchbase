package settings

import (
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
