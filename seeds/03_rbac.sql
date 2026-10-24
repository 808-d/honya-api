-- seeds/03_rbac.sql
-- Base feature catalogue for coupons and its grant to the admin role.
--
-- Idempotent. Run manually (not executed automatically):
--   docker exec -i honya-db psql -U postgres -d honya -v ON_ERROR_STOP=1 \
--     -f - < seeds/03_rbac.sql

BEGIN;

INSERT INTO features (code, name, description)
VALUES
    ('coupon.read',  'Read coupons',   'View coupons and their usage history'),
    ('coupon.write', 'Manage coupons', 'Create, update and delete coupons')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_features (role_id, feature_id)
SELECT r.id, f.id
FROM roles r
JOIN features f ON f.code IN ('coupon.read', 'coupon.write')
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

COMMIT;
