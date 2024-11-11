-- name: CreateRetry :one
INSERT INTO retries (
  client_id,
  application_message,
  packet_identifier,
  qos,
  dup,
  retries,
  ack_status,
  created_at,
  session_id,
  reason_code
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetAllRetries :many
SELECT * FROM retries;

-- name: UpdateRetryAckStatus :exec
UPDATE retries
SET ack_status = ?
WHERE id = ?
RETURNING *;

-- name: GetRetryByClientIdPacketIdentifier :one
SELECT * FROM retries
WHERE client_id = ? AND packet_identifier = ?;

-- name: GetRetryByClientIdPacketIdentifierReasonCode :one
SELECT * FROM retries
WHERE client_id = ? AND packet_identifier = ? and reason_code = ?;

-- name: DeleteRetryById :exec
DELETE FROM retries
WHERE id = ?
RETURNING *;

-- name: DeleteRetryBySessionId :exec
DELETE FROM retries
WHERE session_id = ?
RETURNING *;

-- name: DeleteRetriesOlder :exec
DELETE FROM retries
WHERE created_at < ?
RETURNING *;