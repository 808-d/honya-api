-- name: ListProducts :many
SELECT
    p.id,
    p.public_id::text AS public_id,
    p.title,
    p.cover,
    p.description,
    (COALESCE(pub.public_id::text, ''))::text AS publisher_id,
    p.publish_time,
    p.stock,
    p.price::float8 AS price,
    p.status,
    p.type,
    p.create_time,
    p.update_time
FROM products p
LEFT JOIN publishers pub ON pub.id = p.publisher_id
WHERE p.is_delete = false
  AND (sqlc.narg(filter_status)::smallint IS NULL OR p.status = sqlc.narg(filter_status)::smallint)
  AND (sqlc.narg(filter_type)::smallint IS NULL OR p.type = sqlc.narg(filter_type)::smallint)
  AND (sqlc.narg(search)::text IS NULL OR p.title ILIKE '%' || sqlc.narg(search)::text || '%')
  AND (sqlc.narg(min_price)::float8 IS NULL OR p.price >= sqlc.narg(min_price)::float8)
  AND (sqlc.narg(max_price)::float8 IS NULL OR p.price <= sqlc.narg(max_price)::float8)
ORDER BY p.create_time DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountProducts :one
SELECT count(*)
FROM products p
WHERE p.is_delete = false
  AND (sqlc.narg(filter_status)::smallint IS NULL OR p.status = sqlc.narg(filter_status)::smallint)
  AND (sqlc.narg(filter_type)::smallint IS NULL OR p.type = sqlc.narg(filter_type)::smallint)
  AND (sqlc.narg(search)::text IS NULL OR p.title ILIKE '%' || sqlc.narg(search)::text || '%')
  AND (sqlc.narg(min_price)::float8 IS NULL OR p.price >= sqlc.narg(min_price)::float8)
  AND (sqlc.narg(max_price)::float8 IS NULL OR p.price <= sqlc.narg(max_price)::float8);

-- name: GetProductByPublicID :one
SELECT
    p.id,
    p.public_id::text AS public_id,
    p.title,
    p.cover,
    p.description,
    (COALESCE(pub.public_id::text, ''))::text AS publisher_id,
    p.publish_time,
    p.stock,
    p.price::float8 AS price,
    p.status,
    p.type,
    p.create_time,
    p.update_time
FROM products p
LEFT JOIN publishers pub ON pub.id = p.publisher_id
WHERE p.public_id::text = sqlc.arg(public_id)::text
  AND p.is_delete = false;

-- name: CreateProduct :one
INSERT INTO products (
    title, cover, description, publisher_id, publish_time, stock, price, status, type
)
VALUES (
    sqlc.arg(title),
    sqlc.arg(cover),
    sqlc.arg(description),
    (SELECT id FROM publishers WHERE public_id::text = sqlc.narg(publisher_id)::text),
    sqlc.narg(publish_time),
    sqlc.arg(stock),
    sqlc.arg(price)::float8,
    sqlc.arg(status),
    sqlc.arg(type)
)
RETURNING id, public_id::text AS public_id;

-- name: UpdateProduct :one
UPDATE products
SET
    title = sqlc.arg(title),
    cover = sqlc.arg(cover),
    description = sqlc.arg(description),
    publisher_id = (SELECT id FROM publishers WHERE public_id::text = sqlc.narg(publisher_id)::text),
    publish_time = sqlc.narg(publish_time),
    stock = sqlc.arg(stock),
    price = sqlc.arg(price)::float8,
    status = sqlc.arg(status),
    type = sqlc.arg(type)
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false
RETURNING id, public_id::text AS public_id;

-- name: SoftDeleteProduct :execrows
UPDATE products
SET is_delete = true, delete_at = now()
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;

-- name: GetProductAuthorIDs :many
SELECT a.public_id::text AS public_id
FROM product_authors pa
JOIN authors a ON a.id = pa.author_id
WHERE pa.product_id = $1
ORDER BY a.name;

-- name: GetProductGenreIDs :many
SELECT g.public_id::text AS public_id
FROM product_genres pg
JOIN genres g ON g.id = pg.genre_id
WHERE pg.product_id = $1
ORDER BY g.name;

-- name: GetProductTagIDs :many
SELECT t.public_id::text AS public_id
FROM product_tags pt
JOIN tags t ON t.id = pt.tag_id
WHERE pt.product_id = $1
ORDER BY t.name;

-- name: ClearProductAuthors :exec
DELETE FROM product_authors WHERE product_id = $1;

-- name: ClearProductGenres :exec
DELETE FROM product_genres WHERE product_id = $1;

-- name: ClearProductTags :exec
DELETE FROM product_tags WHERE product_id = $1;

-- name: AddProductAuthors :exec
INSERT INTO product_authors (product_id, author_id)
SELECT sqlc.arg(product_id), a.id
FROM authors a
WHERE a.public_id::text = ANY(sqlc.arg(public_ids)::text[])
  AND a.is_delete = false
ON CONFLICT DO NOTHING;

-- name: AddProductGenres :exec
INSERT INTO product_genres (product_id, genre_id)
SELECT sqlc.arg(product_id), g.id
FROM genres g
WHERE g.public_id::text = ANY(sqlc.arg(public_ids)::text[])
  AND g.is_delete = false
ON CONFLICT DO NOTHING;

-- name: AddProductTags :exec
INSERT INTO product_tags (product_id, tag_id)
SELECT sqlc.arg(product_id), t.id
FROM tags t
WHERE t.public_id::text = ANY(sqlc.arg(public_ids)::text[])
  AND t.is_delete = false
ON CONFLICT DO NOTHING;
