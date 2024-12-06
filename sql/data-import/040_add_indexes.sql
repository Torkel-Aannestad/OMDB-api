BEGIN;
\echo ''
\echo '040_add_indexes'

CREATE INDEX IF NOT EXISTS movies_name_idx ON movies USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS people_links_person_id_idx ON people_links (person_id);
CREATE INDEX IF NOT EXISTS movie_links_movie_id_idx ON movie_links (movie_id);
CREATE INDEX IF NOT EXISTS casts_person_id_idx ON casts (person_id);
CREATE INDEX IF NOT EXISTS casts_movie_id_idx ON casts (movie_id);
CREATE INDEX IF NOT EXISTS movie_keywords_movie_id_idx ON movie_keywords (movie_id);
CREATE INDEX IF NOT EXISTS movie_keywords_category_id_idx ON movie_keywords (category_id);
CREATE INDEX IF NOT EXISTS movie_categories_movie_id_idx ON movie_categories (movie_id);
CREATE INDEX IF NOT EXISTS movie_categories_category_id_idx ON movie_categories (category_id);
CREATE INDEX IF NOT EXISTS trailers_movie_id_idx ON trailers (movie_id);
CREATE INDEX IF NOT EXISTS images_object_id_idx ON images (object_id);
CREATE INDEX IF NOT EXISTS images_object_type_idx ON images (object_type);

COMMIT;