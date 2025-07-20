-- name: DeleteJobByID :one
DELETE FROM jobs
WHERE 
    id = $1
    AND employer_id = $2
RETURNING id;