-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  token,
  created_at,
  updated_at,
  user_id,
  expires_at,
  revoked_at
) VALUES ( 
  $1,
  NOW(),
  NOW(),
  $2,
  NOW() + '60 days',
  NULL
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens 
SET revoked_at = NOW(),
  updated_at = NOW()
WHERE token = $1
RETURNING *;
