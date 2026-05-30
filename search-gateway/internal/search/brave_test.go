package search

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBraveProvider_Search(t *testing.T) {
	mockResponse := `{
		"type": "search",
		"query": {
			"original": "brave search",
			"more_results_available": false
		},
		"web": {
			"type": "search",
			"results": [
				{
					"type": "search_result",
					"title": "search result",
					"url": "https://example.com",
					"description": "first search result"
				},
				{
					"type": "search_result",
					"title": "search result 2",
					"url": "https://example-2.com",
					"description": "second search result"
				}
			]
		}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(mockResponse))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}))
	defer ts.Close()
	provider := NewBraveProvider("foo")
	provider.webSearchURL = ts.URL

	t.Run("Success", func(t *testing.T) {
		expected := Results{
			{
				Title:   "search result",
				URL:     "https://example.com",
				Snippet: "first search result",
			},
			{
				Title:   "search result 2",
				URL:     "https://example-2.com",
				Snippet: "second search result",
			},
		}

		req := Request{
			Query: "brave search",
		}

		results, err := provider.Search(context.Background(), req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(results) != len(expected) {
			t.Fatalf("Expected %d results, got %d", len(expected), len(results))
		}

		for i, result := range results {
			if result.Title != expected[i].Title {
				t.Fatalf("Expected title %s, got %s", expected[i].Title, result.Title)
			}

			if result.URL != expected[i].URL {
				t.Fatalf("Expected URL %s, got %s", expected[i].URL, result.URL)
			}

			if result.Snippet != expected[i].Snippet {
				t.Fatalf("Expected snippet %s, got %s", expected[i].Snippet, result.Snippet)
			}
		}
	})

	t.Run("Error", func(t *testing.T) {
		req := Request{}
		_, err := provider.Search(context.Background(), req)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}

func TestBraveProvider_newRequest(t *testing.T) {
	provider := NewBraveProvider("foo")

	t.Run("region and page offset", func(t *testing.T) {
		req, err := provider.newRequest(context.Background(), Request{
			Query:      "brave search",
			Limit:      10,
			Page:       3,
			Region:     "us-en",
			TimeLimit:  "w",
			SafeSearch: "on",
		}, "https://example.com/search")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if req.Header.Get("x-subscription-token") != "foo" {
			t.Fatalf("expected x-subscription-token foo, got %s", req.Header.Get("x-subscription-token"))
		}

		if req.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("expected Content-Type application/json, got %s", req.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("expected no error reading body, got %v", err)
		}

		var braveReq braveRequest
		if err := json.Unmarshal(body, &braveReq); err != nil {
			t.Fatalf("expected no error unmarshalling body, got %v", err)
		}

		if braveReq.Country != "US" {
			t.Fatalf("expected country US, got %s", braveReq.Country)
		}
		if braveReq.SearchLang != "en" {
			t.Fatalf("expected search_lang en, got %s", braveReq.SearchLang)
		}
		if braveReq.Offset != 2 {
			t.Fatalf("expected offset 2, got %d", braveReq.Offset)
		}
		if braveReq.Freshness != "pw" {
			t.Fatalf("expected freshness pw, got %s", braveReq.Freshness)
		}
		if braveReq.SafeSearch != "strict" {
			t.Fatalf("expected safesearch strict, got %s", braveReq.SafeSearch)
		}
	})

	t.Run("ignores unsupported global region", func(t *testing.T) {
		req, err := provider.newRequest(context.Background(), Request{
			Query:  "brave search",
			Region: "wt-wt",
		}, "https://example.com/search")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("expected no error reading body, got %v", err)
		}

		var braveReq braveRequest
		if err := json.Unmarshal(body, &braveReq); err != nil {
			t.Fatalf("expected no error unmarshalling body, got %v", err)
		}

		if braveReq.Country != "" {
			t.Fatalf("expected empty country, got %s", braveReq.Country)
		}
		if braveReq.SearchLang != "" {
			t.Fatalf("expected empty search_lang, got %s", braveReq.SearchLang)
		}
	})
}
