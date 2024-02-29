-- name: CreateSubscription :exec
INSERT INTO subscriptions (
  client_id,
  topic,
  retain_handling,
  retain_as_published,
  no_local,
  qos,
  protocol_version,
  enabled,
  created_at,
  session_id,
  shared,
  share_name
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetSubscriptionToUnsubscribe :one
SELECT * FROM subscriptions
WHERE share_name = ? AND topic = ? AND client_id = ?
LIMIT 1;

-- name: GetSubscriptions :many
SELECT * FROM subscriptions
WHERE topic IN (?) AND shared = ?;

-- name: DeleteByClientIdTopicShareName :exec
DELETE FROM subscriptions
WHERE share_name = ? AND topic = ? AND client_id = ?
RETURNING *;

-- name: GetAllSubscriptions :many
SELECT * FROM subscriptions;