-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    token,
    created_at,
    updated_at,
    user_id,
    expires_at,
    revoked_at
)
VALUES (
    $1,
    now(),
    now(),
    $2,
    now() + INTERVAL '60 days',
    NULL
)
RETURNING *;

-- name: RetrieveRefreshToken :one
SELECT * from refresh_tokens
WHERE token = $1
AND revoked_at IS NULL
AND expires_at > now();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = now(),
    updated_at = now()
WHERE token = $1;
