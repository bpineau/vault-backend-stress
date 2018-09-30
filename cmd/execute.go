package cmd

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bpineau/vault-backend-stress/pkg/metrics"
	"github.com/bpineau/vault-backend-stress/pkg/metrics/prometheus"
	"github.com/bpineau/vault-backend-stress/pkg/worker/reader"
)

const (
	appName = "vault-backend-stress"
	dumpFmt = `req/s: [success: {{printf "%.0f" .SuccessRate}}, errors: {{printf "%.0f" .ErrorsRate}}], latencies_ms: [p50: {{printf "%.2f" .P50}}, p95: {{printf "%.2f" .P95}}, p99: {{printf "%.2f" .P99}}]` + "\n"
)

var (
	concurrency int
	timeout     int
	jitter      int
	prefix      string
	address     string
	token       string
	lastDump    *metrics.Point

	// RootCmd is the only cobra command so far
	RootCmd = &cobra.Command{
		Use:    appName,
		Short:  "Simple vault backends stress test and benchmark",
		Long:   "Stress vault backends and measure rates and latencies. You can use traditional vault env vars (VAULT_CACERT, VAULT_SKIP_VERIFY, VAULT_ADDR, VAULT_TOKEN, etc.) to configure vault access",
		PreRun: bindArgs,
		RunE:   runE,
	}
)

func runE(cmd *cobra.Command, args []string) (err error) {
	wg := sync.WaitGroup{}
	workers := make([]reader.SecretReader, concurrency)

	col := new(prometheus.Sink)

	for i := range workers {
		err := workers[i].Init(token, address, prefix, timeout, col)
		if err != nil {
			return err
		}
	}

	for i := range workers {
		go workers[i].Start(&wg)
		time.Sleep(time.Duration(jitter) * time.Millisecond)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)

	t := template.Must(template.New("").Parse(dumpFmt))

	for {
		select {
		case <-time.After(10 * time.Second):
			err = dumpMetrics(t, col)
			if err != nil {
				return err
			}
		case <-sigterm:
			for i := range workers {
				workers[i].Stop()
			}
			wg.Wait()
			return nil
		}
	}
}

func dumpMetrics(t *template.Template, col metrics.Sink) (err error) {
	if lastDump == nil {
		lastDump, err = col.Dump()
		return err
	}

	interval := time.Since(lastDump.Date).Seconds()
	res, err := col.Dump()
	if err != nil {
		return err
	}

	res.SuccessRate = (res.SuccessCount - lastDump.SuccessCount) / interval
	res.ErrorsRate = (res.ErrorsCount - lastDump.ErrorsCount) / interval

	lastDump = res

	return t.Execute(os.Stdout, res)
}

// Execute adds all child commands to the root command and sets their flags.
func Execute() error {
	return RootCmd.Execute()
}

func bindPFlag(key string, cmd string) {
	if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(cmd)); err != nil {
		log.Fatal("Failed to bind cli argument:", err)
	}
}

func init() {
	RootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 100, "concurrency level")
	bindPFlag("concurrency", "concurrency")

	RootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 10, "timeout in seconds")
	bindPFlag("timeout", "timeout")

	RootCmd.PersistentFlags().IntVarP(&jitter, "jitter", "j", 0, "start workers jitter ms appart")
	bindPFlag("jitter", "jitter")

	RootCmd.PersistentFlags().StringVarP(&prefix, "prefix", "p", appName+"/", "keys prefix")
	bindPFlag("prefix", "prefix")

	RootCmd.PersistentFlags().StringVarP(&address, "address", "a", "http://127.0.0.1:8200",
		"vault server address")
	bindPFlag("address", "address")
	if err := viper.BindEnv("address", "VAULT_ADDR"); err != nil {
		log.Fatal("Failed to bind env:", err)
	}

	RootCmd.PersistentFlags().StringVarP(&token, "token", "o", "", "vault token")
	bindPFlag("token", "token")
	if err := viper.BindEnv("token", "VAULT_TOKEN"); err != nil {
		log.Fatal("Failed to bind env:", err)
	}
}

func bindArgs(cmd *cobra.Command, args []string) {
	concurrency = viper.GetInt("concurrency")
	timeout = viper.GetInt("timeout")
	jitter = viper.GetInt("jitter")
	prefix = viper.GetString("prefix")
	address = viper.GetString("address")
	token = viper.GetString("token")
}
