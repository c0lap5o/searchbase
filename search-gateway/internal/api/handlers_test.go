package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coolapso/searchbase/search-gateway/internal/scraper"
	"github.com/coolapso/searchbase/search-gateway/internal/search"
	"github.com/gin-gonic/gin"
)

// MockSearchProvider implements search.SearchProvider for testing
type MockSearchProvider struct{}

func (m *MockSearchProvider) Search(ctx context.Context, req search.Request) (search.Results, error) {
	return search.Results{
		{
			Title:   "Test Result",
			URL:     "http://example.com",
			Snippet: "Test snippet",
		},
	}, nil
}

func TestHandleSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock scraper server
	mockScraperServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/extract" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(scraper.ExtractResponse{
				Success: true,
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockScraperServer.Close()

	scraperClient := scraper.NewScraperClient(mockScraperServer.URL)
	logger := slog.Default()
	apiServer := NewAPIServer(&MockSearchProvider{}, scraperClient, logger)

	router := gin.New()
	apiServer.RegisterRoutes(router.Group("/api/v1"))

	// Perform the test request
	reqBody := search.Request{
		Query: "test query",
		Limit: 1,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/search", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", w.Code)
	}

	var resp search.Results
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(resp) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(resp))
	}
}

func TestHandleFetch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock scraper server
	mockScraperServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/extract" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(scraper.ExtractResponse{
				Markdown: "# Fetched Markdown",
				Success:  true,
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockScraperServer.Close()

	scraperClient := scraper.NewScraperClient(mockScraperServer.URL)
	logger := slog.Default()
	apiServer := NewAPIServer(&MockSearchProvider{}, scraperClient, logger)

	router := gin.New()
	apiServer.RegisterRoutes(router.Group("/api/v1"))

	// Perform the test request
	reqBody := FetchRequest{
		URL: "http://example.com",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/fetch", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", w.Code)
	}

	var resp FetchResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Markdown != "# Fetched Markdown" {
		t.Errorf("Expected fetched markdown, got %s", resp.Markdown)
	}
}
