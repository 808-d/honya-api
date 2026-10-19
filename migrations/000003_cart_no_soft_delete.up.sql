-- 000003_cart_no_soft_delete.up.sql
-- Carts are permanent (never soft-deleted) and cart items are hard-deleted.
-- Remove the soft-delete columns from both tables.

ALTER TABLE carts DROP CONSTRAINT carts_check;
ALTER TABLE carts DROP COLUMN is_delete;
ALTER TABLE carts DROP COLUMN delete_at;

ALTER TABLE cart_items DROP CONSTRAINT cart_items_check;
ALTER TABLE cart_items DROP COLUMN is_delete;
ALTER TABLE cart_items DROP COLUMN delete_at;
