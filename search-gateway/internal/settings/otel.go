package settings

import "fmt"

type Otel struct {
	tracing *Tracing
}

type Tracing struct {
	enabled  bool
	endpoint string
}

func (o *Otel) Tracing() *Tracing { return o.tracing }

func (t *Tracing) Enabled() bool    { return t.enabled }
func (t *Tracing) Endpoint() string { return t.endpoint }
func (t *Tracing) validate() error {
	if t.enabled && t.endpoint == "" {
		return fmt.Errorf("tracing endpoint not set: SEARCHBASE_TRACING_ENDPOINT")
	}

	return nil
}
