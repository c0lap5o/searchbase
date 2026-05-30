package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/coolapso/searchbase/search-gateway/internal/settings"
)

type Request struct {
	Query      string `json:"query" binding:"required" example:"who is the president" description:"The search query to look for"`
	Limit      int    `json:"limit,omitempty" example:"3" description:"Optional upper bound for returned search results. Omit or set to 0 to use the provider or search engine default."`
	Engine     string `json:"engine,omitempty" example:"auto" description:"The search engine to query (e.g., 'google', 'duckduckgo'). Defaults to 'auto' (all engines). Note: Only applies if the server uses the 'ddgs' or 'searxng' provider."`
	Region     string `json:"region,omitempty" example:"wt-wt" description:"The region to search in (e.g., 'wt-wt', 'us-en')"`
	TimeLimit  string `json:"timelimit,omitempty" example:"d" description:"Time limit for the search ('d' for day, 'w' for week, 'm' for month, 'y' for year). Leave empty for no limit."`
	SafeSearch string `json:"safesearch,omitempty" example:"moderate" description:"Safe search filtering ('on', 'moderate', 'off')"`
	Page       int    `json:"page,omitempty" example:"1" description:"The page number of results to fetch"`
}

// SearchResult represents a single search hit.
type Result struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type Results []Result

// SearchProvider defines the interface for different search engines.
type SearchProvider interface {
	Search(ctx context.Context, request Request) (Results, error)
}

func (r *Request) Normalize() {
	if r.SafeSearch == "" {
		r.SafeSearch = "moderate"
	}

	if r.TimeLimit != "" {
		switch strings.ToLower(r.TimeLimit) {
		case "d", "day":
			r.TimeLimit = "d"
		case "w", "week":
			r.TimeLimit = "w"
		case "m", "month":
			r.TimeLimit = "m"
		case "y", "year":
			r.TimeLimit = "y"
		default:
			r.TimeLimit = ""
		}
	}
}

func NewProvider(provider *settings.SearchProvider) (SearchProvider, error) {
	switch provider.Name() {
	case "ddgs":
		return NewDDGSProvider(provider.Address()), nil
	case "searxng":
		return NewSearxngProvider(provider.Address()), nil
	case "searchbase_ddg":
		return NewDuckDuckGoProvider(), nil
	case "brave":
		return NewBraveProvider(provider.Token()), nil
	}

	return nil, fmt.Errorf("unsupported provider: %s", provider.Name())
}
