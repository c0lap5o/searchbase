package search

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDDGSProvider_Search(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/search/text" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
			"results": [
				{
					"title": "test",
					"href": "https://example.com",
					"body": "test search result"
				}
			]
		}`))
		if err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer mockServer.Close()

	provider := NewDDGSProvider(mockServer.URL)
	req := Request{
		Query: "test",
		Limit: 10,
	}

	results, err := provider.Search(context.Background(), req)
	if err != nil {
		t.Fatalf("Search() returned unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Title != "test" {
		t.Errorf("expected title 'test', got '%s'", results[0].Title)
	}

	if results[0].URL != "https://example.com" {
		t.Errorf("expected URL 'https://example.com', got '%s'", results[0].URL)
	}

	if results[0].Snippet != "test search result" {
		t.Errorf("expected snippet 'test search result', got '%s'", results[0].Snippet)
	}
}

func TestDDGSProvider_Search_ErrorResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer mockServer.Close()

	provider := NewDDGSProvider(mockServer.URL)
	_, err := provider.Search(context.Background(), Request{Query: "test"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "ddgs search provider returned status 502" {
		t.Fatalf("expected safe provider error, got %v", err)
	}
}
