-- 000007_audit_triggers.down.sql
-- Remove audit triggers and the shared trigger function.

DO $$
DECLARE
    table_name TEXT;
BEGIN
    FOREACH table_name IN ARRAY ARRAY[
        'features',
        'roles',
        'role_features',
        'users',
        'user_roles',
        'authors',
        'publishers',
        'genres',
        'tags',
        'products',
        'product_authors',
        'product_genres',
        'product_tags',
        'carts',
        'cart_items',
        'coupons',
        'coupon_usages',
        'orders',
        'order_items'
    ]
    LOOP
        EXECUTE format(
            'DROP TRIGGER IF EXISTS %I ON %I',
            'audit_' || table_name || '_trg',
            table_name
        );
    END LOOP;
END;
$$;

DROP FUNCTION IF EXISTS audit_row_change();
