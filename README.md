# Go API Gateway

A reverse-proxy API gateway written in Go. It supports YAML-defined routes and clusters, round-robin load balancing, periodic backend health checks, Redis-backed rate limiting, and runtime configuration switching.

## Features

- Prefix-based routes with upstream path rewriting
- Round-robin load balancing across healthy backend servers
- Health checks against `GET /health` every five seconds
- Redis-backed token-bucket rate limiting
- Runtime activation of YAML configuration versions
- Prometheus metrics at `/matrics` on port `2112`

## Requirements

- Go `1.27.2` or later
- Redis
- Reachable backend servers

## Local setup

1. Start Redis on the default local port:

   ```sh
   docker run --name gateway-redis -p 6380:6379 -d redis:7-alpine
   ```

2. From the repository root, start the included mock server:

   ```sh
   go run ./mock-server
   ```

   The mock server listens on `127.0.0.1:9060` and exposes `/health`.

3. Start the gateway in another terminal:

   ```sh
   go run ./cmd/gateway
   ```

The gateway listens on port `8080` by default. The metrics server listens on port `2112`.

Try the configured health route:

```sh
curl -i http://localhost:8080/api/health
```

View gateway metrics:

```sh
curl http://localhost:2112/matrics
```

## Configuration

The gateway loads `configs/gateway/versions/v1.yaml` at startup. Additional `.yaml` or `.yml` files can be activated while the gateway is running.

```yaml
clusters:
  - name: "products_cluster"
    servers:
      - url: "http://127.0.0.1:9001"
      - url: "http://127.0.0.1:9002"

routes:
  - path: "/api/products"
    cluster: "products_cluster"
    target: "/api/products/"

ratelimit:
  max_time_for_token: 6
  number_of_tokens: 10
```

### Routes

- `path` is the request prefix matched by the gateway.
- `cluster` identifies the backend cluster used for the route.
- `target` replaces the matched route prefix.
- A server is eligible to receive traffic only when `GET /health` succeeds.

### Rate limiting

The Redis rate limiter allows `number_of_tokens` requests per client IP and replenishes one token every `max_time_for_token` seconds.

Client IPs are read from `X-Forwarded-For` or `X-Real-Ip`. Configure the proxy in front of the gateway to set these headers so they cannot be spoofed.

## Environment variables

The gateway loads `.env` from the working directory. Environment variables override file values.

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8080` | HTTP gateway port |
| `REDIS_URL` | `redis://localhost:6380/0` | Redis connection URL |

The checked-in `.env` uses `redis://localhost:6379/0`, while the application default uses port `6380`. Export `REDIS_URL=redis://localhost:6380/0` or update `.env` when starting Redis on port `6380`.

Run the gateway from the repository root so the relative configuration directory is available.

## Runtime configuration activation

Activate a configuration file with the gateway control CLI:

```sh
go run ./gatewayctl activate /file_path
```

The CLI sends the absolute file path to the gateway over `/tmp/gateway.sock`. The new configuration replaces the active runtime only after it is loaded and validated.

## Docker Compose

The included Compose file builds the gateway and starts Redis:

```sh
docker compose up --build
```

Redis is published on port `6380`. The Compose file does not currently publish the gateway's `8080` or metrics `2112` ports.

## Tests

Run the Go test suite:

```sh
go test ./...
```

To exercise rate limiting against a running gateway:

```sh
./scripts/ratelimit/rateLimit_test.sh
```
