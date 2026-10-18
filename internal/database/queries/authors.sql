-- name: ListAuthors :many
SELECT id, public_id::text AS public_id, name, create_time, update_time
FROM authors
WHERE is_delete = false
ORDER BY name
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountAuthors :one
SELECT count(*) FROM authors WHERE is_delete = false;

-- name: GetAuthorByPublicID :one
SELECT id, public_id::text AS public_id, name, create_time, update_time
FROM authors
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;

-- name: CreateAuthor :one
INSERT INTO authors (name)
VALUES (sqlc.arg(name))
RETURNING id, public_id::text AS public_id, name, create_time, update_time;

-- name: UpdateAuthor :one
UPDATE authors
SET name = sqlc.arg(name)
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false
RETURNING id, public_id::text AS public_id, name, create_time, update_time;

-- name: SoftDeleteAuthor :execrows
UPDATE authors
SET is_delete = true, delete_at = now()
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;
