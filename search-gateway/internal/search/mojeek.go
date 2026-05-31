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

// MojeekProvider connects the gateway to the Mojeek Search API.
type MojeekProvider struct {
	client  *http.Client
	address string
	tracer  trace.Tracer
	apiKey  string
}

type mojeekResult struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type mojeekResponse struct {
	Response mojeekSearchResponse `json:"response"`
}

type mojeekSearchResponse struct {
	Results []mojeekResult `json:"results"`
}

// NewMojeekProvider creates a Mojeek provider that targets the given the api key
func NewMojeekProvider(apiKey string) *MojeekProvider {
	return &MojeekProvider{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		address: "https://api.mojeek.com",
		tracer:  otel.Tracer("search-gateway"),
		apiKey:  apiKey,
	}
}

// Search sends the query to Mojeek and converts results into Searchbase's
// compact Result shape.
func (p *MojeekProvider) Search(ctx context.Context, req Request) (Results, error) {
	ctx, span := p.tracer.Start(ctx, "MojeekProvider.Search")
	defer span.End()

	span.SetAttributes(
		attribute.String("search.provider", "mojeek"),
	)

	req.Normalize()

	sr, err := p.newRequest(ctx, req)
	if err != nil {
		errMsg := fmt.Errorf("failed to create mojeek search request")
		span.RecordError(errMsg)
		span.SetStatus(codes.Error, "failed to create request")
		return nil, errMsg
	}

	resp, err := p.client.Do(sr)
	if err != nil {
		errMsg := fmt.Errorf("mojeek search provider request failed")
		span.RecordError(errMsg)
		span.SetStatus(codes.Error, "http request failed")
		return nil, errMsg
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Errorf("mojeek search provider returned status %d", resp.StatusCode)
		span.RecordError(errMsg)
		span.SetStatus(codes.Error, "non-200 response")
		return nil, errMsg
	}

	var parsedResp mojeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		errMsg := fmt.Errorf("mojeek search provider returned an invalid response")
		span.RecordError(errMsg)
		span.SetStatus(codes.Error, "failed to decode response")
		return nil, errMsg
	}

	var results Results
	for _, r := range parsedResp.Response.Results {
		results = append(results, Result{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Desc,
		})
	}

	span.SetAttributes(attribute.Int("search.results_count", len(results)))
	return results, nil
}

func (p *MojeekProvider) newRequest(ctx context.Context, req Request) (*http.Request, error) {
	searchUrl, err := url.Parse(strings.TrimRight(p.address, "/") + "/search")
	if err != nil {
		return nil, fmt.Errorf("invalid mojeek address: %w", err)
	}

	q := searchUrl.Query()
	q.Set("api_key", p.apiKey)
	q.Set("q", req.Query)
	q.Set("fmt", "json")

	if req.Limit != 0 {
		q.Set("t", strconv.Itoa(req.Limit))
	}

	switch req.TimeLimit {
	case "d":
		q.Set("since", "day")
	case "m":
		q.Set("since", "month")
	case "y":
		q.Set("since", "year")
	}

	// Safesearch: 0=off, 1=on
	if req.SafeSearch != "" {
		switch req.SafeSearch {
		case "off":
			q.Set("safe", "0")
		case "moderate":
			q.Set("safe", "1")
		case "on", "strict":
			q.Set("safe", "1")
		}
	}

	searchUrl.RawQuery = q.Encode()
	mojeekRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, searchUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return mojeekRequest, nil
}
