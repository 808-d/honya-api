-- name: FindUserByIdentifier :one
SELECT id, public_id::text AS public_id, username, display_name, hash_password,
       email, phone, is_verified
FROM users
WHERE is_delete = false
  AND (email = $1 OR username = $1)
LIMIT 1;

-- name: FindUserByEmail :one
SELECT id, public_id::text AS public_id, username, display_name, hash_password,
       email, phone, is_verified
FROM users
WHERE is_delete = false
  AND email = $1
LIMIT 1;

-- name: EmailExists :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE email = $1 AND is_delete = false
) AS email_exists;

-- name: UsernameExists :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE username = $1 AND is_delete = false
) AS username_exists;

-- name: CreateUser :one
INSERT INTO users (username, display_name, hash_password, email)
VALUES ($1, $2, $3, $4)
RETURNING id, public_id::text AS public_id, username, display_name, email;

-- name: VerifyUser :one
UPDATE users
SET is_verified = true
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false
RETURNING id, public_id::text AS public_id, username, display_name, email, phone;

-- name: AssignRoleByName :exec
INSERT INTO user_roles (user_id, role_id)
SELECT $1, id FROM roles WHERE name = $2 AND is_delete = false
ON CONFLICT DO NOTHING;

-- name: UpdateUserPassword :exec
UPDATE users
SET hash_password = sqlc.arg(hash_password)
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;

-- name: GetUserIDByPublicID :one
SELECT id
FROM users
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;

-- name: GetRoleNamesByUserID :many
SELECT r.name
FROM roles r
JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = $1
  AND r.is_delete = false
ORDER BY r.name;
