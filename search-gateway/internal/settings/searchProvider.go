package settings

// SearchProvider stores startup configuration for the selected search backend.
type SearchProvider struct {
	name    string
	address string
	token   string
}

func (sp *SearchProvider) Name() string    { return sp.name }
func (sp *SearchProvider) Address() string { return sp.address }
func (sp *SearchProvider) Token() string   { return sp.token }
