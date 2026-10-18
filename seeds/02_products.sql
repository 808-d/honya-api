-- seeds/02_products.sql
-- Dummy products (100 books) with author / genre / tag relations.
--
-- Requires seeds/01_catalog.sql to have been applied first.
-- Skips work if the products table is not empty, so it is safe to re-run.
--
-- Run:
--   docker exec -i honya-db psql -U postgres -d honya -v ON_ERROR_STOP=1 \
--     -f - < seeds/02_products.sql

BEGIN;

DO $$
DECLARE
    cover      text := 'data:image/png;base64,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';
    series     text[] := ARRAY[
        'Crimson Blade', 'Echoes of the Void', 'Sakura Requiem', 'Iron Petal',
        'Midnight Ramen Club', 'The Last Alchemist', 'Starlight Courier',
        'Ghost in the Lantern', 'Rusted Crown', 'Paper Tiger', 'Neon Shrine',
        'Dragon''s Apprentice', 'Silent Blade', 'Café at World''s End',
        'Chrono Detective', 'Ashes of Autumn', 'Sky Whale', 'Mecha Academy',
        'Fox Spirit Kitchen', 'The Wandering Swordsmith', 'Glass Ocean',
        'Lullaby for a Robot', 'Bamboo and Bullets', 'Hollow Moon',
        'Tea House of Spirits', 'Zero Gravity High',
        'The Cartographer''s Daughter', 'Oni Detective Agency',
        'Petals in the Rain', 'Frostbound'
    ];
    type_names text[] := ARRAY['manga', 'light novel', 'novel', 'artbook', 'magazine'];
    ns         int := 0;
    na         int := 0;
    ng         int := 0;
    np         int := 0;
BEGIN
    IF EXISTS (SELECT 1 FROM products) THEN
        RAISE NOTICE 'products already exist; skipping product seed';
        RETURN;
    END IF;

    IF (SELECT count(*) FROM publishers) = 0
       OR (SELECT count(*) FROM authors) = 0
       OR (SELECT count(*) FROM genres) = 0 THEN
        RAISE EXCEPTION 'catalog tables are empty; run seeds/01_catalog.sql first';
    END IF;

    ns := array_length(series, 1);
    na := (SELECT count(*) FROM authors);
    ng := (SELECT count(*) FROM genres);
    np := (SELECT count(*) FROM publishers);

    INSERT INTO products (
        title, cover, description, publisher_id, publish_time, stock, price, status, type
    )
    SELECT
        s.title,
        cover,
        s.title || ' is a ' || type_names[(i % 5) + 1] ||
            ' published by ' || pub.name ||
            ', written by ' || aut.name ||
            ', blending ' || gen.name || ' into a story fans keep returning to.',
        pub.id,
        now() - make_interval(days => (i * 37) % 60 + 30),
        (i * 13) % 250,
        round(7 + ((i * 97) % 3800)::numeric / 100, 2),
        CASE WHEN i % 7 = 0 THEN 0 WHEN i % 11 = 0 THEN 2 ELSE 1 END,
        (i % 5)::smallint
    FROM (
        SELECT
            i,
            series[(i % ns) + 1] || ' Vol. ' || ((i / ns) + 1) AS title
        FROM generate_series(0, 99) AS i
    ) s
    CROSS JOIN LATERAL (
        SELECT p.id, p.name FROM publishers p ORDER BY p.id OFFSET (s.i % np) LIMIT 1
    ) pub
    CROSS JOIN LATERAL (
        SELECT a.name FROM authors a ORDER BY a.id OFFSET ((s.i * 7) % na) LIMIT 1
    ) aut
    CROSS JOIN LATERAL (
        SELECT g.name FROM genres g ORDER BY g.id OFFSET ((s.i * 7) % ng) LIMIT 1
    ) gen;

    INSERT INTO product_authors (product_id, author_id)
    SELECT p.id, a.id
    FROM (SELECT id, row_number() OVER (ORDER BY id) - 1 AS i FROM products) p
    CROSS JOIN LATERAL generate_series(0, p.i % 2) AS sel
    JOIN (
        SELECT id, row_number() OVER (ORDER BY id) - 1 AS j, count(*) OVER () AS n FROM authors
    ) a ON a.j = (p.i * 7 + sel) % a.n;

    INSERT INTO product_genres (product_id, genre_id)
    SELECT p.id, g.id
    FROM (SELECT id, row_number() OVER (ORDER BY id) - 1 AS i FROM products) p
    CROSS JOIN LATERAL generate_series(0, p.i % 3) AS sel
    JOIN (
        SELECT id, row_number() OVER (ORDER BY id) - 1 AS j, count(*) OVER () AS n FROM genres
    ) g ON g.j = (p.i * 5 + sel) % g.n;

    INSERT INTO product_tags (product_id, tag_id)
    SELECT p.id, t.id
    FROM (SELECT id, row_number() OVER (ORDER BY id) - 1 AS i FROM products) p
    CROSS JOIN LATERAL generate_series(0, (p.i % 4) - 1) AS sel
    JOIN (
        SELECT id, row_number() OVER (ORDER BY id) - 1 AS j, count(*) OVER () AS n FROM tags
    ) t ON t.j = (p.i * 3 + sel) % t.n;
END $$;

COMMIT;
