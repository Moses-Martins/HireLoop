-- name: UpdateJob :exec
UPDATE jobs
SET 
    title = $1,
    updated_at = NOW(),
    description = $2, 
    location = $3,
    type = $4, 
    salary = $5

WHERE
    id = $6;