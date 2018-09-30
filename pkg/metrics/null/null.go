package null

import (
	"time"

	"github.com/bpineau/vault-backend-stress/pkg/metrics"
)

// Sink just discard metrics
type Sink struct{}

// Observe receive an observation
func (n *Sink) Observe(status metrics.Status, duration time.Duration) {
}

// Dump return metrics state
func (n *Sink) Dump() (*metrics.Point, error) {
	return new(metrics.Point), nil
}
