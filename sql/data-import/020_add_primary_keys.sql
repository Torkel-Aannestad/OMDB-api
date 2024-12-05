BEGIN;
\echo ''
\echo '020_add_primary_keys'

ALTER TABLE people             ADD PRIMARY KEY (id);
ALTER TABLE casts              ADD PRIMARY KEY (id); --(movie_id, person_id, job_id, role, position)
ALTER TABLE jobs               ADD PRIMARY KEY (id);
ALTER TABLE movies             ADD PRIMARY KEY (id);
ALTER TABLE categories         ADD PRIMARY KEY (id);
ALTER TABLE movie_categories   ADD PRIMARY KEY (movie_id, category_id);
ALTER TABLE movie_keywords     ADD PRIMARY KEY (movie_id, category_id);
ALTER TABLE image_ids          ADD PRIMARY KEY (id);
ALTER TABLE image_licenses     ADD PRIMARY KEY (image_id);
ALTER TABLE movie_links        ADD PRIMARY KEY (id); --(movie_id, language, key)
ALTER TABLE people_links       ADD PRIMARY KEY (id); --(person_id, language, key)
ALTER TABLE trailers           ADD PRIMARY KEY (id);
-- ALTER TABLE movie_references   ADD PRIMARY KEY (movie_id, referenced_id, type);

COMMIT;