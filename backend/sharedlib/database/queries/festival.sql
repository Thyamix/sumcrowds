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
