-- name: ApplyJobs :one
INSERT INTO applications (id, applicant_id, job_id, resume_url, status)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

