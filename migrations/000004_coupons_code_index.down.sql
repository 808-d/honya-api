-- 000004_coupons_code_index.down.sql
-- Restore the original case-sensitive, delete-agnostic uniqueness.

DROP INDEX ux_coupons_code_active;

ALTER TABLE coupons ADD CONSTRAINT coupons_code_key UNIQUE (code);
