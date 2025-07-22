-- name: SearchJobs :many
SELECT *
FROM jobs
WHERE (
    title ILIKE '%' || $1 || '%'
    OR description ILIKE '%' || $1 || '%'
)
ORDER BY created_at DESC;
