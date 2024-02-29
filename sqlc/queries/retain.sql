-- name: FindRetains :many
SELECT * FROM retains;

-- name: CreateRetain :exec
INSERT INTO retains (
  client_id,
  topic,
  application_message
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: DeleteRetainById :exec
DELETE FROM retains
WHERE topic = ?
RETURNING *;