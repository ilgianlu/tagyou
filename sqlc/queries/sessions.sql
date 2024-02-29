-- name: CreateSession :exec
INSERT INTO sessions (
  last_seen,
  last_connect,
  expiry_interval,
  client_id,
  connected,
  protocol_version
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetSessionById :one
SELECT * FROM sessions
WHERE id = ?
LIMIT 1;

-- name: GetAllSessions :many
SELECT * FROM sessions;

-- name: SessionExists :one
SELECT * FROM sessions
WHERE client_id = ?
LIMIT 1;

-- name: DeleteSessionByClientId :exec
DELETE FROM sessions
WHERE client_id = ?
RETURNING *;

-- name: DisconnectSessionByClientId :exec
UPDATE sessions
SET
  connected = false,
  last_seen = ?
WHERE id = ?
RETURNING *;