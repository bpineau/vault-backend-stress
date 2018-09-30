package worker

import (
	"sync"

	"github.com/bpineau/vault-backend-stress/pkg/metrics"
)

// Worker actually benchmarks the current vault secrets backend
type Worker interface {
	Init(token string, address string, prefix string, timeout int, collector metrics.Sink) error
	Start(wg *sync.WaitGroup)
	Stop()
}
