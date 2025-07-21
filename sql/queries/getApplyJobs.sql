-- name: GetApplyJobs :one
SELECT * FROM applications
WHERE id = $1;