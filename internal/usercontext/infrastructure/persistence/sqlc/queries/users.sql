-- name: CreateUser :exec
INSERT INTO users (id, email, created_at)
VALUES ($1, $2, $3);

-- name: GetUserByEmail :one
SELECT id, email, created_at
FROM users
WHERE email = $1
LIMIT 1;

