BEGIN;

UPDATE movies c SET parent_id = NULL WHERE NOT EXISTS (SELECT * FROM movies p WHERE p.id = c.parent_id);

DELETE FROM people_aliases     c WHERE NOT EXISTS (SELECT * FROM people p     WHERE p.id = c.person_id);
DELETE FROM people_links       c WHERE NOT EXISTS (SELECT * FROM people p     WHERE p.id = c.person_id);
DELETE FROM casts              c WHERE NOT EXISTS (SELECT * FROM movies p     WHERE p.id = c.movie_id);
DELETE FROM movie_categories   c WHERE NOT EXISTS (SELECT * FROM movies p     WHERE p.id = c.movie_id);
DELETE FROM movie_categories   c WHERE NOT EXISTS (SELECT * FROM categories p WHERE p.id = c.category_id);
DELETE FROM movie_keywords     c WHERE NOT EXISTS (SELECT * FROM movies p     WHERE p.id = c.movie_id);
DELETE FROM movie_keywords     c WHERE NOT EXISTS (SELECT * FROM categories p WHERE p.id = c.category_id);
DELETE FROM trailers           c WHERE NOT EXISTS (SELECT * FROM movies p     WHERE p.id = c.movie_id);
DELETE FROM movie_links        c WHERE NOT EXISTS (SELECT * FROM movies p     WHERE p.id = c.movie_id);
DELETE FROM movie_references   c WHERE NOT EXISTS (SELECT * FROM movies p     WHERE p.id = c.movie_id);
DELETE FROM movie_references   c WHERE NOT EXISTS (SELECT * FROM movies p     WHERE p.id = c.movie_id);
DELETE FROM movie_references   c WHERE NOT EXISTS (SELECT * FROM movies p     WHERE p.id = c.referenced_id);

DELETE FROM image_ids WHERE object_id IS NULL;
DELETE FROM image_licenses c WHERE NOT EXISTS (SELECT * FROM image_ids p     WHERE p.id = c.image_id);

COMMIT;