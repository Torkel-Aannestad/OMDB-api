-- +goose Up
ALTER TABLE movie_aliases_iso  ADD PRIMARY KEY (movie_id, name, language, official_translation);
ALTER TABLE people             ADD PRIMARY KEY (id);
ALTER TABLE people_aliases     ADD PRIMARY KEY (person_id, name);
ALTER TABLE category_names     ADD PRIMARY KEY (category_id, language);
ALTER TABLE categories         ADD PRIMARY KEY (id);
ALTER TABLE image_ids          ADD PRIMARY KEY (id);
ALTER TABLE image_licenses     ADD PRIMARY KEY (image_id);
ALTER TABLE job_names          ADD PRIMARY KEY (job_id, language);
ALTER TABLE jobs               ADD PRIMARY KEY (id);
ALTER TABLE movie_abstracts_de ADD PRIMARY KEY (movie_id);
ALTER TABLE movie_abstracts_en ADD PRIMARY KEY (movie_id);
ALTER TABLE movie_abstracts_fr ADD PRIMARY KEY (movie_id);
ALTER TABLE movie_abstracts_es ADD PRIMARY KEY (movie_id);
ALTER TABLE movie_categories   ADD PRIMARY KEY (movie_id, category_id);
ALTER TABLE movie_countries    ADD PRIMARY KEY (movie_id, country);
ALTER TABLE movie_keywords     ADD PRIMARY KEY (movie_id, category_id);
ALTER TABLE movie_languages    ADD PRIMARY KEY (movie_id, language);
ALTER TABLE movie_links        ADD PRIMARY KEY (movie_id, language, key);
ALTER TABLE movie_references   ADD PRIMARY KEY (movie_id, referenced_id, type);
ALTER TABLE people_links       ADD PRIMARY KEY (person_id, language, key);
ALTER TABLE movies             ADD PRIMARY KEY (id);
ALTER TABLE trailers           ADD PRIMARY KEY (movie_id, trailer_id);
ALTER TABLE casts              ADD PRIMARY KEY (movie_id, person_id, job_id, role, position);


-- +goose Down

ALTER TABLE movie_aliases_iso  DROP CONSTRAINT movie_aliases_iso_pkey;
ALTER TABLE people             DROP CONSTRAINT people_pkey;
ALTER TABLE people_aliases     DROP CONSTRAINT people_aliases_pkey;
ALTER TABLE category_names     DROP CONSTRAINT category_names_pkey;
ALTER TABLE categories         DROP CONSTRAINT categories_pkey;
ALTER TABLE image_ids          DROP CONSTRAINT image_ids_pkey;
ALTER TABLE image_licenses     DROP CONSTRAINT image_licenses_pkey;
ALTER TABLE job_names          DROP CONSTRAINT job_names_pkey;
ALTER TABLE jobs               DROP CONSTRAINT jobs_pkey;
ALTER TABLE movie_abstracts_de DROP CONSTRAINT movie_abstracts_de_pkey;
ALTER TABLE movie_abstracts_en DROP CONSTRAINT movie_abstracts_en_pkey;
ALTER TABLE movie_abstracts_fr DROP CONSTRAINT movie_abstracts_fr_pkey;
ALTER TABLE movie_abstracts_es DROP CONSTRAINT movie_abstracts_es_pkey;
ALTER TABLE movie_categories   DROP CONSTRAINT movie_categories_pkey;
ALTER TABLE movie_countries    DROP CONSTRAINT movie_countries_pkey;
ALTER TABLE movie_keywords     DROP CONSTRAINT movie_keywords_pkey;
ALTER TABLE movie_languages    DROP CONSTRAINT movie_languages_pkey;
ALTER TABLE movie_links        DROP CONSTRAINT movie_links_pkey;
ALTER TABLE movie_references   DROP CONSTRAINT movie_references_pkey;
ALTER TABLE people_links       DROP CONSTRAINT people_links_pkey;
ALTER TABLE movies             DROP CONSTRAINT movies_pkey;
ALTER TABLE trailers           DROP CONSTRAINT trailers_pkey;
ALTER TABLE casts              DROP CONSTRAINT casts_pkey;
