-- name: CreateClient :exec
INSERT INTO clients (
  client_id,
  username,
  password,
  subscribe_acl,
  publish_acl
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetAllClients :many
SELECT * FROM clients;

-- name: GetClientByClientIdUsername :one
SELECT * FROM clients
WHERE client_id = ? AND username = ?;

-- name: GetClientById :one
SELECT * FROM clients
WHERE id = ? LIMIT 1;

-- name: DeleteClientById :exec
DELETE FROM clients
WHERE id = ?
RETURNING *;