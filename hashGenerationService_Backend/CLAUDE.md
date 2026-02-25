# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server
go run ./cmd/server/main.go

# Build
go build ./...

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/service/...
go test ./internal/store/...

# Run a single test
go test ./internal/service/... -run TestGenerateHash

# Vet
go vet ./...
```

Server starts on `:8080`. No environment variables or config files are required.

## Architecture

Three-layer architecture wired together in `cmd/server/main.go`:

```
HTTP request
    → middleware.CORS         (internal/middleware/cors.go)
    → handler.Handler         (internal/handler/handler.go)
    → service.Service         (internal/service/service.go)
    → store.Store interface   (internal/store/store.go)
        └─ InMemoryStore      (internal/store/inmemory.go)
```

**Store** (`internal/store/`) — `Store` is an interface with `Save`, `Exists`, `Get`. `InMemoryStore` is the only implementation; it holds a `map[string]string` (hash → input) guarded by `sync.RWMutex`. Swapping in a persistent store only requires satisfying the interface.

**Service** (`internal/service/service.go`) — owns all business logic: validates that input is alphanumeric-only, calls `utils.GenerateHash` up to `maxRetries` (5) times until a hash that doesn't already exist in the store is found, then saves and returns it. Sentinel errors (`ErrInvalidInput`, `ErrMaxRetriesExceeded`, `ErrHashNotFound`) are defined here and checked by the handler via `errors.Is`.

**Utils** (`internal/utils/hash.go`) — generates a 10-character base-62 hash by SHA-256(input + 16 random bytes). The random salt means the same input produces a different hash each call, which is what enables collision retries.

**Handler** (`internal/handler/handler.go`) — uses the standard `net/http` mux (no framework). Method guards are manual (`r.Method != ...`) because the mux routes by path only. URL parameter extraction for `GET /hash/{hash}` is done with `strings.TrimPrefix`.

## API

| Method | Path | Success | Body |
|--------|------|---------|------|
| `POST` | `/hash` | `201` | `{"input":"...","hash":"..."}` |
| `GET` | `/hash/{hash}` | `200` | `{"hash":"...","input":"..."}` |

Error responses always use `{"error":"..."}`. Status codes: `400` bad input, `404` not found, `409` max retries exceeded, `500` internal error.

## Module

Module name is `hashGenerationService` (see `go.mod`). All internal imports use this as the root prefix, e.g. `hashGenerationService/internal/service`. No external dependencies — stdlib only.
