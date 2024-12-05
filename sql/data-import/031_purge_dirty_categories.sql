BEGIN;
\echo ''
\echo '031_purge_dirty_categories'

CREATE TEMP TABLE dirty_categories (category_id bigint PRIMARY KEY);
INSERT INTO dirty_categories
        SELECT id FROM categories WHERE name LIKE 'Erotic%' OR name = 'Sex';
ANALYZE dirty_categories;

CREATE INDEX t1 ON movies(id);
CREATE INDEX t2 ON movies(parent_id);
CREATE INDEX t3 ON movies(series_id);
CREATE INDEX t4 ON movie_categories(movie_id);
CREATE INDEX t5 ON movie_categories(category_id);
CREATE INDEX t6 ON movie_keywords(movie_id);
CREATE INDEX t7 ON movie_keywords(category_id);
CREATE INDEX t8 ON dirty_categories(category_id);

DELETE FROM movies m USING movie_keywords k JOIN dirty_categories d ON k.category_id = d.category_id WHERE m.id = k.movie_id;
DELETE FROM movies m USING movie_categories k JOIN dirty_categories d ON k.category_id = d.category_id WHERE m.id = k.movie_id;

DROP INDEX t1;
DROP INDEX t2;
DROP INDEX t3;
DROP INDEX t4;
DROP INDEX t5;
DROP INDEX t6;
DROP INDEX t7;
DROP INDEX t8;

COMMIT;