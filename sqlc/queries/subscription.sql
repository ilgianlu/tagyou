-- name: CreateSubscription :one
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

-- name: GetAllSubscriptions :many
SELECT * FROM subscriptions;

-- name: GetSubscriptionsBySessionId :many
SELECT * FROM subscriptions
WHERE session_id = ?;

-- name: GetSharedSubscriptionByNameTopicClientId :one
SELECT * FROM subscriptions
WHERE shared = 1 AND share_name = ? AND topic = ? AND client_id = ?
LIMIT 1;

-- name: GetSubscriptionByTopicClientId :one
SELECT * FROM subscriptions
WHERE shared = 0 AND topic = ? AND client_id = ?
LIMIT 1;

-- name: GetSubscriptions :many
SELECT * FROM subscriptions
WHERE topic IN (sqlc.slice(topics)) AND shared = ?;

-- name: GetSubscriptionsOrdered :many
SELECT * FROM subscriptions
WHERE topic IN (sqlc.slice(topics)) AND shared = ?
ORDER BY share_name;

-- name: DeleteSubscriptionByClientIdTopicShareName :exec
DELETE FROM subscriptions
WHERE share_name = ? AND topic = ? AND client_id = ?
RETURNING *;

