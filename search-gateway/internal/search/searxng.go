package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SearxngProvider connects the gateway to an external SearXNG instance.
type SearxngProvider struct {
	client  *http.Client
	address string
	tracer  trace.Tracer
}

type searxngResult struct {
	URL     string   `json:"url"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Engine  string   `json:"engine"`
	Parsed  []string `json:"parsed_url"`
	Score   float64  `json:"score"`
}

type searxngResponse struct {
	Query   string          `json:"query"`
	Results []searxngResult `json:"results"`
}

// NewSearxngProvider creates a SearXNG provider that targets the given base URL.
func NewSearxngProvider(address string) *SearxngProvider {
	return &SearxngProvider{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		address: address,
		tracer:  otel.Tracer("search-gateway"),
	}
}

// Search sends the query to SearXNG and applies Searchbase's optional local
// result cap because SearXNG does not expose a per-request result count.
func (p *SearxngProvider) Search(ctx context.Context, req Request) (Results, error) {
	ctx, span := p.tracer.Start(ctx, "SearxngProvider.Search")
	defer span.End()

	span.SetAttributes(
		attribute.String("search.provider", "searxng"),
	)

	sr, err := p.newRequest(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(sr)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "http request failed")
		return nil, fmt.Errorf("failed to execute searxng request: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code from provider: %d", resp.StatusCode)
		span.RecordError(err)
		span.SetStatus(codes.Error, "non-200 response")
		return nil, err
	}

	var parsedResp searxngResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to decode response")
		return nil, fmt.Errorf("failed to decode provider response: %w", err)
	}

	var results Results
	for i, r := range parsedResp.Results {
		if req.Limit != 0 {
			if i >= req.Limit {
				break
			}
		}
		results = append(results, Result{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Content,
		})
	}

	span.SetAttributes(attribute.Int("search.results_count", len(results)))
	return results, nil
}

func (p *SearxngProvider) newRequest(ctx context.Context, req Request) (*http.Request, error) {
	searchUrl, err := url.Parse(strings.TrimRight(p.address, "/") + "/search")
	if err != nil {
		return nil, fmt.Errorf("invalid searxng address: %w", err)
	}

	q := searchUrl.Query()
	q.Set("q", req.Query)
	q.Set("format", "json")

	if req.Engine != "" && req.Engine != "auto" {
		q.Set("engines", req.Engine)
	}

	if req.Page > 0 {
		q.Set("pageno", strconv.Itoa(req.Page))
	}

	if req.TimeLimit != "" {
		switch req.TimeLimit {
		case "d":
			q.Set("time_range", "day")
		case "w":
			q.Set("time_range", "week")
		case "m":
			q.Set("time_range", "month")
		case "y":
			q.Set("time_range", "year")
		}
	}

	// Safesearch: 0=none, 1=moderate, 2=strict
	if req.SafeSearch != "" {
		switch req.SafeSearch {
		case "off":
			q.Set("safesearch", "0")
		case "moderate":
			q.Set("safesearch", "1")
		case "on", "strict":
			q.Set("safesearch", "2")
		}
	}

	searchUrl.RawQuery = q.Encode()
	searxngRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, searchUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return searxngRequest, nil
}
