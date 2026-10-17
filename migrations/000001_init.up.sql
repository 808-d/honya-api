-- 000001_init.up.sql
-- Core schema for honya: manga / light novel / book e-commerce
--
-- Conventions:
--   * Internal PK: BIGINT GENERATED ALWAYS AS IDENTITY
--   * External ID: public_id UUID on exposed entities
--   * Audit: create_time, update_time
--   * Soft delete: is_delete, delete_at
--   * Money: NUMERIC(14,2)
--   * No PostgreSQL enums; use CHECK + COMMENT for numeric statuses/types
--
-- Audit/soft-delete columns are declared directly in each table.
-- PostgreSQL table inheritance is intentionally not used.

-- ============================================================
-- Shared trigger function
-- ============================================================

CREATE OR REPLACE FUNCTION set_update_time()
RETURNS TRIGGER AS $$
BEGIN
    NEW.update_time := now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================================
-- RBAC: roles <-> features
-- ============================================================

CREATE TABLE roles (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id   UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',

    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete   BOOLEAN NOT NULL DEFAULT false,
    delete_at   TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL))
);

CREATE TABLE features (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id   UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    code        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',

    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete   BOOLEAN NOT NULL DEFAULT false,
    delete_at   TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL))
);

CREATE TABLE role_features (
    role_id    BIGINT NOT NULL
    REFERENCES roles (id) ON DELETE CASCADE,
    feature_id BIGINT NOT NULL
    REFERENCES features (id) ON DELETE CASCADE,

    PRIMARY KEY (role_id, feature_id)
);

-- ============================================================
-- Users
-- ============================================================

CREATE TABLE users (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id     UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    username      TEXT NOT NULL UNIQUE,
    display_name  TEXT NOT NULL DEFAULT '',
    hash_password TEXT NOT NULL,
    email         TEXT NOT NULL UNIQUE,
    phone         TEXT,

    create_time   TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time   TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete     BOOLEAN NOT NULL DEFAULT false,
    delete_at     TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL))
);

CREATE TABLE user_roles (
    user_id BIGINT NOT NULL
    REFERENCES users (id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL
    REFERENCES roles (id) ON DELETE CASCADE,

    PRIMARY KEY (user_id, role_id)
);

-- ============================================================
-- Catalog: authors / publishers / genres / tags
-- ============================================================

CREATE TABLE authors (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id   UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name        TEXT NOT NULL,

    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete   BOOLEAN NOT NULL DEFAULT false,
    delete_at   TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL))
);

CREATE TABLE publishers (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id   UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name        TEXT NOT NULL,

    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete   BOOLEAN NOT NULL DEFAULT false,
    delete_at   TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL))
);

CREATE TABLE genres (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id   UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name        TEXT NOT NULL,

    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete   BOOLEAN NOT NULL DEFAULT false,
    delete_at   TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL))
);

CREATE TABLE tags (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id   UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    name        TEXT NOT NULL UNIQUE,

    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete   BOOLEAN NOT NULL DEFAULT false,
    delete_at   TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL))
);

-- ============================================================
-- Products
--
-- status: 0=draft, 1=active, 2=archived
-- type:   0=manga, 1=light_novel, 2=novel,
--         3=artbook, 4=magazine
-- ============================================================

CREATE TABLE products (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id    UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    title        TEXT NOT NULL,
    cover        TEXT NOT NULL DEFAULT '',
    description  TEXT NOT NULL DEFAULT '',
    publisher_id BIGINT REFERENCES publishers (id),
    publish_time TIMESTAMPTZ,
    stock        INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    price        NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (price >= 0),
    status       SMALLINT NOT NULL DEFAULT 0
    CHECK (status IN (0, 1, 2)),
    type         SMALLINT NOT NULL DEFAULT 0
    CHECK (type BETWEEN 0 AND 4),

    create_time  TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time  TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete    BOOLEAN NOT NULL DEFAULT false,
    delete_at    TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL))
);

COMMENT ON COLUMN products.status IS
'0=draft, 1=active, 2=archived';

COMMENT ON COLUMN products.type IS
'0=manga, 1=light_novel, 2=novel, 3=artbook, 4=magazine';

CREATE TABLE product_authors (
    product_id BIGINT NOT NULL
    REFERENCES products (id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL
    REFERENCES authors (id),

    PRIMARY KEY (product_id, author_id)
);

CREATE TABLE product_genres (
    product_id BIGINT NOT NULL
    REFERENCES products (id) ON DELETE CASCADE,
    genre_id BIGINT NOT NULL
    REFERENCES genres (id),

    PRIMARY KEY (product_id, genre_id)
);

CREATE TABLE product_tags (
    product_id BIGINT NOT NULL
    REFERENCES products (id) ON DELETE CASCADE,
    tag_id BIGINT NOT NULL
    REFERENCES tags (id),

    PRIMARY KEY (product_id, tag_id)
);

-- ============================================================
-- Cart
-- One cart per registered user
-- ============================================================

CREATE TABLE carts (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     BIGINT NOT NULL UNIQUE
    REFERENCES users (id),

    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete   BOOLEAN NOT NULL DEFAULT false,
    delete_at   TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL))
);

CREATE TABLE cart_items (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cart_id     BIGINT NOT NULL
    REFERENCES carts (id) ON DELETE CASCADE,
    product_id  BIGINT NOT NULL
    REFERENCES products (id),
    quantity    INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),

    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete   BOOLEAN NOT NULL DEFAULT false,
    delete_at   TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL)),
    UNIQUE (cart_id, product_id)
);

-- ============================================================
-- Coupons
-- type: fixed | percent
-- ============================================================

CREATE TABLE coupons (
    id               BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id        UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    code             TEXT NOT NULL UNIQUE,
    type SMALLINT NOT NULL DEFAULT 0
    CHECK (type IN (0, 1)),
    value            NUMERIC(14,2) NOT NULL CHECK (value >= 0),
    min_order_amount NUMERIC(14,2) NOT NULL DEFAULT 0
    CHECK (min_order_amount >= 0),
    max_discount     NUMERIC(14,2)
    CHECK (max_discount IS NULL OR max_discount >= 0),
    start_at         TIMESTAMPTZ,
    expire_at        TIMESTAMPTZ,
    usage_limit      INTEGER
    CHECK (usage_limit IS NULL OR usage_limit >= 0),
    used_count       INTEGER NOT NULL DEFAULT 0
    CHECK (used_count >= 0),
    per_user_limit   INTEGER NOT NULL DEFAULT 1
    CHECK (per_user_limit > 0),
    is_active        BOOLEAN NOT NULL DEFAULT true,

    create_time      TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time      TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete        BOOLEAN NOT NULL DEFAULT false,
    delete_at        TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL)),
    CHECK (expire_at IS NULL OR start_at IS NULL OR expire_at > start_at),
    CHECK (type <> 1 OR value <= 100),
    CHECK (usage_limit IS NULL OR used_count <= usage_limit)

);

    COMMENT ON COLUMN coupons.type IS
    '0=fixed, 1=percent';
-- ============================================================
-- Orders
--
-- status:
--   0=pending, 1=confirmed, 2=processing, 3=shipping,
--   4=completed, 5=cancelled, 6=refunded
--
-- payment_status:
--   unpaid | paid | failed | refunded
-- ============================================================

CREATE TABLE orders (
    id               BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id        UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    user_id          BIGINT NOT NULL REFERENCES users (id),
    coupon_id        BIGINT REFERENCES coupons (id),

    subtotal         NUMERIC(14,2) NOT NULL DEFAULT 0
    CHECK (subtotal >= 0),
    discount_amount  NUMERIC(14,2) NOT NULL DEFAULT 0
    CHECK (discount_amount >= 0),
    shipping_fee     NUMERIC(14,2) NOT NULL DEFAULT 0
    CHECK (shipping_fee >= 0),
    total            NUMERIC(14,2) NOT NULL DEFAULT 0
    CHECK (total >= 0),

    status           SMALLINT NOT NULL DEFAULT 0
    CHECK (status BETWEEN 0 AND 6),
    payment_status SMALLINT NOT NULL DEFAULT 0
    CHECK (payment_status BETWEEN 0 AND 3),

    shipping_name    TEXT,
    shipping_phone   TEXT,
    shipping_address TEXT,
    note             TEXT,

    create_time      TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time      TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete        BOOLEAN NOT NULL DEFAULT false,
    delete_at        TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL))
);

COMMENT ON COLUMN orders.payment_status IS
'0=unpaid, 1=paid, 2=failed, 3=refunded';

COMMENT ON COLUMN orders.status IS
'0=pending, 1=confirmed, 2=processing, 3=shipping, 4=completed, 5=cancelled, 6=refunded';

CREATE TABLE order_items (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id      BIGINT NOT NULL
    REFERENCES orders (id),
    product_id    BIGINT
    REFERENCES products (id) ON DELETE SET NULL,
    product_title TEXT NOT NULL,
    unit_price    NUMERIC(14,2) NOT NULL CHECK (unit_price >= 0),
    quantity      INTEGER NOT NULL CHECK (quantity > 0),
    subtotal      NUMERIC(14,2) NOT NULL CHECK (subtotal >= 0),

    create_time   TIMESTAMPTZ NOT NULL DEFAULT now(),
    update_time   TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_delete     BOOLEAN NOT NULL DEFAULT false,
    delete_at     TIMESTAMPTZ,

    CHECK (is_delete = (delete_at IS NOT NULL))
);

-- ============================================================
-- Coupon usage history
-- ============================================================

CREATE TABLE coupon_usages (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    coupon_id   BIGINT NOT NULL REFERENCES coupons (id),
    user_id     BIGINT NOT NULL REFERENCES users (id),
    order_id    BIGINT NOT NULL UNIQUE REFERENCES orders (id),
    create_time TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================
-- Audit logs: append-only history
-- ============================================================

CREATE TABLE audit_logs (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    actor_user_id BIGINT REFERENCES users (id) ON DELETE SET NULL,
    entity_table  TEXT NOT NULL,
    entity_id     BIGINT,
    action        TEXT NOT NULL,
    old_data      JSONB,
    new_data      JSONB,
    create_time   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================
-- Indexes
-- PostgreSQL does not automatically index referencing FK columns.
-- ============================================================

CREATE INDEX idx_role_features_feature_id
ON role_features (feature_id);

CREATE INDEX idx_user_roles_role_id
ON user_roles (role_id);

CREATE INDEX idx_products_publisher_id
ON products (publisher_id);

CREATE INDEX idx_products_status_type
ON products (status, type);

CREATE INDEX idx_product_authors_author_id
ON product_authors (author_id);

CREATE INDEX idx_product_genres_genre_id
ON product_genres (genre_id);

CREATE INDEX idx_product_tags_tag_id
ON product_tags (tag_id);

CREATE INDEX idx_cart_items_product_id
ON cart_items (product_id);

CREATE INDEX idx_orders_user_created
ON orders (user_id, create_time DESC);

CREATE INDEX idx_orders_status
ON orders (status);

CREATE INDEX idx_orders_coupon_id
ON orders (coupon_id);

CREATE INDEX idx_order_items_product_id
ON order_items (product_id);

CREATE INDEX idx_coupon_usages_user_coupon
ON coupon_usages (user_id, coupon_id);

CREATE INDEX idx_audit_logs_actor_created
ON audit_logs (actor_user_id, create_time DESC);

CREATE INDEX idx_audit_logs_entity
ON audit_logs (entity_table, entity_id);

-- ============================================================
-- update_time triggers
-- ============================================================

DO $$
    DECLARE
    t TEXT;
    BEGIN
	FOREACH t IN ARRAY ARRAY[
	'roles',
	'features',
	'users',
	'authors',
	'publishers',
	'genres',
	'tags',
	'products',
	'carts',
	'cart_items',
	'coupons',
	'orders',
	'order_items'
	]
	LOOP
	    EXECUTE format(
		'CREATE TRIGGER %I
		BEFORE UPDATE ON %I
		FOR EACH ROW
		    EXECUTE FUNCTION set_update_time()',
		    'trg_' || t || '_update_time',
		    t
		);
	END LOOP;
	END;
	$$;
