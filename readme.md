# Behavix AI

Production-grade SaaS analytics backend for ingesting and storing product usage events and serving behavioral insights. Built in **Go** with **Gin** and **PostgreSQL**, following **clean architecture**.

---

## Tech stack

| Layer        | Technology        |
|-------------|--------------------|
| Language    | Go 1.21+           |
| HTTP        | Gin                |
| Database    | PostgreSQL (pgx/v5)|
| Config      | Environment (.env) |
| Logging     | Zap                |

---

## Architecture

The project uses **clean architecture**: domain at the centre, application services around it, and infrastructure (HTTP, DB) on the outside. Dependencies point inward; the domain does not depend on frameworks or databases.

```
                    ┌─────────────────────────────────────────┐
                    │              main.go                     │
                    │  (wiring: config, DB, repos, services,   │
                    │   HTTP router, server)                   │
                    └────────────────────┬────────────────────┘
                                         │
         ┌──────────────────────────────┼──────────────────────────────┐
         │                              │                              │
         ▼                              ▼                              ▼
┌─────────────────┐          ┌─────────────────┐          ┌─────────────────┐
│  Infrastructure │          │    Service      │          │     Domain      │
│  (adapters)     │─────────▶│  (use cases)    │─────────▶│  (entities,     │
│                 │          │                 │          │   ports)        │
└─────────────────┘          └─────────────────┘          └─────────────────┘
  • HTTP (Gin, handlers)       • ingestion.Service           • event.Event, Record
  • Postgres (repos)            • insight.Service            • tenant.Tenant
  • Logger                     • IngestBatch, List           • user.User, insight.Insight
                                • Validate & enrich           • event/tenant/insight.Repository
```

- **Domain** (`internal/domain/`): Entities and repository interfaces (ports). No external dependencies.
- **Service** (`internal/service/`): Application use cases. Depends only on domain interfaces.
- **Infrastructure** (`internal/infrastructure/`): HTTP server, middleware, Postgres repositories. Implements domain ports and calls into services.
- **Config / Logger**: Loaded at startup; used by main and infrastructure.

---

## Project structure

```
behavix-ai/
├── main.go                    # Single entrypoint: server or migrate-up subcommand
├── go.mod
├── Makefile                   # make migrate-up
├── .env.example               # Copy to .env and set DB_* / PORT
│
├── internal/
│   ├── app/
│   │   └── migrate.go         # Migration runner (applies pending SQL files)
│   ├── config/
│   │   └── config.go          # Load from env; build PostgresDSN from DB_*
│   ├── logger/
│   │   └── logger.go          # Zap production logger
│   │
│   ├── domain/                # Core business entities and ports
│   │   ├── event/
│   │   │   ├── model.go       # Event, BatchRequest, Record
│   │   │   └── repository.go # InsertBatch, List
│   │   ├── tenant/
│   │   │   ├── model.go       # Tenant
│   │   │   └── repository.go  # GetTenantIDByAPIKey, GetByAPIKey
│   │   ├── user/
│   │   │   └── model.go       # User
│   │   └── insight/
│   │       ├── model.go       # Insight
│   │       └── repository.go  # List
│   │
│   ├── service/               # Application use cases
│   │   ├── ingestion/
│   │   │   └── service.go     # IngestBatch (validate, enrich, persist)
│   │   └── insight/
│   │       └── service.go     # List (insights feed)
│   │
│   └── infrastructure/        # Adapters (HTTP, DB)
│       ├── http/
│       │   ├── server/        # Router, middleware (auth, logging)
│       │   │   ├── router.go
│       │   │   └── middleware.go
│       │   ├── ingestion/
│       │   │   └── handler.go # POST /api/v1/events
│       │   ├── insight/
│       │   │   └── handler.go # GET /api/v1/insights
│       │   └── eventsdebug/
│       │       └── handler.go # GET /api/v1/events (debug list)
│       └── postgres/
│           ├── db.go
│           ├── tenant_repo.go
│           ├── event_repo.go
│           └── insight_repo.go
│
├── pkg/
│   └── response/              # HTTP helpers (Accepted, Error)
│       └── response.go
│
├── migrations/                # SQL migrations (run in order by name)
│   ├── 001_create_tenants.sql
│   ├── 002_create_events.sql
│   ├── 003_create_users.sql
│   ├── 004_create_insights.sql
│   └── 005_alter_events_schema.sql
│
└── scripts/
    └── dev.sh                 # Run server with .env loaded
```

---

## Setup

### Prerequisites

- **Go 1.21+**
- **PostgreSQL** (local or remote)

### 1. Clone and install

```bash
cd behavix-ai
go mod download
```

### 2. Environment

```bash
cp .env.example .env
```

Edit `.env`:

| Variable      | Description        | Example    |
|---------------|--------------------|------------|
| `PORT`        | HTTP server port   | `8080`     |
| `DB_HOST`     | Postgres host      | `localhost`|
| `DB_PORT`     | Postgres port      | `5432`     |
| `DB_USER`     | Postgres user      | `postgres` (or your OS user on macOS/Homebrew) |
| `DB_PASSWORD` | Postgres password  | (empty for local) |
| `DB_NAME`     | Database name      | `behavix`  |
| `DB_SSLMODE`  | SSL mode           | `disable`  |

You can instead set `POSTGRES_DSN` to a full connection string; it overrides the `DB_*` values.

### 3. Create the database

Using the same user as `DB_USER`:

```bash
createdb -h localhost behavix
# or: createdb -U postgres -h localhost behavix
```

Or inside `psql`:

```bash
psql -h localhost -c "CREATE DATABASE behavix;"
```

### 4. Run migrations

```bash
make migrate-up
```

This applies only migrations that have not been applied yet (tracked in `schema_migrations`).

### 5. Run the server

```bash
go run .
# or, with .env loaded from repo root:
./scripts/dev.sh
```

The server listens on `http://localhost:8080` (or your `PORT`).

---

## API

All `/api/v1` routes require **Authorization: Bearer &lt;API_KEY&gt;** (except where noted). The API key must exist in the `tenants` table; the middleware resolves it to `tenant_id` and attaches it to the request context.

### Health

- **GET** `/health`  
  No auth.  
  Response: `{"status":"ok"}`.

### Event ingestion

- **POST** `/api/v1/events`  
  Accept product usage events. Required fields per event: `user_id`, `event_name`. Server sets `id`, `tenant_id`, `received_at`.

**Request body:**

```json
{
  "events": [
    {
      "user_id": "user_123",
      "session_id": "session_456",
      "event_name": "project_created",
      "event_timestamp": "2026-03-06T10:00:00Z",
      "properties": { "project_type": "demo" },
      "context": { "platform": "web" }
    }
  ]
}
```

**Response:** `202 Accepted` with `{"status":"accepted"}`.  
Validation errors (e.g. missing `user_id` or `event_name`) return `400 Bad Request`.

### Insights feed

- **GET** `/api/v1/insights`  
  Returns recent behavioral insights for the tenant.

**Query params:**

| Param   | Description        | Default |
|---------|--------------------|---------|
| `limit` | Max number of rows | 20      |

**Response:** `200 OK` with a JSON array of insights:

```json
[
  {
    "id": "uuid",
    "type": "activation_drop",
    "title": "Activation dropped 31%",
    "description": "Activation dropped significantly compared to the 7-day average.",
    "severity": "high",
    "metric_value": 0.29,
    "baseline_value": 0.42,
    "created_at": "2026-03-06T12:00:00Z"
  }
]
```

### Events debug (list)

- **GET** `/api/v1/events`  
  List ingested events for the tenant (for debugging). Same auth as above.

**Query params:**

| Param        | Description              |
|-------------|--------------------------|
| `limit`     | Max events (default 20)  |
| `event_name`| Filter by event name     |
| `user_id`   | Filter by user id        |

**Response:** `200 OK` with a JSON array of event objects (id, tenant_id, user_id, session_id, event_name, event_timestamp, received_at, properties, context).

---

## Commands

| Command              | Description |
|----------------------|-------------|
| `go run .`           | Start the HTTP server (default). |
| `go run . migrate-up`| Apply all pending migrations. |
| `make migrate-up`    | Same as `go run . migrate-up`. |
| `./scripts/dev.sh`   | Start server with `.env` loaded. |

---

## Database

| Table        | Purpose |
|-------------|---------|
| **tenants** | SaaS companies; `id`, `name`, `api_key` (unique), `created_at`. Used for API-key auth. |
| **users**   | End users of the product (lifecycle); `id`, `tenant_id`, `external_user_id`, `email`, `created_at`, `metadata` (JSONB). Unique on `(tenant_id, external_user_id)`. |
| **events**  | Product usage events; `id`, `tenant_id`, `user_id`, `session_id` (nullable), `event_name`, `event_timestamp`, `received_at`, `properties` (JSONB), `context` (JSONB). Indexes on `(tenant_id, event_timestamp)`, `(tenant_id, user_id)`, `(tenant_id, event_name)`. |
| **insights**| Generated behavioral insights; `id`, `tenant_id`, `type`, `title`, `description`, `severity`, `metric_value`, `baseline_value`, `metadata` (JSONB), `created_at`. Index on `(tenant_id, created_at DESC)`. |

Migrations are in `migrations/`; order is by filename. Applied versions are stored in `schema_migrations`.

---

## Troubleshooting

| Issue | What to do |
|-------|------------|
| **Connection refused** | Start PostgreSQL (e.g. `brew services start postgresql@16`). |
| **Database "behavix" does not exist** | Create it: `createdb -h localhost behavix`. |
| **role "postgres" does not exist** | On macOS/Homebrew the default role is often your OS username. Set `DB_USER` in `.env` to that (e.g. `DB_USER=sumitnegi`). |
| **Password authentication failed** | Set `DB_PASSWORD` in `.env` to match the Postgres user. |
| **Invalid API key** (401 on `/api/v1/*`) | Insert a row in `tenants` with a valid `api_key` and use it in `Authorization: Bearer <key>`. |
