package search

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDuckDuckGoProvider_Search(t *testing.T) {
	mockHTML := `
<!DOCTYPE html>
<html>
<body>
	<div class="result">
		<div class="result__body">
			<h2 class="result__title">
				<a class="result__a" href="//duckduckgo.com/l/?uddg=https://example.com/1">Example 1</a>
			</h2>
			<a class="result__snippet" href="#">Snippet 1</a>
		</div>
	</div>
	<div class="result result--ad">
		<div class="result__body">
			<h2 class="result__title">
				<a class="result__a" href="https://example.com/ad">Ad</a>
			</h2>
			<a class="result__snippet" href="#">Ad Snippet</a>
		</div>
	</div>
	<div class="result">
		<div class="result__body">
			<h2 class="result__title">
				<a class="result__a" href="https://example.com/2">Example 2</a>
			</h2>
			<a class="result__snippet" href="#">Snippet 2</a>
		</div>
	</div>
</body>
</html>
`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(mockHTML))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}))
	defer ts.Close()

	provider := NewDuckDuckGoProvider()
	provider.searchURL = ts.URL

	req := Request{
		Query: "test",
		Limit: 2,
	}

	results, err := provider.Search(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0].Title != "Example 1" || results[0].URL != "https://example.com/1" || results[0].Snippet != "Snippet 1" {
		t.Errorf("unexpected first result: %+v", results[0])
	}

	if results[1].Title != "Example 2" || results[1].URL != "https://example.com/2" || results[1].Snippet != "Snippet 2" {
		t.Errorf("unexpected second result: %+v", results[1])
	}
}
