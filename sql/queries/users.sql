-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, hashed_password, email)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)

RETURNING *;

-- name: ResetUser :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, hashed_password, email FROM users
WHERE email = $1;

-- name: UpdateUserByID :one
UPDATE users
SET updated_at = now(),
    hashed_password = COALESCE($2, hashed_password),
    email = COALESCE($3, email)
WHERE id = $1
RETURNING *;
