-- ============================================================
-- Coupons
-- ============================================================

-- name: CreateCoupon :one
INSERT INTO coupons (
    code, type, value, min_order_amount, max_discount,
    start_at, expire_at, usage_limit, per_user_limit, is_active
)
VALUES (
    sqlc.arg(code),
    sqlc.arg(type),
    sqlc.arg(value)::float8,
    sqlc.arg(min_order_amount)::float8,
    sqlc.narg(max_discount)::float8,
    sqlc.narg(start_at),
    sqlc.narg(expire_at),
    sqlc.narg(usage_limit),
    sqlc.arg(per_user_limit),
    sqlc.arg(is_active)
)
RETURNING id, public_id::text AS public_id, code, type,
          value::float8 AS value,
          min_order_amount::float8 AS min_order_amount,
          max_discount,
          start_at, expire_at, usage_limit, used_count,
          per_user_limit, is_active, create_time, update_time;

-- name: UpdateCoupon :one
UPDATE coupons
SET code             = sqlc.arg(code),
    type             = sqlc.arg(type),
    value            = sqlc.arg(value)::float8,
    min_order_amount = sqlc.arg(min_order_amount)::float8,
    max_discount     = sqlc.narg(max_discount)::float8,
    start_at         = sqlc.narg(start_at),
    expire_at        = sqlc.narg(expire_at),
    usage_limit      = sqlc.narg(usage_limit),
    per_user_limit   = sqlc.arg(per_user_limit),
    is_active        = sqlc.arg(is_active)
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false
RETURNING id, public_id::text AS public_id, code, type,
          value::float8 AS value,
          min_order_amount::float8 AS min_order_amount,
          max_discount,
          start_at, expire_at, usage_limit, used_count,
          per_user_limit, is_active, create_time, update_time;

-- name: GetCouponByPublicID :one
SELECT id, public_id::text AS public_id, code, type,
       value::float8 AS value,
       min_order_amount::float8 AS min_order_amount,
       max_discount,
       start_at, expire_at, usage_limit, used_count,
       per_user_limit, is_active, create_time, update_time
FROM coupons
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;

-- name: GetCouponByCode :one
SELECT id, public_id::text AS public_id, code, type,
       value::float8 AS value,
       min_order_amount::float8 AS min_order_amount,
       max_discount,
       start_at, expire_at, usage_limit, used_count,
       per_user_limit, is_active, create_time, update_time
FROM coupons
WHERE upper(code) = upper(sqlc.arg(code)::text)
  AND is_delete = false;

-- name: GetCouponByCodeForUpdate :one
SELECT id, public_id::text AS public_id, code, type,
       value::float8 AS value,
       min_order_amount::float8 AS min_order_amount,
       max_discount,
       start_at, expire_at, usage_limit, used_count,
       per_user_limit, is_active, create_time, update_time
FROM coupons
WHERE upper(code) = upper(sqlc.arg(code)::text)
  AND is_delete = false
FOR UPDATE;

-- name: ListCoupons :many
SELECT id, public_id::text AS public_id, code, type,
       value::float8 AS value,
       min_order_amount::float8 AS min_order_amount,
       max_discount,
       start_at, expire_at, usage_limit, used_count,
       per_user_limit, is_active, create_time, update_time
FROM coupons
WHERE is_delete = false
  AND (sqlc.narg(search)::text IS NULL OR code ILIKE '%' || sqlc.narg(search)::text || '%')
  AND (sqlc.narg(filter_type)::smallint IS NULL OR type = sqlc.narg(filter_type)::smallint)
  AND (sqlc.narg(filter_active)::boolean IS NULL OR is_active = sqlc.narg(filter_active)::boolean)
  AND (NOT sqlc.arg(active_now)::boolean
       OR (is_active = true
           AND (start_at IS NULL OR start_at <= now())
           AND (expire_at IS NULL OR expire_at >= now())))
ORDER BY create_time DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountCoupons :one
SELECT count(*)
FROM coupons
WHERE is_delete = false
  AND (sqlc.narg(search)::text IS NULL OR code ILIKE '%' || sqlc.narg(search)::text || '%')
  AND (sqlc.narg(filter_type)::smallint IS NULL OR type = sqlc.narg(filter_type)::smallint)
  AND (sqlc.narg(filter_active)::boolean IS NULL OR is_active = sqlc.narg(filter_active)::boolean)
  AND (NOT sqlc.arg(active_now)::boolean
       OR (is_active = true
           AND (start_at IS NULL OR start_at <= now())
           AND (expire_at IS NULL OR expire_at >= now())));

-- name: SoftDeleteCoupon :execrows
UPDATE coupons
SET is_delete = true, delete_at = now()
WHERE public_id::text = sqlc.arg(public_id)::text
  AND is_delete = false;

-- ============================================================
-- Redemption
-- ============================================================

-- name: CountCouponUsagesByUser :one
SELECT count(*)
FROM coupon_usages
WHERE coupon_id = sqlc.arg(coupon_id)
  AND user_id = sqlc.arg(user_id);

-- name: IncrementCouponUsed :execrows
UPDATE coupons
SET used_count = used_count + 1
WHERE id = sqlc.arg(id)
  AND is_delete = false
  AND (usage_limit IS NULL OR used_count < usage_limit);

-- name: DecrementCouponUsed :exec
UPDATE coupons
SET used_count = used_count - 1
WHERE id = sqlc.arg(id)
  AND used_count > 0;

-- name: InsertCouponUsage :exec
INSERT INTO coupon_usages (coupon_id, user_id, order_id)
VALUES (sqlc.arg(coupon_id), sqlc.arg(user_id), sqlc.arg(order_id));

-- name: DeleteCouponUsageByOrder :one
DELETE FROM coupon_usages
WHERE order_id = sqlc.arg(order_id)
RETURNING coupon_id;

-- name: ListCouponUsages :many
SELECT cu.id,
       u.public_id::text AS user_id,
       u.username,
       o.public_id::text AS order_id,
       cu.create_time
FROM coupon_usages cu
JOIN users u  ON u.id = cu.user_id
JOIN orders o ON o.id = cu.order_id
WHERE cu.coupon_id = sqlc.arg(coupon_id)
ORDER BY cu.create_time DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountCouponUsages :one
SELECT count(*)
FROM coupon_usages
WHERE coupon_id = sqlc.arg(coupon_id);
