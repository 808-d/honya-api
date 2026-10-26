-- ============================================================
-- Orders
-- ============================================================

-- name: LockCartForCheckout :many
SELECT ci.id AS cart_item_id,
       ci.cart_id,
       p.id AS product_id,
       p.public_id::text AS product_public_id,
       p.title AS product_title,
       p.cover AS product_cover,
       p.price::float8 AS unit_price,
       ci.quantity,
       p.stock,
       p.status,
       p.is_delete
FROM cart_items ci
JOIN carts c ON c.id = ci.cart_id
JOIN users u ON u.id = c.user_id
JOIN products p ON p.id = ci.product_id
WHERE u.public_id::text = sqlc.arg(user_public_id)::text
  AND u.is_delete = false
ORDER BY ci.id
FOR UPDATE OF ci;

-- name: CreateOrder :one
INSERT INTO orders (
    user_id, coupon_id,
    subtotal, discount_amount, shipping_fee, total,
    status, payment_status, payment_method,
    shipping_name, shipping_phone, shipping_address, note
)
VALUES (
    sqlc.arg(user_id),
    sqlc.narg(coupon_id),
    sqlc.arg(subtotal)::float8,
    sqlc.arg(discount_amount)::float8,
    sqlc.arg(shipping_fee)::float8,
    sqlc.arg(total)::float8,
    sqlc.arg(status),
    sqlc.arg(payment_status),
    sqlc.arg(payment_method),
    sqlc.narg(shipping_name),
    sqlc.narg(shipping_phone),
    sqlc.narg(shipping_address),
    sqlc.narg(note)
)
RETURNING id, public_id::text AS public_id, create_time, update_time;

-- name: CreateOrderItem :exec
INSERT INTO order_items (
    order_id, product_id, product_title, product_cover,
    unit_price, quantity, subtotal
)
VALUES (
    sqlc.arg(order_id),
    sqlc.arg(product_id),
    sqlc.arg(product_title),
    sqlc.arg(product_cover),
    sqlc.arg(unit_price)::float8,
    sqlc.arg(quantity),
    sqlc.arg(subtotal)::float8
);

-- name: DecrementProductStock :execrows
UPDATE products
SET stock = stock - sqlc.arg(quantity),
    update_time = now()
WHERE id = sqlc.arg(product_id)
  AND is_delete = false
  AND stock >= sqlc.arg(quantity);

-- name: RestoreProductStock :exec
UPDATE products
SET stock = stock + sqlc.arg(quantity),
    update_time = now()
WHERE id = sqlc.arg(product_id);

-- name: DeleteCartItems :exec
DELETE FROM cart_items
WHERE cart_id = sqlc.arg(cart_id);

-- ============================================================
-- Reads
-- ============================================================

-- name: GetOrderByPublicID :one
SELECT o.id,
       o.public_id::text AS public_id,
       o.user_id,
       COALESCE(c.code, '') AS coupon_code,
       o.subtotal::float8 AS subtotal,
       o.discount_amount::float8 AS discount_amount,
       o.shipping_fee::float8 AS shipping_fee,
       o.total::float8 AS total,
       o.status,
       o.payment_status,
       o.payment_method,
       COALESCE(o.shipping_name, '') AS shipping_name,
       COALESCE(o.shipping_phone, '') AS shipping_phone,
       COALESCE(o.shipping_address, '') AS shipping_address,
       COALESCE(o.note, '') AS note,
       o.paid_at,
       o.create_time,
       o.update_time
FROM orders o
LEFT JOIN coupons c ON c.id = o.coupon_id
WHERE o.public_id::text = sqlc.arg(public_id)::text
  AND o.is_delete = false;

-- name: GetOrderByPublicIDForUser :one
SELECT o.id,
       o.public_id::text AS public_id,
       o.user_id,
       COALESCE(c.code, '') AS coupon_code,
       o.subtotal::float8 AS subtotal,
       o.discount_amount::float8 AS discount_amount,
       o.shipping_fee::float8 AS shipping_fee,
       o.total::float8 AS total,
       o.status,
       o.payment_status,
       o.payment_method,
       COALESCE(o.shipping_name, '') AS shipping_name,
       COALESCE(o.shipping_phone, '') AS shipping_phone,
       COALESCE(o.shipping_address, '') AS shipping_address,
       COALESCE(o.note, '') AS note,
       o.paid_at,
       o.create_time,
       o.update_time
FROM orders o
LEFT JOIN coupons c ON c.id = o.coupon_id
WHERE o.public_id::text = sqlc.arg(public_id)::text
  AND o.is_delete = false
  AND o.user_id = (
      SELECT id
      FROM users
      WHERE public_id::text = sqlc.arg(user_public_id)::text
        AND is_delete = false
  );

-- name: GetOrderForUpdate :one
SELECT o.id,
       o.public_id::text AS public_id,
       o.user_id,
       COALESCE(c.code, '') AS coupon_code,
       o.subtotal::float8 AS subtotal,
       o.discount_amount::float8 AS discount_amount,
       o.shipping_fee::float8 AS shipping_fee,
       o.total::float8 AS total,
       o.status,
       o.payment_status,
       o.payment_method,
       COALESCE(o.shipping_name, '') AS shipping_name,
       COALESCE(o.shipping_phone, '') AS shipping_phone,
       COALESCE(o.shipping_address, '') AS shipping_address,
       COALESCE(o.note, '') AS note,
       o.paid_at,
       o.create_time,
       o.update_time
FROM orders o
LEFT JOIN coupons c ON c.id = o.coupon_id
WHERE o.public_id::text = sqlc.arg(public_id)::text
  AND o.is_delete = false
FOR UPDATE OF o;

-- name: GetOrderForUpdateForUser :one
SELECT o.id,
       o.public_id::text AS public_id,
       o.user_id,
       COALESCE(c.code, '') AS coupon_code,
       o.subtotal::float8 AS subtotal,
       o.discount_amount::float8 AS discount_amount,
       o.shipping_fee::float8 AS shipping_fee,
       o.total::float8 AS total,
       o.status,
       o.payment_status,
       o.payment_method,
       COALESCE(o.shipping_name, '') AS shipping_name,
       COALESCE(o.shipping_phone, '') AS shipping_phone,
       COALESCE(o.shipping_address, '') AS shipping_address,
       COALESCE(o.note, '') AS note,
       o.paid_at,
       o.create_time,
       o.update_time
FROM orders o
LEFT JOIN coupons c ON c.id = o.coupon_id
WHERE o.public_id::text = sqlc.arg(public_id)::text
  AND o.is_delete = false
  AND o.user_id = (
      SELECT id
      FROM users
      WHERE public_id::text = sqlc.arg(user_public_id)::text
        AND is_delete = false
  )
FOR UPDATE OF o;

-- name: ListOrderItems :many
SELECT oi.id,
       oi.product_id AS product_internal_id,
       COALESCE(p.public_id::text, '')::text AS product_id,
       oi.product_title,
       oi.product_cover,
       oi.unit_price::float8 AS unit_price,
       oi.quantity,
       oi.subtotal::float8 AS subtotal
FROM order_items oi
LEFT JOIN products p ON p.id = oi.product_id
WHERE oi.order_id = sqlc.arg(order_id)
  AND oi.is_delete = false
ORDER BY oi.id;

-- ============================================================
-- Listing
-- ============================================================

-- name: ListOrdersForUser :many
SELECT o.id,
       o.public_id::text AS public_id,
       o.user_id,
       COALESCE(c.code, '') AS coupon_code,
       o.subtotal::float8 AS subtotal,
       o.discount_amount::float8 AS discount_amount,
       o.shipping_fee::float8 AS shipping_fee,
       o.total::float8 AS total,
       o.status,
       o.payment_status,
       o.payment_method,
       COALESCE(o.shipping_name, '') AS shipping_name,
       COALESCE(o.shipping_phone, '') AS shipping_phone,
       COALESCE(o.shipping_address, '') AS shipping_address,
       COALESCE(o.note, '') AS note,
       o.paid_at,
       o.create_time,
       o.update_time
FROM orders o
LEFT JOIN coupons c ON c.id = o.coupon_id
WHERE o.is_delete = false
  AND o.user_id = (
      SELECT id
      FROM users
      WHERE public_id::text = sqlc.arg(user_public_id)::text
        AND is_delete = false
  )
  AND (sqlc.narg(filter_status)::smallint IS NULL
       OR o.status = sqlc.narg(filter_status)::smallint)
ORDER BY o.create_time DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountOrdersForUser :one
SELECT count(*)
FROM orders o
WHERE o.is_delete = false
  AND o.user_id = (
      SELECT id
      FROM users
      WHERE public_id::text = sqlc.arg(user_public_id)::text
        AND is_delete = false
  )
  AND (sqlc.narg(filter_status)::smallint IS NULL
       OR o.status = sqlc.narg(filter_status)::smallint);

-- name: ListOrders :many
SELECT o.id,
       o.public_id::text AS public_id,
       o.user_id,
       COALESCE(c.code, '') AS coupon_code,
       o.subtotal::float8 AS subtotal,
       o.discount_amount::float8 AS discount_amount,
       o.shipping_fee::float8 AS shipping_fee,
       o.total::float8 AS total,
       o.status,
       o.payment_status,
       o.payment_method,
       COALESCE(o.shipping_name, '') AS shipping_name,
       COALESCE(o.shipping_phone, '') AS shipping_phone,
       COALESCE(o.shipping_address, '') AS shipping_address,
       COALESCE(o.note, '') AS note,
       o.paid_at,
       o.create_time,
       o.update_time
FROM orders o
LEFT JOIN coupons c ON c.id = o.coupon_id
WHERE o.is_delete = false
  AND (sqlc.narg(filter_status)::smallint IS NULL
       OR o.status = sqlc.narg(filter_status)::smallint)
  AND (sqlc.narg(filter_payment_status)::smallint IS NULL
       OR o.payment_status = sqlc.narg(filter_payment_status)::smallint)
ORDER BY o.create_time DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountOrders :one
SELECT count(*)
FROM orders o
WHERE o.is_delete = false
  AND (sqlc.narg(filter_status)::smallint IS NULL
       OR o.status = sqlc.narg(filter_status)::smallint)
  AND (sqlc.narg(filter_payment_status)::smallint IS NULL
       OR o.payment_status = sqlc.narg(filter_payment_status)::smallint);

-- ============================================================
-- Mutations
-- ============================================================

-- name: UpdateOrderStatus :execrows
UPDATE orders
SET status = sqlc.arg(status),
    update_time = now()
WHERE id = sqlc.arg(id)
  AND status = sqlc.arg(expected_status)
  AND is_delete = false;

-- name: UpdateOrderPayment :execrows
UPDATE orders
SET payment_status = sqlc.arg(payment_status),
    paid_at = COALESCE(sqlc.narg(paid_at), paid_at),
    update_time = now()
WHERE id = sqlc.arg(id)
  AND payment_status = sqlc.arg(expected_status)
  AND is_delete = false;
