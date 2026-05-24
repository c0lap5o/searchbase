---
title: "API Reference"
weight: 30
---

You can also build your own scripts or use LangChain / LlamaIndex by hitting the standard REST API.

**Note:** A fully interactive Swagger UI is available at `http://localhost:8080/docs` to test and explore the API.

## Web Search
**Endpoint:** `POST /api/v1/search`

**Request:**
```json
{
  "query": "kubernetes 1.30 release notes",
  "limit": 5,
  "provider": "auto",
  "engine": "auto",
  "region": "wt-wt",
  "timelimit": "d",
  "safesearch": "moderate",
  "page": 1
}
```

**Additional Optional Parameters:**
- `provider`: Ignored in this version. The search provider is configured at the server level via environment variables.
- `engine`: The search engine to query (e.g., `"auto"`, `"google"`, `"duckduckgo"`). Defaults to `"auto"` (searches all engines and deduplicates). **Note:** Only applies if the server is configured to use the `ddgs` or `searxng` provider.
- `region`: The region to search in (e.g., `"wt-wt"`, `"us-en"`).
- `timelimit`: Time limit for the search (`"d"`=day, `"w"`=week, `"m"`=month, `"y"`=year). Leave empty for no limit.
- `safesearch`: Safe search filtering (`"on"`, `"moderate"`, `"off"`). Defaults to `"moderate"`.
- `page`: The page number of results to fetch. Defaults to `1`.

**Response:**
```json
[
  {
    "title": "Kubernetes 1.30: ...",
    "url": "https://kubernetes.io/...",
    "snippet": "Short description from the search engine..."
  }
]
```

## Fetch URL
**Endpoint:** `POST /api/v1/fetch`

**Request:**
```json
{
  "url": "https://kubernetes.io/blog/...",
  "js_render": false
}
```
*(Set `js_render: true` if you need the worker to execute JavaScript for Single Page Applications, though it is slower).*

**Response:**
```json
{
  "title": "",
  "url": "https://kubernetes.io/blog/...",
  "snippet": "",
  "markdown": "## Kubernetes 1.30\n\nThe new release features..."
}
```
