DELETE FROM role_features rf
USING roles r, features f
WHERE rf.role_id = r.id
  AND rf.feature_id = f.id
  AND r.name = 'admin'
  AND f.code IN (
      'dashboard.access',
      'catalog.read', 'catalog.write',
      'role.read', 'role.write',
      'user.read', 'user.write'
  );

DELETE FROM features
WHERE code IN (
    'dashboard.access',
    'catalog.read', 'catalog.write',
    'role.read', 'role.write',
    'user.read', 'user.write'
);
