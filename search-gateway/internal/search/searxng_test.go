package search

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearxngProvider_Search_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters mapping
		if r.URL.Query().Get("q") != "test query" {
			t.Errorf("Expected query 'test query', got '%s'", r.URL.Query().Get("q"))
		}
		if r.URL.Query().Get("format") != "json" {
			t.Errorf("Expected format 'json', got '%s'", r.URL.Query().Get("format"))
		}
		if r.URL.Query().Get("engines") != "wikipedia" {
			t.Errorf("Expected engines 'wikipedia', got '%s'", r.URL.Query().Get("engines"))
		}
		if r.URL.Query().Get("pageno") != "3" {
			t.Errorf("Expected pageno '3', got '%s'", r.URL.Query().Get("pageno"))
		}
		if r.URL.Query().Get("time_range") != "week" {
			t.Errorf("Expected time_range 'week', got '%s'", r.URL.Query().Get("time_range"))
		}
		if r.URL.Query().Get("safesearch") != "2" {
			t.Errorf("Expected safesearch '2', got '%s'", r.URL.Query().Get("safesearch"))
		}

		w.WriteHeader(http.StatusOK)
		mockResp := searxngResponse{
			Query: "test query",
			Results: []searxngResult{
				{URL: "https://example.com/1", Title: "Result 1", Content: "Snippet 1"},
				{URL: "https://example.com/2", Title: "Result 2", Content: "Snippet 2"},
			},
		}
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer mockServer.Close()

	provider := NewSearxngProvider(mockServer.URL)
	req := Request{
		Query:      "test query",
		Limit:      5,
		Engine:     "wikipedia",
		Page:       3,
		TimeLimit:  "w",
		SafeSearch: "strict",
	}

	results, err := provider.Search(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0].Title != "Result 1" || results[0].URL != "https://example.com/1" || results[0].Snippet != "Snippet 1" {
		t.Errorf("First result mismatch: %+v", results[0])
	}
}

func TestSearxngProvider_Search_ErrorResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	provider := NewSearxngProvider(mockServer.URL)
	req := Request{Query: "test"}

	_, err := provider.Search(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestSearxngProvider_Search_Limit(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		mockResp := searxngResponse{
			Results: []searxngResult{
				{URL: "1"}, {URL: "2"}, {URL: "3"}, {URL: "4"},
			},
		}
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer mockServer.Close()

	provider := NewSearxngProvider(mockServer.URL)
	req := Request{Query: "test", Limit: 2} // Request limit 2

	results, err := provider.Search(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results due to limit, got %d", len(results))
	}
}
