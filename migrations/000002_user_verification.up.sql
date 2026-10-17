-- 000002_user_verification.up.sql
-- Email verification flag on users + seeded default roles.

ALTER TABLE users
    ADD COLUMN is_verified BOOLEAN NOT NULL DEFAULT false;

INSERT INTO roles (name, description)
VALUES
    ('customer', 'Default role for self-registered customers'),
    ('admin', 'Full administrative access')
ON CONFLICT (name) DO NOTHING;
