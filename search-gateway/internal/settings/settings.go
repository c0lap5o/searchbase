package settings

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Settings struct {
	port               string
	crawlWorkerAddress string
	searchProvider     *SearchProvider
	ginMode            string
	address            string
	logLevel           string
	otel               *Otel
}

func (s *Settings) Port() string                    { return s.port }
func (s *Settings) CrawlWorkerAddress() string      { return s.crawlWorkerAddress }
func (s *Settings) SearchProvider() *SearchProvider { return s.searchProvider }
func (s *Settings) GinMode() string                 { return s.ginMode }
func (s *Settings) Address() string                 { return s.address }
func (s *Settings) LogLevel() string                { return s.logLevel }
func (s *Settings) Otel() *Otel                     { return s.otel }

// NewSettings initializes a new Settings instance with default values from the environment
// searchbase address, should be left empty if running behind a reverse proxy for routing to use relative paths
func NewSettings() (*Settings, error) {
	v := viper.New()
	v.SetEnvPrefix("SEARCHBASE")
	v.AutomaticEnv()

	v.SetDefault("PORT", "8080")
	v.SetDefault("CRAWL_WORKER_ADDRESS", "http://localhost:8000")
	v.SetDefault("SEARCH_PROVIDER", "searchbase_ddg")
	v.SetDefault("LOG_LEVEL", "error")
	v.SetDefault("ENVIRONMENT", "production")

	providerName := v.GetString("SEARCH_PROVIDER")
	var providerAddress string
	var token string
	switch providerName {
	case "ddgs":
		providerAddress = v.GetString("DDGS_PROVIDER_ADDRESS")
	case "searxng":
		providerAddress = v.GetString("SEARXNG_PROVIDER_ADDRESS")
	case "brave":
		token = v.GetString("BRAVE_API_TOKEN")
	}

	s := &Settings{
		port:               v.GetString("PORT"),
		crawlWorkerAddress: v.GetString("CRAWL_WORKER_ADDRESS"),
		address:            v.GetString("ADDRESS"),
		logLevel:           strings.ToLower(v.GetString("LOG_LEVEL")),
		ginMode:            gin.ReleaseMode,
		searchProvider: &SearchProvider{
			name:    providerName,
			address: providerAddress,
			token:   token,
		},
		otel: &Otel{
			tracing: &Tracing{
				enabled:  v.GetBool("TRACING_ENABLED"),
				endpoint: v.GetString("TRACING_ENDPOINT"),
			},
		},
	}

	if strings.ToLower(v.GetString("ENVIRONMENT")) == "dev" {
		s.ginMode = gin.DebugMode
	}

	if err := s.searchProvider.validate(); err != nil {
		return nil, err
	}

	if err := s.otel.tracing.validate(); err != nil {
		return nil, err
	}

	return s, nil
}
