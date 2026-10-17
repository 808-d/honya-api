-- 000002_user_verification.down.sql

DELETE FROM roles WHERE name IN ('customer', 'admin');

ALTER TABLE users DROP COLUMN IF EXISTS is_verified;
