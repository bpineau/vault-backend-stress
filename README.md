# vault-backend-stress

Stress Hashicorp Vault backends, measuring reads rates and latencies.

## Usage

`vault-backend-stress` will write (then mostly read, and delete) random keys names (to ensure a proper distribution, ie. on s3 and gcs) under a path starting with the provided `prefix`.

You can use traditional Vault env vars (VAULT_CACERT, VAULT_SKIP_VERIFY, VAULT_CAPATH, VAULT_ADDR, VAULT_TOKEN, etc.) to fine tune connections to Vault.

```
Usage:
  vault-backend-stress [flags]

Flags:
  -a, --address string    vault server
  -c, --concurrency int   concurrency level (default 100)
  -h, --help              help for vault-backend-stress
  -j, --jitter int        start workers jitter ms appart
  -p, --prefix string     keys prefix (default "vault-backend-stress/")
  -t, --timeout int       timeout in seconds (default 10)
  -o, --token string      vault token
```

## Example

```
env VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=XXX \
  ./vault-backend-stress -c 1000 -p secrets/stresstest

req/s: [success: 4356, errors: 0], latencies_ms: [p50: 98.10, p95: 1047.13, p99: 1279.47]
req/s: [success: 3607, errors: 0], latencies_ms: [p50: 142.75, p95: 441.56, p99: 1279.47]
req/s: [success: 3031, errors: 0], latencies_ms: [p50: 167.77, p95: 472.44, p99: 1245.48]
req/s: [success: 2951, errors: 0], latencies_ms: [p50: 194.86, p95: 498.07, p99: 1226.55]
req/s: [success: 3062, errors: 0], latencies_ms: [p50: 225.86, p95: 492.64, p99: 1210.45]
req/s: [success: 5789, errors: 0], latencies_ms: [p50: 194.86, p95: 507.00, p99: 1226.55]
req/s: [success: 4562, errors: 0], latencies_ms: [p50: 167.77, p95: 507.00, p99: 1245.48]
```
