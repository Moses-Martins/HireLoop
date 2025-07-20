-- name: CreateJobs :one
INSERT INTO jobs (id, created_at, updated_at, title, description, location, type, salary, employer_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5, 
    $6
)
RETURNING *;



