# golabs-api

## Quick Start

```bash
cp .env.example .env           # configure DB creds and JWT secret
make docker-up                 # start DB + API in Docker
```

---

## Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Compile the binary to `./bin/api` |
| `make run` | Build and run locally |
| `make dev` | Run with hot-reload via `air` |
| `make test` | Run all tests |
| `make lint` | Run `go vet` + `staticcheck` |
| `make docker-up` | Build image and start all containers |
| `make docker-down` | Stop containers |
| `make migrate` | Apply SQL migrations from `deployments/database/init/` |
| `make tidy` | Run `go mod tidy` |

---

## API Versioning

All routes are prefixed with `/api/v1/`.

Legacy `/health` is kept; new probes are at `/healthz/live` and `/healthz/ready`.
