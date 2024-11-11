-- name: CreateSession :one
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

-- name: UpdateSession :one
UPDATE sessions
SET last_seen = ?,
    last_connect = ?,
    expiry_interval = ?,
    client_id = ?,
    connected = ?,
    protocol_version = ?
WHERE id = ?
RETURNING *;

-- name: GetSessionById :one
SELECT * FROM sessions
WHERE id = ?
LIMIT 1;

-- name: GetDisconnectedSessions :many
SELECT * FROM sessions
WHERE connected = 0;

-- name: GetAllSessions :many
SELECT * FROM sessions;

-- name: GetSessionByClientId :one
SELECT * FROM sessions
WHERE client_id = ?
LIMIT 1;

-- name: DeleteSessionById :exec
DELETE FROM sessions
WHERE id = ?
RETURNING *;

-- name: DeleteSessionByClientId :exec
DELETE FROM sessions
WHERE client_id = ?
RETURNING *;

-- name: DisconnectSessionByClientId :exec
UPDATE sessions
SET
  connected = false,
  last_seen = ?
WHERE client_id = ?
RETURNING *;