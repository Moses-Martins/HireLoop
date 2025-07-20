-- name: CreateJobs :one
INSERT INTO jobs (id, title, description, location, type, salary, employer_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5, 
    $6
)
RETURNING *;



