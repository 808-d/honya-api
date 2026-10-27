-- name: UserHasFeature :one
SELECT EXISTS (
    SELECT 1
    FROM users u
    JOIN user_roles ur    ON ur.user_id = u.id
    JOIN roles r          ON r.id = ur.role_id AND r.is_delete = false
    JOIN role_features rf ON rf.role_id = r.id
    JOIN features f       ON f.id = rf.feature_id AND f.is_delete = false
    WHERE u.public_id::text = sqlc.arg(public_id)::text
      AND u.is_delete = false
      AND f.code = sqlc.arg(code)::text
) AS has_feature;

-- ============================================================
-- Roles
-- ============================================================

-- name: ListRoles :many
SELECT r.id,
       r.public_id::text AS public_id,
       r.name,
       r.description,
       r.create_time,
       r.update_time,
       COALESCE(
           (SELECT array_agg(f.code ORDER BY f.code)
            FROM role_features rf
            JOIN features f ON f.id = rf.feature_id AND f.is_delete = false
            WHERE rf.role_id = r.id),
           ARRAY[]::text[]
       )::text[] AS features,
       (SELECT count(*) FROM user_roles ur WHERE ur.role_id = r.id)::bigint AS member_count
FROM roles r
WHERE r.is_delete = false
ORDER BY r.name;

-- name: GetRoleByPublicID :one
SELECT r.id,
       r.public_id::text AS public_id,
       r.name,
       r.description,
       r.create_time,
       r.update_time,
       COALESCE(
           (SELECT array_agg(f.code ORDER BY f.code)
            FROM role_features rf
            JOIN features f ON f.id = rf.feature_id AND f.is_delete = false
            WHERE rf.role_id = r.id),
           ARRAY[]::text[]
       )::text[] AS features,
       (SELECT count(*) FROM user_roles ur WHERE ur.role_id = r.id)::bigint AS member_count
FROM roles r
WHERE r.public_id::text = sqlc.arg(public_id)::text
  AND r.is_delete = false;

-- name: CreateRole :one
INSERT INTO roles (name, description)
VALUES (sqlc.arg(name), sqlc.arg(description))
RETURNING id, public_id::text AS public_id, name, description, create_time, update_time;

-- name: UpdateRole :one
UPDATE roles
SET name        = sqlc.arg(name),
    description = sqlc.arg(description)
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false
RETURNING id, public_id::text AS public_id, name, description, create_time, update_time;

-- name: SoftDeleteRole :execrows
UPDATE roles
SET is_delete = true, delete_at = now()
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;

-- ============================================================
-- Features
-- ============================================================

-- name: ListFeatures :many
SELECT id, public_id::text AS public_id, code, name, description
FROM features
WHERE is_delete = false
ORDER BY code;

-- name: DeleteRoleFeatures :exec
DELETE FROM role_features WHERE role_id = sqlc.arg(role_id);

-- name: InsertRoleFeature :execrows
INSERT INTO role_features (role_id, feature_id)
SELECT sqlc.arg(role_id), id
FROM features
WHERE public_id::text = sqlc.arg(feature_public_id)::text
  AND is_delete = false
ON CONFLICT DO NOTHING;

-- ============================================================
-- Users and their roles
-- ============================================================

-- name: ListUsers :many
SELECT u.id,
       u.public_id::text AS public_id,
       u.username,
       u.display_name,
       u.email,
       u.phone,
       u.is_verified,
       u.create_time,
       COALESCE(
           (SELECT array_agg(r.name ORDER BY r.name)
            FROM user_roles ur
            JOIN roles r ON r.id = ur.role_id AND r.is_delete = false
            WHERE ur.user_id = u.id),
           ARRAY[]::text[]
       )::text[] AS roles
FROM users u
WHERE u.is_delete = false
  AND (sqlc.narg(search)::text IS NULL
       OR u.username ILIKE '%' || sqlc.narg(search)::text || '%'
       OR u.email ILIKE '%' || sqlc.narg(search)::text || '%'
       OR u.display_name ILIKE '%' || sqlc.narg(search)::text || '%')
ORDER BY u.create_time DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountUsers :one
SELECT count(*)
FROM users u
WHERE u.is_delete = false
  AND (sqlc.narg(search)::text IS NULL
       OR u.username ILIKE '%' || sqlc.narg(search)::text || '%'
       OR u.email ILIKE '%' || sqlc.narg(search)::text || '%'
       OR u.display_name ILIKE '%' || sqlc.narg(search)::text || '%');

-- name: GetUserByPublicID :one
SELECT u.id,
       u.public_id::text AS public_id,
       u.username,
       u.display_name,
       u.email,
       u.phone,
       u.is_verified,
       u.create_time,
       COALESCE(
           (SELECT array_agg(r.name ORDER BY r.name)
            FROM user_roles ur
            JOIN roles r ON r.id = ur.role_id AND r.is_delete = false
            WHERE ur.user_id = u.id),
           ARRAY[]::text[]
       )::text[] AS roles
FROM users u
WHERE u.public_id::text = sqlc.arg(public_id)::text
  AND u.is_delete = false;

-- name: DeleteUserRoles :exec
DELETE FROM user_roles WHERE user_id = sqlc.arg(user_id);

-- name: InsertUserRole :execrows
INSERT INTO user_roles (user_id, role_id)
SELECT sqlc.arg(user_id), id
FROM roles
WHERE public_id::text = sqlc.arg(role_public_id)::text
  AND is_delete = false
ON CONFLICT DO NOTHING;

-- name: GetRoleNamesByPublicID :many
SELECT r.name
FROM roles r
JOIN user_roles ur ON ur.role_id = r.id
JOIN users u ON u.id = ur.user_id
WHERE u.public_id::text = sqlc.arg(public_id)::text
  AND u.is_delete = false
  AND r.is_delete = false
ORDER BY r.name;

-- name: GetFeatureCodesByPublicID :many
SELECT DISTINCT f.code
FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id AND r.is_delete = false
JOIN role_features rf ON rf.role_id = r.id
JOIN features f ON f.id = rf.feature_id AND f.is_delete = false
WHERE u.public_id::text = sqlc.arg(public_id)::text
  AND u.is_delete = false
ORDER BY f.code;

-- name: CountRolesByPublicIDs :one
SELECT count(*)
FROM roles
WHERE is_delete = false
  AND public_id::text = ANY(sqlc.arg(public_ids)::text[]);

-- name: CountFeaturesByPublicIDs :one
SELECT count(*)
FROM features
WHERE is_delete = false
  AND public_id::text = ANY(sqlc.arg(public_ids)::text[]);

-- name: UserHasRole :one
SELECT EXISTS (
    SELECT 1
    FROM users u
    JOIN user_roles ur ON ur.user_id = u.id
    JOIN roles r ON r.id = ur.role_id AND r.is_delete = false
    WHERE u.public_id::text = sqlc.arg(user_public_id)::text
      AND u.is_delete = false
      AND r.public_id::text = sqlc.arg(role_public_id)::text
) AS has_role;

-- name: GetFeatureCodesByPublicIDExcludingRole :many
SELECT DISTINCT f.code
FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id AND r.is_delete = false
JOIN role_features rf ON rf.role_id = r.id
JOIN features f ON f.id = rf.feature_id AND f.is_delete = false
WHERE u.public_id::text = sqlc.arg(public_id)::text
  AND u.is_delete = false
  AND r.id <> sqlc.arg(role_id)
ORDER BY f.code;
