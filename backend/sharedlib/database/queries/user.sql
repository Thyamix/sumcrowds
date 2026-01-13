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
