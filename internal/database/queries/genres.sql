-- name: ListGenres :many
SELECT id, public_id::text AS public_id, name, create_time, update_time
FROM genres
WHERE is_delete = false
ORDER BY name
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountGenres :one
SELECT count(*) FROM genres WHERE is_delete = false;

-- name: GetGenreByPublicID :one
SELECT id, public_id::text AS public_id, name, create_time, update_time
FROM genres
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;

-- name: CreateGenre :one
INSERT INTO genres (name)
VALUES (sqlc.arg(name))
RETURNING id, public_id::text AS public_id, name, create_time, update_time;

-- name: UpdateGenre :one
UPDATE genres
SET name = sqlc.arg(name)
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false
RETURNING id, public_id::text AS public_id, name, create_time, update_time;

-- name: SoftDeleteGenre :execrows
UPDATE genres
SET is_delete = true, delete_at = now()
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;
