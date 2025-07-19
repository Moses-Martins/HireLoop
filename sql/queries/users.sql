-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, hashed_password, role)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3, 
    $4
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;
