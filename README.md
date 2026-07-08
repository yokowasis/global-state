# bima-state

A global key-value state server backed by Redis, with a TypeScript client library.

## Stack

- **Server:** Go (`net/http` + `go-redis/v9`)
- **Storage:** Redis 7 (persistent via Docker volume)
- **Client:** TypeScript (compiled to ES6)
- **Containerization:** Docker + Docker Compose

## Architecture

```
Client (TS/JS)  ──POST──►  Go Server  ──►  Redis
                   :8080              :6379
```

The client sends JSON payloads to `/?_=/get` or `/?_=/set`. The Go server reads/writes from Redis and returns JSON responses.

## Getting Started

```bash
# Start services (app + redis)
docker compose up -d

# Stop
docker compose down
```

Server runs on `http://localhost:8080`.

## API

All requests use `POST` with `Content-Type: application/json`.

### Set

```
POST /?_=/set
{"key": "foo", "value": "bar"}
```

```json
{"success": true, "key": "foo", "value": "bar"}
```

### Get

```
POST /?_=/get
{"key": "foo"}
```

```json
{"key": "foo", "value": "bar"}
```

Returns `{"key": "foo", "value": null}` if the key does not exist.

### Errors

| Condition | Status | Body |
|---|---|---|
| Invalid JSON | `400` | `{"error": "invalid JSON"}` |
| Unknown action | `404` | `{"error": "unknown action: ..."}` |
| Wrong method | `405` | `{"error": "method not allowed"}` |

## Client Library

```typescript
import { getState, setState } from "bima-state";

await setState("user", { name: "Alice", age: 30 });
const result = await getState("user");
// { key: "user", value: { name: "Alice", age: 30 } }
```

The client currently points to a Cloudflare Worker URL. To use your local server, update the URL in `src/index.ts` to `http://localhost:8080`.

## Persistence

Redis writes to `/data/dump.rdb` inside the container. A Docker named volume (`redis-data`) is mounted at `/data`, so data survives container restarts.

Default save triggers (configurable via `redis.conf`):

| Writes | Window |
|---|---|
| ≥ 1 | 900s (15 min) |
| ≥ 10 | 300s (5 min) |
| ≥ 10000 | 60s (1 min) |

## Configuration

| Variable | Default | Description |
|---|---|---|
| `REDIS_ADDR` | `localhost:6379` | Redis host:port |
| `PORT` | `8080` | Server listen port |

## Project Structure

```
├── server/
│   ├── main.go
│   ├── go.mod
│   └── go.sum
├── src/
│   └── index.ts
├── dist/
│   ├── index.js
│   └── index.d.ts
├── Dockerfile
├── docker-compose.yml
└── package.json
```
