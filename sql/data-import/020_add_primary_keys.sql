BEGIN;

ALTER TABLE people             ADD PRIMARY KEY (id);
ALTER TABLE people_aliases     ADD PRIMARY KEY (person_id, name);
ALTER TABLE categories         ADD PRIMARY KEY (id);
ALTER TABLE image_ids          ADD PRIMARY KEY (id);
ALTER TABLE image_licenses     ADD PRIMARY KEY (image_id);
ALTER TABLE jobs               ADD PRIMARY KEY (id);
ALTER TABLE movie_categories   ADD PRIMARY KEY (movie_id, category_id);
ALTER TABLE movie_keywords     ADD PRIMARY KEY (movie_id, category_id);
ALTER TABLE movie_links        ADD PRIMARY KEY (movie_id, language, key);
ALTER TABLE movie_references   ADD PRIMARY KEY (movie_id, referenced_id, type);
ALTER TABLE people_links       ADD PRIMARY KEY (person_id, language, key);
ALTER TABLE movies             ADD PRIMARY KEY (id);
ALTER TABLE trailers           ADD PRIMARY KEY (movie_id, trailer_id);
ALTER TABLE casts              ADD PRIMARY KEY (movie_id, person_id, job_id, role, position);

COMMIT;