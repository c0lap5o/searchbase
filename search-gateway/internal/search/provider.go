package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/coolapso/searchbase/search-gateway/internal/settings"
)

type Request struct {
	Query string `json:"query" binding:"required" example:"who is the president" description:"The search query to look for"`
	Limit int    `json:"limit" example:"3" description:"Number of search results to return (default 3, max 10)"`
	//TODO: this ignored right now, Add support for multiple search providers
	Provider   string `json:"provider" example:"auto" description:"Ignored in this version. The search provider is configured at the server level."`
	Engine     string `json:"engine" example:"auto" description:"The search engine to query (e.g., 'google', 'duckduckgo'). Defaults to 'auto' (all engines). Note: Only applies if the server uses the 'ddgs' or 'searxng' provider."`
	Region     string `json:"region" example:"wt-wt" description:"The region to search in (e.g., 'wt-wt', 'us-en')"`
	TimeLimit  string `json:"timelimit" example:"d" description:"Time limit for the search ('d' for day, 'w' for week, 'm' for month, 'y' for year). Leave empty for no limit."`
	SafeSearch string `json:"safesearch" example:"moderate" description:"Safe search filtering ('on', 'moderate', 'off')"`
	Page       int    `json:"page" example:"1" description:"The page number of results to fetch"`
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
	if r.Limit <= 0 {
		r.Limit = 5
	} else if r.Limit > 10 {
		r.Limit = 10
	}
	if r.Provider == "" {
		r.Provider = "auto"
	}
	if r.Engine == "" {
		r.Engine = "auto"
	}
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
	}

	return nil, fmt.Errorf("unsupported provider: %s", provider)
}
