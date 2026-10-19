-- name: GetOrCreateCart :one
WITH inserted AS (
    INSERT INTO carts (user_id)
    SELECT u.id
    FROM users u
    WHERE u.public_id::text = sqlc.arg(user_public_id)::text
      AND u.is_delete = false
    ON CONFLICT (user_id) DO NOTHING
    RETURNING id
)
SELECT id FROM inserted
UNION ALL
SELECT c.id
FROM carts c
JOIN users u ON u.id = c.user_id
WHERE u.public_id::text = sqlc.arg(user_public_id)::text
LIMIT 1;

-- name: GetCartItems :many
SELECT ci.id, p.cover, p.title, ci.quantity, ci.create_time
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.cart_id = sqlc.arg(cart_id)
  AND p.is_delete = false
ORDER BY ci.create_time;

-- name: CountCartItems :one
SELECT COALESCE(SUM(ci.quantity), 0)::bigint AS cart_count
FROM cart_items ci
JOIN products p ON p.id = ci.product_id AND p.is_delete = false
WHERE ci.cart_id = sqlc.arg(cart_id);

-- name: AddToCart :one
WITH upserted AS (
    INSERT INTO cart_items (cart_id, product_id, quantity)
    VALUES (sqlc.arg(cart_id), sqlc.arg(product_id), sqlc.arg(quantity))
    ON CONFLICT (cart_id, product_id)
    DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
    RETURNING id, product_id, quantity, create_time
)
SELECT u.id, p.cover, p.title, u.quantity, u.create_time
FROM upserted u
JOIN products p ON p.id = u.product_id;

-- name: RemoveFromCart :one
WITH deleted AS (
    DELETE FROM cart_items ci
    WHERE ci.id = sqlc.arg(id)
      AND ci.cart_id = sqlc.arg(cart_id)
    RETURNING ci.id, ci.product_id, ci.quantity, ci.create_time
)
SELECT d.id, p.cover, p.title, d.quantity, d.create_time
FROM deleted d
JOIN products p ON p.id = d.product_id;
