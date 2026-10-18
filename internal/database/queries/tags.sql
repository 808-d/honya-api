-- name: ListTags :many
SELECT id, public_id::text AS public_id, name, create_time, update_time
FROM tags
WHERE is_delete = false
ORDER BY name
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountTags :one
SELECT count(*) FROM tags WHERE is_delete = false;

-- name: GetTagByPublicID :one
SELECT id, public_id::text AS public_id, name, create_time, update_time
FROM tags
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;

-- name: CreateTag :one
INSERT INTO tags (name)
VALUES (sqlc.arg(name))
RETURNING id, public_id::text AS public_id, name, create_time, update_time;

-- name: UpdateTag :one
UPDATE tags
SET name = sqlc.arg(name)
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false
RETURNING id, public_id::text AS public_id, name, create_time, update_time;

-- name: SoftDeleteTag :execrows
UPDATE tags
SET is_delete = true, delete_at = now()
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;
