-- 000005_orders_cod.down.sql
-- Remove cash-on-delivery support and order item snapshots.

DROP INDEX idx_orders_created;

ALTER TABLE order_items
    DROP COLUMN product_cover;

ALTER TABLE orders
    DROP COLUMN paid_at,
    DROP COLUMN payment_method;
