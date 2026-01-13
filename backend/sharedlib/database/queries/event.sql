-- name: GetActiveEventIdByFestivalCode :one
SELECT e.id FROM event e
JOIN festival f ON f.id = e.festival_id
WHERE e.active = TRUE AND f.code = $1;

-- name: GetActiveEventIdByFestivalId :one
SELECT id FROM event WHERE active = TRUE AND festival_id = $1;

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
SELECT active.value, active.time, active.event_id FROM active WHERE active.event_id = $1;

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
SELECT a.value, a.time,
  COALESCE((SELECT gm.gauge_max FROM gauge_max gm
    WHERE gm.event_id = a.event_id AND gm.time <= a.time
    ORDER BY gm.time DESC LIMIT 1), 0) AS current_gauge_max
FROM active a WHERE a.event_id = $1
ORDER BY a.id LIMIT $2 OFFSET $3;

-- name: GetArchiveEntriesChunk :many
SELECT ar.value, ar.time,
  COALESCE((SELECT gm.gauge_max FROM gauge_max gm
    WHERE gm.event_id = ar.event_id AND gm.time <= ar.time
    ORDER BY gm.time DESC LIMIT 1), 0) AS current_gauge_max
FROM archive ar WHERE ar.event_id = $1
ORDER BY ar.id LIMIT $2 OFFSET $3;

-- name: UpdateEventLastUsedAt :exec
UPDATE event SET last_used_at = $1 WHERE id = $2;

-- name: CheckEventExists :one
SELECT COUNT(1) FROM event WHERE id = $1;
