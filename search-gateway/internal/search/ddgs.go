package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type DDGSProvider struct {
	client  *http.Client
	address string
	tracer  trace.Tracer
}

type ddgsSearchRequest struct {
	Query      string `json:"query"`
	Region     string `json:"region,omitempty"`
	Safesearch string `json:"safesearch,omitempty"`
	TimeLimit  string `json:"timelimit,omitempty"`
	MaxResults int    `json:"max_results,omitempty"`
	Page       int    `json:"page,omitempty"`
	Backend    string `json:"backend,omitempty"`
}

type ddgsResult struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Body  string `json:"body"`
}

type ddgsResponse struct {
	Results []ddgsResult `json:"results"`
}

func NewDDGSProvider(address string) *DDGSProvider {
	return &DDGSProvider{
		// TODO: Add any custom proxy/timeout configs here later
		client:  &http.Client{},
		address: address,
		tracer:  otel.Tracer("search-gateway"),
	}
}

func (p *DDGSProvider) Search(ctx context.Context, req Request) (Results, error) {
	ctx, span := p.tracer.Start(ctx, "DDGSProvider.Search")
	defer span.End()

	span.SetAttributes(
		attribute.String("search.provider", "ddgs"),
	)

	searchUrl := fmt.Sprintf("%s/search/text", p.address)

	ddgsReq := ddgsSearchRequest{
		Query:      req.Query,
		Region:     req.Region,
		Safesearch: req.SafeSearch,
		TimeLimit:  req.TimeLimit,
		MaxResults: req.Limit,
		Page:       req.Page,
		Backend:    req.Engine,
	}

	jsonReq, err := json.Marshal(ddgsReq)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to marshal request")
		return nil, fmt.Errorf("failed to marshal DDGS request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, searchUrl, bytes.NewBuffer(jsonReq))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create http request")
		return nil, fmt.Errorf("failed to create ddgs request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "http request failed")
		return nil, fmt.Errorf("failed to execute ddgs request: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code from provider: %d", resp.StatusCode)
		span.RecordError(err)
		span.SetStatus(codes.Error, "non-200 response")

		return nil, err
	}

	var ddgsResults ddgsResponse
	if err := json.NewDecoder(resp.Body).Decode(&ddgsResults); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to decode response")
		return nil, fmt.Errorf("failed to decode provider response: %w", err)
	}

	results := make(Results, len(ddgsResults.Results))
	for i, r := range ddgsResults.Results {
		results[i] = Result{
			Title:   r.Title,
			URL:     r.Href,
			Snippet: r.Body,
		}
	}

	return results, nil
}
