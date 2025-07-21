-- name: UpdateApplyJob :exec
UPDATE applications
SET 
    resume_url = $1
WHERE
    id = $2;


