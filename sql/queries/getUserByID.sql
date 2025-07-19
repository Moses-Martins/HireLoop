-- name: GetUserByID :one
SELECT 
    id,
    created_at,
    updated_at,
    name,
    email,
    hashed_password,
    role

FROM 
    users
WHERE 
    id = $1;