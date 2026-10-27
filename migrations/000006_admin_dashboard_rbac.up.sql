-- Dashboard access plus RBAC/catalog feature catalogue.
--
-- Roles are free-form labels; access is decided by feature codes. This keeps
-- the seed self-contained (it also ensures the order/coupon codes exist) and
-- grants everything to the bootstrap "admin" role.

INSERT INTO features (code, name, description)
VALUES
    ('dashboard.access', 'Access dashboard',  'Enter the admin dashboard'),
    ('catalog.read',     'Read catalog',      'View products, authors, genres, publishers and tags'),
    ('catalog.write',    'Manage catalog',    'Create, update and delete catalog records'),
    ('role.read',        'Read roles',        'View roles and the feature catalogue'),
    ('role.write',       'Manage roles',      'Create, update and delete roles and assign features'),
    ('user.read',        'Read users',        'View users and their roles'),
    ('user.write',       'Manage user roles', 'Assign roles to users'),
    ('order.read',       'Read orders',       'View orders and their details'),
    ('order.write',      'Manage orders',     'Update order status and payment'),
    ('coupon.read',      'Read coupons',      'View coupons and their usage history'),
    ('coupon.write',     'Manage coupons',    'Create, update and delete coupons')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_features (role_id, feature_id)
SELECT r.id, f.id
FROM roles r
JOIN features f
  ON f.code IN (
      'dashboard.access',
      'catalog.read', 'catalog.write',
      'role.read', 'role.write',
      'user.read', 'user.write',
      'order.read', 'order.write',
      'coupon.read', 'coupon.write'
  )
WHERE r.name = 'admin'
  AND r.is_delete = false
ON CONFLICT DO NOTHING;
