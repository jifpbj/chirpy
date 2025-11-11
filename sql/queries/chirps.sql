-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)

RETURNING *;

-- name: RetrieveChirps :many
SELECT * from chirps
ORDER BY created_at ASC;

-- name: RetrieveChirpByID :one
SELECT * from chirps
WHERE id = $1;
