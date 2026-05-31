package settings

import "fmt"

// SearchProvider stores startup configuration for the selected search backend.
type SearchProvider struct {
	name    string
	address string
	token   string
}

func (sp *SearchProvider) Name() string    { return sp.name }
func (sp *SearchProvider) Address() string { return sp.address }
func (sp *SearchProvider) Token() string   { return sp.token }
func (sp *SearchProvider) validate() error {
	switch sp.name {
	case "ddgs":
		if sp.address == "" {
			return fmt.Errorf("provider address not set: searchbase provider ddgs requires SEARCHBASE_DDGS_PROVIDER_ADDRESS")
		}
		return nil
	case "searxng":
		if sp.address == "" {
			return fmt.Errorf("provider address not set: searchbase provider searxng requires SEARCHBASE_SEARXNG_PROVIDER_ADDRESS")
		}
		return nil
	case "searchbase_ddg":
		return nil
	case "brave":
		if sp.token == "" {
			return fmt.Errorf("provider API Token not set: searchbase provider brave requires SEARCHBASE_BRAVE_API_TOKEN")
		}
		return nil
	}

	return fmt.Errorf("unsupported search provider: %s", sp.name)
}
