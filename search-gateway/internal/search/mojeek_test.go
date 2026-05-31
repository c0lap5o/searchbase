package search

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMojeekProvider_Search_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("api_key") != "foo" {
			t.Errorf("Expected api key 'foo', got '%s'", r.URL.Query().Get("api_key"))
		}
		if r.URL.Query().Get("q") != "test query" {
			t.Errorf("Expected query 'test query', got '%s'", r.URL.Query().Get("q"))
		}
		if r.URL.Query().Get("fmt") != "json" {
			t.Errorf("Expected format 'json', got '%s'", r.URL.Query().Get("fmt"))
		}
		if r.URL.Query().Get("t") != "5" {
			t.Errorf("Expected limit '5', got '%s'", r.URL.Query().Get("t"))
		}
		if r.URL.Query().Get("since") != "day" {
			t.Errorf("Expected since 'day', got '%s'", r.URL.Query().Get("since"))
		}
		if r.URL.Query().Get("safe") != "1" {
			t.Errorf("Expected safe '1', got '%s'", r.URL.Query().Get("safe"))
		}

		w.WriteHeader(http.StatusOK)
		mockResp := map[string]any{
			"response": map[string]any{
				"status": "OK",
				"results": []map[string]any{
					{
						"url":   "https://example.com/1",
						"title": "Result 1",
						"desc":  "Snippet 1",
						"size":  "1kb",
					},
					{
						"url":   "https://example.com/2",
						"title": "Result 2",
						"desc":  "Snippet 2",
						"score": 12.34,
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer mockServer.Close()

	provider := NewMojeekProvider("foo")
	provider.address = mockServer.URL
	provider.apiKey = "foo"
	req := Request{
		Query:      "test query",
		Limit:      5,
		TimeLimit:  "day",
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

func TestMojeekProvider_Search_ErrorResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	provider := NewMojeekProvider("foo")
	provider.address = mockServer.URL
	req := Request{Query: "test"}

	_, err := provider.Search(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "mojeek search provider returned status 500" {
		t.Fatalf("Expected safe provider error, got %v", err)
	}
}

func TestMojeekProvider_Search_NetworkErrorDoesNotLeakQuery(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	provider := NewMojeekProvider("foo")
	provider.address = mockServer.URL
	mockServer.Close()

	_, err := provider.Search(context.Background(), Request{Query: "secret private query"})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "mojeek search provider request failed" {
		t.Fatalf("Expected safe network error, got %v", err)
	}

	if strings.Contains(err.Error(), "secret private query") {
		t.Fatalf("Error leaked query: %v", err)
	}
}

func TestMojeekProvider_Search_ForwardsLimit(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("t") != "2" {
			t.Errorf("Expected limit '2', got '%s'", r.URL.Query().Get("t"))
		}

		w.WriteHeader(http.StatusOK)
		mockResp := map[string]any{
			"response": map[string]any{
				"results": []map[string]any{
					{"url": "1"}, {"url": "2"},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer mockServer.Close()

	provider := NewMojeekProvider("foo")
	provider.address = mockServer.URL
	req := Request{Query: "test", Limit: 2}

	results, err := provider.Search(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results returned by Mojeek, got %d", len(results))
	}
}

func TestMojeekProvider_Search_OmitsLimitWhenUnset(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("t") {
			t.Errorf("Expected no limit param, got '%s'", r.URL.Query().Get("t"))
		}

		w.WriteHeader(http.StatusOK)
		mockResp := map[string]any{
			"response": map[string]any{
				"results": []map[string]any{
					{"url": "1"}, {"url": "2"}, {"url": "3"}, {"url": "4"},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer mockServer.Close()

	provider := NewMojeekProvider("foo")
	provider.address = mockServer.URL
	req := Request{Query: "test"}

	results, err := provider.Search(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 4 {
		t.Fatalf("Expected all mock results when limit is unset, got %d", len(results))
	}
}
