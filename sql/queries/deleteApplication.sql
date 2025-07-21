-- name: DeleteAppByID :one
DELETE FROM applications
WHERE 
    id = $1
    AND applicant_id = $2
RETURNING id;