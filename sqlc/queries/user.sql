-- name: CreateUser :exec
INSERT INTO users (
  username,
  password,
  created_at
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ?;

-- name: DeleteUserById :exec
DELETE FROM users
WHERE id = ?
RETURNING *;