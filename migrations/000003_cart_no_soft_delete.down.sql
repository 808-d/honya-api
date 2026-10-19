-- 000003_cart_no_soft_delete.down.sql
-- Restore the soft-delete columns removed by the up migration.

ALTER TABLE carts ADD COLUMN is_delete BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE carts ADD COLUMN delete_at TIMESTAMPTZ;
ALTER TABLE carts ADD CONSTRAINT carts_check CHECK (is_delete = (delete_at IS NOT NULL));

ALTER TABLE cart_items ADD COLUMN is_delete BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE cart_items ADD COLUMN delete_at TIMESTAMPTZ;
ALTER TABLE cart_items ADD CONSTRAINT cart_items_check CHECK (is_delete = (delete_at IS NOT NULL));
