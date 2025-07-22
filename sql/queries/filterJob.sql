-- name: FiltersJobs :many
SELECT *
FROM jobs
WHERE
    location ILIKE '%' || $1 || '%'
    AND ($2 = '' OR type = $2)
    AND ($3 = 1 OR salary >= $3)
    AND ($4 = 1 OR salary <= $4)
ORDER BY created_at DESC;