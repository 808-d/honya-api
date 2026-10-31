-- 000007_audit_triggers.up.sql
-- Capture every row create/update/delete in audit_logs.
--
-- The actor's public id is supplied per transaction with:
--   SELECT set_config('app.actor_public_id', '<public-id>', true);
-- Unauthenticated and system writes leave it empty, so actor_user_id is NULL.

CREATE OR REPLACE FUNCTION audit_row_change() RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    actor_public_id TEXT;
    actor_id        BIGINT;
    old_payload     JSONB;
    new_payload     JSONB;
    entity_id       BIGINT;
    action          TEXT;
BEGIN
    actor_public_id := NULLIF(current_setting('app.actor_public_id', true), '');
    IF actor_public_id IS NOT NULL THEN
        SELECT id
        INTO actor_id
        FROM users
        WHERE public_id::text = actor_public_id
          AND is_delete = false;
    END IF;

    IF TG_OP = 'INSERT' THEN
        old_payload := NULL;
        new_payload := to_jsonb(NEW) - 'hash_password' - 'cover' - 'product_cover';
        entity_id := NULLIF(new_payload ->> 'id', '')::bigint;
        action := 'create';
    ELSIF TG_OP = 'UPDATE' THEN
        old_payload := to_jsonb(OLD) - 'hash_password' - 'cover' - 'product_cover';
        new_payload := to_jsonb(NEW) - 'hash_password' - 'cover' - 'product_cover';
        entity_id := NULLIF(new_payload ->> 'id', '')::bigint;
        IF (old_payload ->> 'is_delete')::boolean IS FALSE
           AND (new_payload ->> 'is_delete')::boolean IS TRUE THEN
            action := 'delete';
        ELSE
            action := 'update';
        END IF;
    ELSE
        old_payload := to_jsonb(OLD) - 'hash_password' - 'cover' - 'product_cover';
        new_payload := NULL;
        entity_id := NULLIF(old_payload ->> 'id', '')::bigint;
        action := 'delete';
    END IF;

    INSERT INTO audit_logs (
        actor_user_id,
        entity_table,
        entity_id,
        action,
        old_data,
        new_data
    ) VALUES (
        actor_id,
        TG_TABLE_NAME,
        entity_id,
        action,
        old_payload,
        new_payload
    );

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$;

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
        EXECUTE format(
            'CREATE TRIGGER %I
             AFTER INSERT OR UPDATE OR DELETE ON %I
             FOR EACH ROW
             EXECUTE FUNCTION audit_row_change()',
            'audit_' || table_name || '_trg',
            table_name
        );
    END LOOP;
END;
$$;
