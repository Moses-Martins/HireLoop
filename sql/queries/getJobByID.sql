-- name: GetJobsByID :one
SELECT * FROM jobs
WHERE id = $1;