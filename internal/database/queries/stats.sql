-- ============================================================
-- Dashboard stats
-- ============================================================

-- name: ProductStats :one
SELECT count(*)::bigint AS total,
       count(*) FILTER (WHERE status = 1)::bigint AS active,
       count(*) FILTER (WHERE status = 0)::bigint AS draft,
       count(*) FILTER (WHERE status = 2)::bigint AS archived,
       count(*) FILTER (WHERE stock <= sqlc.arg(low_stock)::int)::bigint AS low_stock
FROM products
WHERE is_delete = false;

-- name: OrderStats :one
SELECT count(*)::bigint AS total,
       COALESCE(sum(total) FILTER (WHERE payment_status = 1), 0)::float8 AS paid_revenue
FROM orders
WHERE is_delete = false;

-- name: OrderStatusCounts :many
SELECT status, count(*)::bigint AS count
FROM orders
WHERE is_delete = false
GROUP BY status;

-- name: CouponStats :one
SELECT count(*)::bigint AS total,
       count(*) FILTER (
           WHERE is_active
             AND (start_at IS NULL OR start_at <= now())
             AND (expire_at IS NULL OR expire_at >= now())
       )::bigint AS active
FROM coupons
WHERE is_delete = false;

-- name: UserStats :one
SELECT count(*)::bigint AS total
FROM users
WHERE is_delete = false;
