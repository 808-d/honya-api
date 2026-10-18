-- name: ListPublishers :many
SELECT id, public_id::text AS public_id, name, create_time, update_time
FROM publishers
WHERE is_delete = false
ORDER BY name
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountPublishers :one
SELECT count(*) FROM publishers WHERE is_delete = false;

-- name: GetPublisherByPublicID :one
SELECT id, public_id::text AS public_id, name, create_time, update_time
FROM publishers
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;

-- name: CreatePublisher :one
INSERT INTO publishers (name)
VALUES (sqlc.arg(name))
RETURNING id, public_id::text AS public_id, name, create_time, update_time;

-- name: UpdatePublisher :one
UPDATE publishers
SET name = sqlc.arg(name)
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false
RETURNING id, public_id::text AS public_id, name, create_time, update_time;

-- name: SoftDeletePublisher :execrows
UPDATE publishers
SET is_delete = true, delete_at = now()
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;
