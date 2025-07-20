-- name: GetUserByEmail :one
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
    email = $1;