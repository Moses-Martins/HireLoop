-- name: GetApplyByJobID :many
SELECT * FROM applications
WHERE job_id = $1;