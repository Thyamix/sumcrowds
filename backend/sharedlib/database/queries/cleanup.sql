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
