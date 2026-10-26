-- seeds/05_orders_rbac.sql
-- Order features and their grant to the admin role.
--
-- Idempotent. Run manually (not executed automatically):
--   docker exec -i honya-db psql -U postgres -d honya -v ON_ERROR_STOP=1 \
--     -f - < seeds/05_orders_rbac.sql

BEGIN;

INSERT INTO features (code, name, description)
VALUES
    ('order.read',  'Read orders',   'View orders and their details'),
    ('order.write', 'Manage orders', 'Update order status and payment')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_features (role_id, feature_id)
SELECT r.id, f.id
FROM roles r
JOIN features f ON f.code IN ('order.read', 'order.write')
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

COMMIT;
