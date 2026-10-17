# AI Gateway

A lightweight AI gateway written in Go that exposes an HTTP API for chat completions and routes requests to configurable LLM providers. The project includes API-key authentication, per-user rate limiting, token-balance tracking, Redis response caching, and provider retry/fallback logic.

> **Status:** Work in progress. Review configuration and security-sensitive behavior before using this project in production.

## Features

- **HTTP API** for chat completions.
- **API-key authentication** using `Authorization: Bearer <api-key>`.
- **Multiple LLM providers:** OpenRouter and a local Ollama provider.
- **Retry and fallback:** retries failed provider calls and can fall back to the default provider.
- **Redis caching** for chat responses, with a configurable TTL.
- **Per-user rate limiting** with a configurable request limit and time window.
- **Token accounting:** estimates token usage when provider usage data is unavailable and deducts usage from a user's balance.
- **PostgreSQL persistence** for user records and token balances.
- **Graceful HTTP server shutdown** on termination signals.
- **Environment-variable configuration** with YAML configuration support.

## Tech stack

- Go
- `net/http`
- PostgreSQL
- Redis
- OpenRouter API
- Ollama

## Architecture

The code is organized by responsibility:

```text
.
├── adapter/redis/          # Redis client adapter
├── config/                 # Configuration loading and defaults
├── delivery/httpserver/    # HTTP server, handlers, and middleware
├── entity/                 # Domain entities
├── param/                  # HTTP request/response structures
├── pkg/
│   ├── cache/               # Cache-key generation
│   ├── estimatetoken/       # Approximate token estimation
│   └── ratelimiter/         # Per-user rate limiter
├── provider/
│   ├── ollama/              # Local Ollama provider
│   └── openrouter/          # OpenRouter provider
├── repository/
│   ├── migrator/            # Database migration support
│   ├── postgres/            # PostgreSQL repository
│   └── redis/               # Redis cache repository
├── service/
│   ├── aiservice/           # Provider selection, retries, caching
│   └── userservice/         # User creation, API keys, token balances
├── config.yml
├── docker-compose.yml
├── go.mod
└── main.go
```

## Requirements

- Go **1.26.8** as specified in `go.mod` (or a compatible toolchain)
- Docker and Docker Compose
- An OpenRouter API key if you want to use the OpenRouter provider
- Ollama and the configured local model if you want to use the local provider

## Getting started

### 1. Clone the repository

```bash
git clone https://github.com/yonesbeheshti/ai-gateway.git
cd ai-gateway
```

### 2. Configure environment variables

The application loads `config.yml` and then applies environment variables prefixed with `AI_GATEWAY_`.

Set your OpenRouter key without putting it directly into source control:

```bash
export AI_GATEWAY_OPENROUTER_API_KEY="your-openrouter-api-key"
```

The nested form is also supported:

```bash
export AI_GATEWAY_OPENROUTER__API_KEY="your-openrouter-api-key"
```

Use your own PostgreSQL and Redis settings if they differ from the defaults in `config.yml`. For example:

```bash
export AI_GATEWAY_DBCONFIG_HOST="localhost"
export AI_GATEWAY_DBCONFIG_PORT="5432"
export AI_GATEWAY_DBCONFIG_DB_NAME="ai_gateway"
export AI_GATEWAY_DBCONFIG_USERNAME="postgres"
export AI_GATEWAY_DBCONFIG_PASSWORD="postgres"

export AI_GATEWAY_REDIS_HOST="localhost"
export AI_GATEWAY_REDIS_PORT="6379"
export AI_GATEWAY_REDIS_DB="0"
```

The default HTTP port is `8080`. The default rate-limit configuration is `10` requests per `60s` window.

> **Security:** Never commit `.env` files, provider keys, database passwords, or real user API keys. If a secret has already been committed or shared, revoke/rotate it. Use environment variables or a secret manager for deployments.

### 3. Start PostgreSQL and Redis

The repository includes a Compose file for local dependencies:

```bash
docker compose up -d
```

Check that the services are running:

```bash
docker compose ps
```

The Compose file configures PostgreSQL on port `5432` and Redis on port `6379`. Confirm the database name and health-check settings in `docker-compose.yml` match your configuration before relying on the health check.

### 4. Configure the application

Review `config.yml` and update values for your environment. The current default configuration includes:

```yaml
httpserver:
  address: ""
  port: 8080

dbconfig:
  host: localhost
  port: 5432
  db_name: ai_gateway
  username: postgres
  password: postgres
  sslmode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

ratelimiter:
  limit: 10
  window: 60s
```

The AI service cache TTL defaults to one hour. The provider API key should be supplied through `AI_GATEWAY_OPENROUTER_API_KEY` rather than committed to `config.yml`.

### 5. Run the application

```bash
go mod download
go run .
```

The server listens on the configured address and port. With the default settings, the base URL is:

```text
http://localhost:8080
```

Ensure the PostgreSQL schema/migrations have been applied before using endpoints that access user data.

## API

### Health check

```http
GET /health_check
```

Use this endpoint to check whether the HTTP server is reachable.

### Login / API-key creation

```http
POST /api/login
Content-Type: application/json
```

Request:

```json
{
  "email": "developer@example.com"
}
```

The current implementation creates a user and API key for a new email address, or returns the existing user's key and balance. This is a development-oriented flow, **not a production authentication system**: it does not verify email ownership.

### Chat completions

```http
POST /api/chat/completions
Authorization: Bearer sk-gw-YOUR_API_KEY
Content-Type: application/json
```

Request:

```json
{
  "model": "openrouter",
  "messages": [
    {
      "role": "user",
      "content": "Explain rate limiting in simple terms."
    }
  ],
  "no_cache": false
}
```

Example with `curl`:

```bash
curl -X POST "http://localhost:8080/api/chat/completions" \
  -H "Authorization: Bearer sk-gw-YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "openrouter",
    "messages": [
      {
        "role": "user",
        "content": "Explain rate limiting in simple terms."
      }
    ],
    "no_cache": false
  }'
```

The response contains the generated message and token usage. The exact provider model is currently selected in the provider implementation, so the request's `model` field selects a configured gateway provider rather than necessarily choosing an arbitrary upstream model.

**Provider selection note:** the current code registers `openrouter` and `ollama` as provider keys and uses OpenRouter as the default provider. The OpenRouter and Ollama implementations currently contain hard-coded upstream model names; make these configurable if clients need to select upstream models.

## Configuration reference

| Setting | Purpose | Default / example |
|---|---|---|
| `httpserver.address` | Address to bind the HTTP server to | `""` |
| `httpserver.port` | HTTP port | `8080` |
| `dbconfig.host` | PostgreSQL host | `localhost` |
| `dbconfig.port` | PostgreSQL port | `5432` |
| `dbconfig.db_name` | PostgreSQL database name | `ai_gateway` |
| `dbconfig.username` | PostgreSQL username | `postgres` |
| `dbconfig.password` | PostgreSQL password | `postgres` |
| `dbconfig.sslmode` | PostgreSQL SSL mode | `disable` |
| `redis.host` | Redis host | `localhost` |
| `redis.port` | Redis port | `6379` |
| `redis.db` | Redis database index | `0` |
| `ratelimiter.limit` | Requests allowed per window per user | `10` |
| `ratelimiter.window` | Rate-limit window | `60s` |
| `aiservice.cache_ttl` | Cache entry lifetime | `1h` |

## Development

Run the test suite:

```bash
go test ./...
```

Format Go code:

```bash
gofmt -w .
```

## Limitations and security considerations

This is a learning/development project and should not be considered production-ready without further work. In particular:

- **Authentication:** the login endpoint currently uses an email address to retrieve or create an API key. Add email verification or another identity-verving mechanism before deployment.
- **Rate limiting:** the current limiter stores counters in process memory. Limits reset when the process restarts and are not shared across multiple gateway instances.
- **Token estimation:** the fallback estimator uses a rough character-based approximation, not a model-specific tokenizer. Billing should use authoritative provider usage where available.
- **Caching:** review cache-key composition and cache behavior carefully before enabling it for production traffic. Do not cache sensitive or user-specific responses without an appropriate isolation policy.
- **Secrets:** rotate any credentials that have been exposed, and keep keys/passwords out of Git.
- **Observability:** add structured logs, metrics, tracing, and alerting before production use.
- **Validation and error handling:** validate request fields and upstream responses, and avoid returning internal error details to clients.
- **Database migrations:** confirm migrations are run as part of setup/deployment.

## Roadmap ideas

- Add tests for HTTP handlers, provider failures, rate limiting, and cache behavior.
- Make provider models and provider selection configurable.
- Move rate-limit state to Redis for multi-instance deployments.
- Improve token counting and usage reporting.
- Add request validation, structured logging, metrics, and tracing.
- Harden authentication and API-key management.
