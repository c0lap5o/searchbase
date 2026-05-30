package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type braveRequest struct {
	Query           string   `json:"q"`
	Country         string   `json:"country,omitempty"`
	SearchLang      string   `json:"search_lang,omitempty"`
	Count           int      `json:"count,omitempty"`
	Offset          int      `json:"offset,omitempty"`
	SafeSearch      string   `json:"safesearch,omitempty"`
	Freshness       string   `json:"freshness,omitempty"`
	TextDecorations bool     `json:"text_decorations"`
	ResultFilter    []string `json:"result_filter,omitempty"`
	ExtraSnippets   bool     `json:"extra_snippets"`
}

type BraveProvider struct {
	client       *http.Client
	address      string
	webSearchURL string
	tracer       trace.Tracer
	token        string
}

type braveResponse struct {
	Web braveWebResults `json:"web"`
}

type braveWebResults struct {
	Results []braveResult `json:"results"`
}

type braveResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

const (
	braveAPI = "https://api.search.brave.com"
)

var (
	braveCountries = []string{"AR", "AU", "AT", "BE", "BR", "CA", "CL", "DK", "FI", "FR", "DE", "GR", "HK", "IN", "ID", "IT", "JP", "KR", "MY", "MX", "NL", "NZ", "NO", "CN", "PL", "PT", "PH", "RU", "SA", "ZA", "ES", "SE", "CH", "TW", "TR", "GB", "US", "ALL"}

	braveLangs = []string{"ar", "eu", "bn", "bg", "ca", "zh-hans", "zh-hant", "hr", "cs", "da", "nl", "en", "en-gb", "et", "fi", "fr", "gl", "de", "el", "gu", "he", "hi", "hu", "is", "it", "jp", "kn", "ko", "lv", "lt", "ms", "ml", "mr", "nb", "pl", "pt-br", "pt-pt", "pa", "ro", "ru", "sr", "sk", "sl", "es", "sv", "ta", "te", "th", "tr", "uk", "vi"}
)

func NewBraveProvider(token string) *BraveProvider {
	return &BraveProvider{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		address:      braveAPI,
		webSearchURL: fmt.Sprintf("%s/res/v1/web/search", braveAPI),
		tracer:       otel.Tracer("search-gateway"),
		token:        token,
	}
}

func (b *BraveProvider) Search(ctx context.Context, req Request) (Results, error) {
	ctx, span := b.tracer.Start(ctx, "BraveProvider.Search")
	defer span.End()

	span.SetAttributes(
		attribute.String("search.provider", "brave"),
	)

	request, err := b.newRequest(ctx, req, b.webSearchURL)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpResp, err := b.client.Do(request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to execute request")
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	//nolint:errcheck
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, fmt.Sprintf("unexpected status code: %d", httpResp.StatusCode))
		return nil, fmt.Errorf("brave API returned status %d", httpResp.StatusCode)
	}

	var resp braveResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to decode response")
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	results := make([]Result, 0, len(resp.Web.Results))
	for _, r := range resp.Web.Results {
		results = append(results, Result{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Description,
		})
	}

	return results, nil
}

func (b *BraveProvider) newRequest(ctx context.Context, req Request, searchUrl string) (*http.Request, error) {
	var braveReq braveRequest
	if req.Query == "" {
		return nil, fmt.Errorf("query is required")
	}
	braveReq.Query = req.Query
	braveReq.Country, braveReq.SearchLang = parseRegion(req.Region)

	braveReq.Count = min(req.Limit, 20)

	if req.Page > 1 {
		braveReq.Offset = min(req.Page-1, 9)
	}
	braveReq.SafeSearch = req.SafeSearch
	if braveReq.SafeSearch == "on" {
		braveReq.SafeSearch = "strict"
	}
	if req.TimeLimit != "" {
		switch strings.ToLower(req.TimeLimit) {
		case "d":
			braveReq.Freshness = "pd"
		case "w":
			braveReq.Freshness = "pw"
		case "m":
			braveReq.Freshness = "pm"
		case "y":
			braveReq.Freshness = "py"
		default:
			braveReq.Freshness = ""
		}
	}
	braveReq.TextDecorations = false
	braveReq.ResultFilter = []string{"web"}
	braveReq.ExtraSnippets = false

	body, err := json.Marshal(braveReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal brave request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		searchUrl,
		strings.NewReader(string(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("x-subscription-token", b.token)

	return httpReq, nil
}

func parseRegion(region string) (string, string) {
	if region == "" {
		return "", ""

	}
	parts := strings.Split(region, "-")
	var country, language string
	if len(parts) != 2 {
		return "", ""
	}

	regionCountry := strings.ToUpper(parts[0])
	if slices.Contains(braveCountries, regionCountry) {
		country = regionCountry
	}

	regionLanguage := strings.ToLower(parts[1])
	if slices.Contains(braveLangs, regionLanguage) {
		language = regionLanguage
	}

	return country, language
}
