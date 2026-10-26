-- 000005_orders_cod.up.sql
-- Cash-on-delivery support and order item snapshots.

ALTER TABLE orders
    ADD COLUMN payment_method SMALLINT NOT NULL DEFAULT 0
    CHECK (payment_method = 0),
    ADD COLUMN paid_at TIMESTAMPTZ;

COMMENT ON COLUMN orders.payment_method IS
'0=cod';

COMMENT ON COLUMN orders.paid_at IS
'When the order was marked paid, if any';

ALTER TABLE order_items
    ADD COLUMN product_cover TEXT NOT NULL DEFAULT '';

CREATE INDEX idx_orders_created
ON orders (create_time DESC);
