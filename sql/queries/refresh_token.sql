-- name: InsertRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id, created_at, updated_at, expires_at, revoked_at)
VALUES (
    $1,
    $2,
    NOW(),
    NOW(),
    NOW() + INTERVAL '60 days',
    NULL
);

-- name: GetUserFromRefreshToken :one
SELECT user_id
FROM refresh_tokens
WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL;

-- name: UpdateRevokedRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1 AND revoked_at IS NULL;
