package search

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// DuckDuckGoProvider implements SearchProvider by scraping DuckDuckGo HTML.
type DuckDuckGoProvider struct {
	client    *http.Client
	searchURL string
	tracer    trace.Tracer
}

// NewDuckDuckGoProvider creates a new instance of the provider.
func NewDuckDuckGoProvider() *DuckDuckGoProvider {
	return &DuckDuckGoProvider{
		client: &http.Client{
			// TODO: Add any custom proxy/timeout configs here later
		},
		searchURL: "https://html.duckduckgo.com/html/",
		tracer:    otel.Tracer("search-gateway"),
	}
}

// Search scrapes html.duckduckgo.com for the top N results.
func (p *DuckDuckGoProvider) Search(ctx context.Context, sr Request) (Results, error) {
	ctx, span := p.tracer.Start(ctx, "DuckDuckGoProvider.Search")
	defer span.End()

	req, err := p.newRequest(ctx, sr)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create new request")
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "search request failed")
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		span.RecordError(err)
		span.SetStatus(codes.Error, "non-200 response")
		return nil, fmt.Errorf("unexpected status code from search provider: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to parse search results")
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	var results Results

	// DDG HTML results are usually wrapped in `.result__snippet` and `.result__title`
	doc.Find(".result__body").Each(func(i int, s *goquery.Selection) {
		if sr.Limit != 0 && len(results) >= sr.Limit {
			return
		}

		if s.Closest(".result").HasClass("result--ad") {
			return
		}

		titleEl := s.Find(".result__title .result__a")
		title := strings.TrimSpace(titleEl.Text())
		link, exists := titleEl.Attr("href")
		if !exists || title == "" {
			return
		}

		// DDG proxy links often need to be decoded
		if strings.HasPrefix(link, "//duckduckgo.com/l/?uddg=") {
			u, err := url.Parse("https:" + link)
			if err == nil {
				if uddg := u.Query().Get("uddg"); uddg != "" {
					link = uddg
				}
			}
		}

		snippet := strings.TrimSpace(s.Find(".result__snippet").Text())

		results = append(results, Result{
			Title:   title,
			URL:     link,
			Snippet: snippet,
		})
	})
	span.SetAttributes(attribute.Int("search.results_count", len(results)))

	return results, nil
}

func (p *DuckDuckGoProvider) newRequest(ctx context.Context, req Request) (*http.Request, error) {

	reqBody := url.Values{}
	reqBody.Set("q", req.Query)
	reqBody.Set("b", "")
	reqBody.Set("kl", "")

	ddgRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, p.searchURL, strings.NewReader(reqBody.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	ddgRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ddgRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	return ddgRequest, nil
}
