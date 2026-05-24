package mcp

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/coolapso/searchbase/search-gateway/internal/scraper"
	"github.com/coolapso/searchbase/search-gateway/internal/search"
	"github.com/mark3labs/mcp-go/mcp"
)

// MockSearchProvider implements search.SearchProvider for testing
type MockSearchProvider struct{}

func (m *MockSearchProvider) Search(ctx context.Context, req search.Request) (search.Results, error) {
	return search.Results{
		{
			Title:   "Test Title",
			URL:     "http://example.com/mock",
			Snippet: "Test snippet",
		},
	}, nil
}

func TestHandleWebSearch(t *testing.T) {
	// Mock scraper
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
	srv := NewServer(&MockSearchProvider{}, scraperClient, logger)

	// Call HandleWebSearch directly
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "web_search",
			Arguments: map[string]any{
				"query": "test query",
				"limit": 1,
			},
		},
	}

	res, err := srv.HandleWebSearch(context.Background(), req)
	if err != nil {
		t.Fatalf("HandleWebSearch returned unexpected error: %v", err)
	}

	if res.IsError {
		t.Fatalf("Expected success, got tool error")
	}

	if len(res.Content) == 0 {
		t.Fatalf("Expected content, got none")
	}

	textContent := res.Content[0].(mcp.TextContent).Text
	if !strings.Contains(textContent, "Test Title") {
		t.Errorf("Expected title in response, got %s", textContent)
	}
}

func TestHandleFetchURL(t *testing.T) {
	// Mock scraper
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
	srv := NewServer(&MockSearchProvider{}, scraperClient, logger)

	// Call HandleFetchURL directly
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "fetch_url",
			Arguments: map[string]any{
				"url": "http://example.com",
			},
		},
	}

	res, err := srv.HandleFetchURL(context.Background(), req)
	if err != nil {
		t.Fatalf("HandleFetchURL returned unexpected error: %v", err)
	}

	if res.IsError {
		t.Fatalf("Expected success, got tool error")
	}

	if len(res.Content) == 0 {
		t.Fatalf("Expected content, got none")
	}

	textContent := res.Content[0].(mcp.TextContent).Text
	if !strings.Contains(textContent, "# Fetched Markdown") {
		t.Errorf("Expected markdown in response, got %s", textContent)
	}
}
