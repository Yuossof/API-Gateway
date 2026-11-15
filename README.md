# API Gateway

A Go reverse-proxy gateway that routes requests to configured backend clusters. It checks backend health, balances requests across healthy servers, rate-limits clients with Redis, and supports switching configuration versions while running.

## Features

- Prefix-based routes with configurable upstream path rewriting
- Round-robin load balancing across healthy servers
- Backend health checks (`GET /health` on each configured server, every five seconds)
- Redis-backed token-bucket rate limiting
- Runtime configuration activation over a local Unix socket
<<<<<<< HEAD
- Metrics endpoint at `/matrics` on port `2112`
=======
- Prometheus metrics at `/matrics` on port `2112`
>>>>>>> 6f82c121fd53d7d42d588dfbe6d6c0b7fdc3687f

## Requirements

- Go `1.27.2` or later
- Redis, reachable at the configured `REDIS_URL`
- Backend servers listed in the active gateway configuration

The default Redis URL is `redis://localhost:6380/0`. For a local Redis instance on that port:

```sh
docker run --name gateway-redis -p 6380:6379 -d redis:7-alpine
```

## Run locally

1. Start Redis.
2. Review `configs/gateway/versions/v1.yaml` and update its backend URLs for your environment. The checked-in products cluster uses example network addresses; the health cluster points to `127.0.0.1:9060`.
3. Ensure every backend responds successfully to `GET /health`. The gateway marks a server unhealthy if this check fails.
4. Start the optional local mock backend in another terminal if you want to use the included health route:

   ```sh
   go run ./mock-server
   ```

5. Start the gateway from the repository root:

   ```sh
   go run ./cmd/gateway
   ```

The gateway listens on port `8080` by default. Try the configured health route:

```sh
curl -i http://localhost:8080/api/health
```

The mock server listens on port `9060` and responds to `/health`. To build the gateway binary instead:

```sh
go build -o gateway ./cmd/gateway
./gateway
```

## Environment

The gateway loads a `.env` file from its working directory when present. Environment variables take precedence over defaults:

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8080` | HTTP port for gateway traffic |
| `REDIS_URL` | `redis://localhost:6380/0` | Redis connection URL |

Run the gateway from the repository root so its relative configuration directory (`configs/gateway/versions/`) is found.

## Gateway configuration

Configuration versions are YAML files in `configs/gateway/versions/`. The gateway activates `v1.yaml` at startup. A version contains backend clusters, routes, and rate-limit settings; for example:

```yaml
clusters:
  - name: products
    servers:
      - url: "http://127.0.0.1:9001"
      - url: "http://127.0.0.1:9002"

routes:
  - path: "/api/products"
    cluster: products
    target: "/products"

ratelimit:
  max_time_for_token: 6
  number_of_tokens: 10
```

Each route matches its configured path and subpaths. The matched prefix is replaced by `target` while the remaining path suffix is preserved. Backend URLs should be reachable by the gateway and expose a successful `GET /health` endpoint.

The rate limiter allows an initial burst of `number_of_tokens` requests per client IP and replenishes one token every `max_time_for_token` seconds. It requires a working Redis connection. When running behind a proxy, the gateway uses `X-Forwarded-For` (first address) or `X-Real-Ip` to identify the client; configure the proxy to overwrite these headers with trusted values.

### Activate a configuration version

While the gateway is running, activate another YAML file from the versions directory using the included CLI:

```sh
go run ./gatewayctl activate v2.yaml
```

Activation sends a request to `/tmp/gateway.sock`, validates and builds the new runtime, then switches traffic to it. Only version file names ending in `.yaml` are accepted. `v1.yaml` and `v2.yaml` currently contain the same example configuration.

<<<<<<< HEAD
## Metrics

The gateway exposes metrics at:
=======
## Metrics and monitoring

The gateway serves Prometheus metrics at:
>>>>>>> 6f82c121fd53d7d42d588dfbe6d6c0b7fdc3687f

```sh
curl http://localhost:2112/matrics
```

<<<<<<< HEAD
## Docker Compose

Build and start the gateway with Redis:

```sh
docker compose up --build
```

The gateway is available on port `8080`, its metrics endpoint on port `2112`, and Redis on port `6380`.
=======
The repository includes Prometheus and Grafana services for local monitoring:

```sh
docker compose up -d
```

Prometheus is exposed on port `9090` and Grafana on port `3000`. The included Grafana defaults are intended for local development only; change them before using this setup outside a trusted local environment.
>>>>>>> 6f82c121fd53d7d42d588dfbe6d6c0b7fdc3687f

## Tests

Run the Go test suite from the repository root:

```sh
go test ./...
```

An additional rate-limit request script is available at `scripts/ratelimit/rateLimit_test.sh`; it sends requests to a running gateway at `http://localhost:8080/api/health`.
# API-Gateway
