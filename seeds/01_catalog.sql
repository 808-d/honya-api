-- seeds/01_catalog.sql
-- Dummy catalog data (publishers, authors, genres, tags) for local development.
--
-- Idempotent: rows are inserted only when a matching name does not already
-- exist, so it is safe to re-run.
--
-- Run:
--   docker exec -i honya-db psql -U postgres -d honya -v ON_ERROR_STOP=1 \
--     -f - < seeds/01_catalog.sql

BEGIN;

INSERT INTO publishers (name)
SELECT v.name
FROM (VALUES
    ('Shueisha'),
    ('Kodansha'),
    ('Kadokawa Shoten'),
    ('Yen Press'),
    ('Viz Media'),
    ('Seven Seas Entertainment'),
    ('Dark Horse Manga'),
    ('Square Enix'),
    ('Tokyopop'),
    ('Vertical'),
    ('J-Novel Club'),
    ('Futabasha')
) AS v(name)
WHERE NOT EXISTS (SELECT 1 FROM publishers p WHERE p.name = v.name);

INSERT INTO authors (name)
SELECT v.name
FROM (VALUES
    ('Eiichiro Oda'),
    ('Masashi Kishimoto'),
    ('Tite Kubo'),
    ('Hajime Isayama'),
    ('Akira Toriyama'),
    ('Naoki Urasawa'),
    ('Hiromu Arakawa'),
    ('Yoshihiro Togashi'),
    ('Rumiko Takahashi'),
    ('Tsugumi Ohba'),
    ('Kentaro Miura'),
    ('CLAMP'),
    ('Natsuki Takaya'),
    ('Yana Toboso'),
    ('Koyoharu Gotouge'),
    ('Gege Akutami'),
    ('Reki Kawahara'),
    ('Nisio Isin'),
    ('Fuminori Teshima'),
    ('Kugane Maruyama'),
    ('Rifujin na Magonote'),
    ('Yuu Kamiya'),
    ('Isuna Hasekura'),
    ('Natsume Akatsuki'),
    ('Sui Ishida'),
    ('ONE'),
    ('Yusuke Murata'),
    ('Kafka Asagiri'),
    ('Sango Harukawa'),
    ('Tappei Nagatsuki')
) AS v(name)
WHERE NOT EXISTS (SELECT 1 FROM authors a WHERE a.name = v.name);

INSERT INTO genres (name)
SELECT v.name
FROM (VALUES
    ('Action'),
    ('Adventure'),
    ('Comedy'),
    ('Drama'),
    ('Fantasy'),
    ('Horror'),
    ('Mystery'),
    ('Romance'),
    ('Sci-Fi'),
    ('Slice of Life'),
    ('Sports'),
    ('Supernatural'),
    ('Thriller'),
    ('Isekai'),
    ('Mecha'),
    ('Historical'),
    ('Psychological'),
    ('Music')
) AS v(name)
WHERE NOT EXISTS (SELECT 1 FROM genres g WHERE g.name = v.name);

INSERT INTO tags (name)
SELECT v.name
FROM (VALUES
    ('Shounen'),
    ('Seinen'),
    ('Shoujo'),
    ('Josei'),
    ('Ongoing'),
    ('Completed'),
    ('Anime Adaptation'),
    ('Light Novel'),
    ('Best Seller'),
    ('Award Winner'),
    ('New Release'),
    ('Cult Classic'),
    ('Web Novel'),
    ('Full Color')
) AS v(name)
WHERE NOT EXISTS (SELECT 1 FROM tags t WHERE t.name = v.name);

COMMIT;
