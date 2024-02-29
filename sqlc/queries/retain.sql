-- name: GetAllRetains :many
SELECT * FROM retains;

-- name: CreateRetain :exec
INSERT INTO retains (
  client_id,
  topic,
  application_message,
  created_at
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteRetainByClientIdTopic :exec
DELETE FROM retains
WHERE client_id = ? AND topic = ?
RETURNING *;