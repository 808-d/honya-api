-- 000001_init.down.sql
-- Drop dependent tables first, then referenced tables.

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS coupon_usages;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS coupons;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;

DROP TABLE IF EXISTS product_tags;
DROP TABLE IF EXISTS product_genres;
DROP TABLE IF EXISTS product_authors;
DROP TABLE IF EXISTS products;

DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_features;

DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS publishers;
DROP TABLE IF EXISTS authors;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS features;
DROP TABLE IF EXISTS roles;

DROP FUNCTION IF EXISTS set_update_time();
