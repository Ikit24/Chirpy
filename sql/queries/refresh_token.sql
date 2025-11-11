-- name: InsertRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id, expires_at, revoked_at)
VALUES (
    $1,
    $2,
    NOW() + INTERVAL '60 days',
    NULL
    );
