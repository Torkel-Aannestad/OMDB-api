-- +goose Up
BEGIN;

-- movie references
ALTER TABLE movies             ADD FOREIGN KEY (parent_id) REFERENCES movies (id) ON DELETE cascade,
                               ADD FOREIGN KEY (series_id) REFERENCES movies (id) ON DELETE cascade;
ALTER TABLE people_aliases     ADD FOREIGN KEY (person_id) REFERENCES people (id) ON DELETE cascade;
ALTER TABLE people_links       ADD FOREIGN KEY (person_id) REFERENCES people (id) ON DELETE cascade;
ALTER TABLE casts              ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade,
                               ADD FOREIGN KEY (person_id) REFERENCES people (id) ON DELETE cascade,
                               ADD FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE cascade;
ALTER TABLE movie_categories   ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade,
                               ADD FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE cascade;
ALTER TABLE movie_keywords     ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade,
                               ADD FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE cascade;
ALTER TABLE trailers           ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade;
ALTER TABLE movie_links        ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade;
ALTER TABLE movie_aliases_iso  ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade;
ALTER TABLE movie_languages    ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade;
ALTER TABLE movie_countries    ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade;
ALTER TABLE movie_references   ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade,
                               ADD FOREIGN KEY (referenced_id) REFERENCES movies (id) ON DELETE cascade;
ALTER TABLE movie_abstracts_de ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade;
ALTER TABLE movie_abstracts_en ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade;
ALTER TABLE movie_abstracts_fr ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade;
ALTER TABLE movie_abstracts_es ADD FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE cascade;

-- other references
ALTER TABLE image_licenses     ADD FOREIGN KEY (image_id) REFERENCES image_ids (id) ON DELETE cascade;
ALTER TABLE job_names          ADD FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE cascade;
ALTER TABLE categories         ADD FOREIGN KEY (parent_id) REFERENCES categories (id) ON DELETE cascade,
                               ADD FOREIGN KEY (root_id) REFERENCES categories (id) ON DELETE cascade;
ALTER TABLE category_names     ADD FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE cascade;

COMMIT;


-- +goose Down
BEGIN;

ALTER TABLE movies             DROP CONSTRAINT IF EXISTS movies_parent_id_fkey,
                               DROP CONSTRAINT IF EXISTS movies_series_id_fkey;
ALTER TABLE people_aliases     DROP CONSTRAINT IF EXISTS people_aliases_person_id_fkey;
ALTER TABLE people_links       DROP CONSTRAINT IF EXISTS people_links_person_id_fkey;
ALTER TABLE casts              DROP CONSTRAINT IF EXISTS casts_movie_id_fkey,
                               DROP CONSTRAINT IF EXISTS casts_person_id_fkey,
                               DROP CONSTRAINT IF EXISTS casts_job_id_fkey;
ALTER TABLE movie_categories   DROP CONSTRAINT IF EXISTS movie_categories_movie_id_fkey,
                               DROP CONSTRAINT IF EXISTS movie_categories_category_id_fkey;
ALTER TABLE movie_keywords     DROP CONSTRAINT IF EXISTS movie_keywords_movie_id_fkey,
                               DROP CONSTRAINT IF EXISTS movie_keywords_category_id_fkey;
ALTER TABLE trailers           DROP CONSTRAINT IF EXISTS trailers_movie_id_fkey;
ALTER TABLE movie_links        DROP CONSTRAINT IF EXISTS movie_links_movie_id_fkey;
ALTER TABLE movie_aliases_iso  DROP CONSTRAINT IF EXISTS movie_aliases_iso_movie_id_fkey;
ALTER TABLE movie_languages    DROP CONSTRAINT IF EXISTS movie_languages_movie_id_fkey;
ALTER TABLE movie_countries    DROP CONSTRAINT IF EXISTS movie_countries_movie_id_fkey;
ALTER TABLE movie_references   DROP CONSTRAINT IF EXISTS movie_references_movie_id_fkey,
                               DROP CONSTRAINT IF EXISTS movie_references_referenced_id_fkey;
ALTER TABLE movie_abstracts_de DROP CONSTRAINT IF EXISTS movie_abstracts_de_movie_id_fkey;
ALTER TABLE movie_abstracts_en DROP CONSTRAINT IF EXISTS movie_abstracts_en_movie_id_fkey;
ALTER TABLE movie_abstracts_fr DROP CONSTRAINT IF EXISTS movie_abstracts_fr_movie_id_fkey;
ALTER TABLE movie_abstracts_es DROP CONSTRAINT IF EXISTS movie_abstracts_es_movie_id_fkey;

ALTER TABLE image_licenses     DROP CONSTRAINT IF EXISTS image_licenses_image_id_fkey;
ALTER TABLE job_names          DROP CONSTRAINT IF EXISTS job_names_job_id_fkey;
ALTER TABLE categories         DROP CONSTRAINT IF EXISTS categories_parent_id_fkey,
                               DROP CONSTRAINT IF EXISTS categories_root_id_fkey;
ALTER TABLE category_names     DROP CONSTRAINT IF EXISTS category_names_category_id_fkey;

COMMIT;