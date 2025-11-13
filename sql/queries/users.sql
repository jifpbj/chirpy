-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, hashed_password, email, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2,
    FALSE
)

RETURNING *;

-- name: ResetUser :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUserByID :one
UPDATE users
SET updated_at = now(),
    hashed_password = COALESCE($2, hashed_password),
    email = COALESCE($3, email)
WHERE id = $1
RETURNING *;

-- name: UpgradeUserByID :exec
UPDATE users
SET updated_at = now(),
    is_chirpy_red = true
WHERE id = $1;
