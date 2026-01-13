# SUM-14: Switch to SQLC - Implementation Plan

## Overview

Migrate the backend database layer from raw SQL queries with `database/sql` to [SQLC](https://sqlc.dev/), a compile-time SQL-to-Go code generator that provides type-safe database access.

## Current State

### Database Files
- `backend/sharedlib/database/database.go` - DB connection and initialization
- `backend/sharedlib/database/init.sql` - Schema definitions (9 tables)
- `backend/counter/internal/database/festival.go` - Festival CRUD operations
- `backend/counter/internal/database/event.go` - Event operations and archiving
- `backend/counter/internal/database/user.go` - User, token, and session management
- `backend/cleanup/internal/database/counter.go` - Cleanup/expiration queries

### Tables (from init.sql)
1. `festival` - Core festival data
2. `event` - Events per festival
3. `active` - Active event values
4. `archive` - Archived event values
5. `gauge_max` - Maximum gauge history
6. `app_user` - Users
7. `refresh_token` - JWT refresh tokens
8. `access_token` - JWT access tokens
9. `festival_access` - User-festival access records

### Current Query Patterns
- Simple SELECT/INSERT/UPDATE/DELETE queries
- Joins between tables (festival-event, festival_access-festival)
- Transactions for multi-step operations
- Pagination with LIMIT/OFFSET

## Implementation Steps

### Phase 1: Setup and Configuration

1. **Install SQLC**
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

2. **Create SQLC configuration file**
   Create `backend/sqlc.yaml`:
   ```yaml
   version: "2"
   sql:
     - engine: "postgresql"
       queries: "counter/internal/database/queries/"
       schema: "sharedlib/database/"
       gen:
         go:
           package: "sqlcdb"
           out: "counter/internal/database/sqlcdb"
           sql_package: "pgx/v5"
           emit_json_tags: true
           emit_prepared_queries: false
           emit_interface: true
           emit_exact_table_names: false
   ```

3. **Create queries directory structure**
   ```
   backend/
   ├── counter/internal/database/
   │   ├── queries/
   │   │   ├── festival.sql
   │   │   ├── event.sql
   │   │   ├── user.sql
   │   │   └── cleanup.sql
   │   └── sqlcdb/  (generated)
   ```

### Phase 2: Convert Queries

#### festival.sql
```sql
-- name: GetFestival :one
SELECT id, last_used_at, created_at, pin, password, code
FROM festival WHERE code = $1;

-- name: CreateFestival :one
INSERT INTO festival (code, password, pin, created_at, last_used_at)
VALUES ($1, $2, $3, $4, $4) RETURNING id;

-- name: IsFestivalCodeExists :one
SELECT EXISTS(SELECT 1 FROM festival WHERE code = $1);

-- name: UpdateFestivalLastUsedAt :exec
UPDATE festival SET last_used_at = $1 WHERE id = $2;
```

#### event.sql
```sql
-- name: GetActiveEventIdByFestivalCode :one
SELECT e.id FROM event e
JOIN festival f ON f.id = e.festival_id
WHERE e.active = TRUE AND f.code = $1;

-- name: GetTotalFromEvent :one
SELECT total FROM event WHERE id = $1;

-- name: GetLatestGaugeMax :one
SELECT gauge_max FROM gauge_max
WHERE event_id = $1
ORDER BY time DESC LIMIT 1;

-- name: InsertActiveValue :exec
INSERT INTO active (value, time, event_id) VALUES ($1, $2, $3);

-- name: UpdateEventTotal :exec
UPDATE event SET total = total + $1 WHERE id = $2;

-- name: CreateEvent :one
INSERT INTO event (created_at, last_used_at, festival_id, active, total)
VALUES ($1, $1, $2, TRUE, 0) RETURNING id;

-- name: DeactivateEvent :exec
UPDATE event SET active = FALSE, last_used_at = $1 WHERE id = $2;

-- name: InsertGaugeMax :exec
INSERT INTO gauge_max (gauge_max, time, event_id) VALUES ($1, $2, $3);

-- name: ArchiveActiveToArchive :exec
INSERT INTO archive (value, time, event_id)
SELECT value, time, event_id FROM active WHERE event_id = $1;

-- name: DeleteActiveByEventId :exec
DELETE FROM active WHERE event_id = $1;

-- name: GetArchivedEventIds :many
SELECT id, last_used_at FROM event
WHERE active = FALSE AND festival_id = $1
ORDER BY id DESC;

-- name: CountActiveEntries :one
SELECT COUNT(*) FROM active WHERE event_id = $1;

-- name: CountArchiveEntries :one
SELECT COUNT(*) FROM archive WHERE event_id = $1;

-- name: GetActiveEntriesChunk :many
SELECT t.value, t.time,
  COALESCE((SELECT gm.gauge_max FROM gauge_max gm
    WHERE gm.event_id = t.event_id AND gm.time <= t.time
    ORDER BY gm.time DESC LIMIT 1), 0) AS current_gauge_max
FROM active t WHERE t.event_id = $1
ORDER BY t.id LIMIT $2 OFFSET $3;

-- name: GetArchiveEntriesChunk :many
SELECT t.value, t.time,
  COALESCE((SELECT gm.gauge_max FROM gauge_max gm
    WHERE gm.event_id = t.event_id AND gm.time <= t.time
    ORDER BY gm.time DESC LIMIT 1), 0) AS current_gauge_max
FROM archive t WHERE t.event_id = $1
ORDER BY t.id LIMIT $2 OFFSET $3;

-- name: UpdateEventLastUsedAt :exec
UPDATE event SET last_used_at = $1 WHERE id = $2;
```

#### user.sql
```sql
-- name: CreateUser :one
INSERT INTO app_user DEFAULT VALUES RETURNING id;

-- name: CreateRefreshToken :exec
INSERT INTO refresh_token (token, expires_at, user_id, revoked)
VALUES ($1, $2, $3, $4);

-- name: CreateAccessToken :exec
INSERT INTO access_token (token, expires_at, user_id, revoked)
VALUES ($1, $2, $3, $4);

-- name: GetRefreshToken :one
SELECT id, user_id, token, expires_at, revoked
FROM refresh_token WHERE token = $1;

-- name: GetAccessToken :one
SELECT id, token, expires_at, user_id, revoked
FROM access_token WHERE token = $1;

-- name: UpdateRefreshToken :exec
UPDATE refresh_token SET expires_at = $1, revoked = $2 WHERE id = $3;

-- name: UpdateAccessToken :exec
UPDATE access_token SET expires_at = $1, revoked = $2 WHERE id = $3;

-- name: DeleteAccessToken :exec
DELETE FROM access_token WHERE token = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_token WHERE token = $1;

-- name: GetUserIdFromAccessToken :one
SELECT user_id FROM access_token WHERE token = $1;

-- name: CreateFestivalAccess :exec
INSERT INTO festival_access (festival_id, user_id, last_used_at, revoked)
VALUES ($1, $2, $3, $4);

-- name: GetFestivalAccess :one
SELECT id, festival_id, last_used_at, user_id
FROM festival_access
WHERE user_id = $1 AND festival_id = $2;

-- name: UpdateFestivalAccessLastUsedAt :exec
UPDATE festival_access SET last_used_at = $1 WHERE id = $2;

-- name: GetUserRecentSessions :many
SELECT f.code, fa.last_used_at
FROM festival_access fa
JOIN festival f ON f.id = fa.festival_id
WHERE fa.user_id = $1 AND fa.revoked = false
ORDER BY fa.last_used_at DESC
LIMIT $2 OFFSET $3;
```

#### cleanup.sql
```sql
-- name: DeleteExpiredFestivals :exec
DELETE FROM festival WHERE last_used_at < $1;

-- name: DeleteExpiredEvents :exec
DELETE FROM event WHERE last_used_at < $1;

-- name: DeleteExpiredFestivalAccess :exec
DELETE FROM festival_access WHERE last_used_at < $1;

-- name: DeleteExpiredAccessTokens :exec
DELETE FROM access_token WHERE expires_at < $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_token WHERE expires_at < $1;
```

### Phase 3: Database Connection Update

1. **Update database.go to use pgx**
   - Replace `database/sql` with `github.com/jackc/pgx/v5/pgxpool`
   - Update connection handling for connection pool
   - Keep backward compatibility during migration

2. **Create SQLC Queries wrapper**
   ```go
   package database

   import (
       "github.com/jackc/pgx/v5/pgxpool"
       "github.com/thyamix/sumcrowds/backend/counter/internal/database/sqlcdb"
   )

   var Pool *pgxpool.Pool
   var Queries *sqlcdb.Queries

   func InitDB() {
       // ... existing connection logic
       Pool, err = pgxpool.New(ctx, connStr)
       Queries = sqlcdb.New(Pool)
   }
   ```

### Phase 4: Refactor Existing Code

1. **Replace manual queries with SQLC calls**
   - Update each function to use generated queries
   - Maintain same function signatures for backward compatibility
   - Handle transaction wrapping where needed

2. **Transaction handling**
   - SQLC supports `WithTx()` for transactions
   - Example:
     ```go
     tx, _ := Pool.Begin(ctx)
     defer tx.Rollback(ctx)
     qtx := Queries.WithTx(tx)
     // use qtx for queries
     tx.Commit(ctx)
     ```

### Phase 5: Testing and Migration

1. **Add integration tests**
   - Test all SQLC-generated queries
   - Compare behavior with existing implementation

2. **Gradual rollout**
   - Keep old functions alongside new ones initially
   - Switch handlers one at a time
   - Remove old code once validated

## Benefits

1. **Type Safety** - Compile-time SQL validation
2. **No Runtime Errors** - SQL errors caught at build time
3. **Better Performance** - Optimized query execution
4. **Maintainability** - SQL files are easy to read and modify
5. **IDE Support** - Better autocomplete for generated Go code
6. **Schema Changes** - Automatic detection of breaking changes

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Breaking existing functionality | Gradual migration with parallel implementations |
| Learning curve | Follow SQLC documentation, simple query patterns |
| pgx vs database/sql differences | Test thoroughly before full migration |

## Estimated Effort

- Phase 1 (Setup): Small
- Phase 2 (Convert Queries): Medium
- Phase 3 (Connection Update): Small
- Phase 4 (Refactor): Medium-Large
- Phase 5 (Testing): Medium

## Dependencies

- `github.com/sqlc-dev/sqlc` - CLI tool
- `github.com/jackc/pgx/v5` - PostgreSQL driver (recommended by SQLC)

## References

- [SQLC Documentation](https://docs.sqlc.dev/)
- [SQLC PostgreSQL Guide](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html)
- [pgx Documentation](https://github.com/jackc/pgx)
