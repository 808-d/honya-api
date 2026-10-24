-- seeds/04_coupons.sql
-- Development coupons: one fixed amount and one percentage.
--
-- Idempotent. Run manually (not executed automatically):
--   docker exec -i honya-db psql -U postgres -d honya -v ON_ERROR_STOP=1 \
--     -f - < seeds/04_coupons.sql

BEGIN;

INSERT INTO coupons (
    code, type, value, min_order_amount, max_discount,
    start_at, expire_at, usage_limit, per_user_limit, is_active
)
SELECT v.code, v.type, v.value, v.min_order_amount, v.max_discount,
       v.start_at, v.expire_at, v.usage_limit, v.per_user_limit, v.is_active
FROM (VALUES
    ('WELCOME10', 0::smallint, 10.00::numeric, 50.00::numeric, NULL::numeric,
     NULL::timestamptz, NULL::timestamptz, 1000::integer, 1::integer, true),
    ('SUMMER20', 1::smallint, 20.00::numeric, 100.00::numeric, 30.00::numeric,
     '2026-06-01T00:00:00Z'::timestamptz, '2026-08-31T23:59:59Z'::timestamptz,
     500::integer, 1::integer, true)
) AS v(code, type, value, min_order_amount, max_discount,
       start_at, expire_at, usage_limit, per_user_limit, is_active)
WHERE NOT EXISTS (
    SELECT 1 FROM coupons c
    WHERE upper(c.code) = upper(v.code) AND c.is_delete = false
);

COMMIT;
