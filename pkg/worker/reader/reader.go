package reader

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bpineau/vault-backend-stress/pkg/metrics"
	vault "github.com/hashicorp/vault/api"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// SecretReader stress test secrets reads
type SecretReader struct {
	client    *vault.Client
	key       string
	ctx       context.Context
	cancel    context.CancelFunc
	collector metrics.Sink
}

// Init build a reader
func (r *SecretReader) Init(token string, address string, prefix string, timeout int,
	collector metrics.Sink) error {

	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.collector = collector

	client, err := vault.NewClient(&vault.Config{
		Address:    address,
		Timeout:    time.Duration(timeout) * time.Second,
		MaxRetries: 0,
	})

	if err != nil {
		return err
	}

	r.client = client
	r.client.SetToken(token)

	secretData := map[string]interface{}{
		"value": "world",
		"foo":   "bar",
		"age":   "-1",
	}

	r.key = fmt.Sprintf("%s/%s", prefix, randString(64))
	_, err = r.client.Logical().Write(r.key, secretData)

	if err != nil {
		return err
	}

	return nil
}

// Start launch reader
func (r *SecretReader) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-r.ctx.Done():
			_, _ = r.client.Logical().Delete(r.key)
			return
		default:
		}

		before := time.Now()
		_, err := r.client.Logical().Read(r.key)

		if err != nil {
			r.collector.Observe(metrics.Error, time.Since(before))
		} else {
			r.collector.Observe(metrics.Success, time.Since(before))
		}
	}
}

// Stop cancel a reader
func (r *SecretReader) Stop() {
	r.cancel()
}

func randString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
