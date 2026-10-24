-- 000004_coupons_code_index.up.sql
-- Coupon codes are case-insensitive and reusable once a coupon is soft-deleted.
-- Replace the plain UNIQUE constraint with a partial, case-insensitive index.

ALTER TABLE coupons DROP CONSTRAINT coupons_code_key;

CREATE UNIQUE INDEX ux_coupons_code_active
    ON coupons (upper(code))
    WHERE is_delete = false;
