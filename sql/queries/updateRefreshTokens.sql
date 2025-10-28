-- name: RevokeToken :exec
UPDATE refresh_tokens
SET
    revoked_at = now(),
    updated_at = now()
WHERE
    token = $1;