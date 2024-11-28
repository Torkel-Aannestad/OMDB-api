BEGIN; 
\echo ''
\echo '024_add_not_null_constraints'

ALTER TABLE movies
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN date SET NOT NULL,
    ALTER COLUMN kind SET NOT NULL, 
    ALTER COLUMN runtime SET NOT NULL,
    ALTER COLUMN budget SET NOT NULL,
    ALTER COLUMN revenue SET NOT NULL,
    ALTER COLUMN homepage SET NOT NULL,
    ALTER COLUMN vote_average SET NOT NULL,
    ALTER COLUMN votes_count SET NOT NULL,
    ALTER COLUMN abstract SET NOT NULL;

ALTER TABLE people
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN birthday SET NOT NULL,
    ALTER COLUMN deathday SET NOT NULL,
    ALTER COLUMN gender SET NOT NULL;

ALTER TABLE people_aliases
    ALTER COLUMN person_id SET NOT NULL,
    ALTER COLUMN name SET NOT NULL;

ALTER TABLE people_links
    ALTER COLUMN source SET NOT NULL,
    ALTER COLUMN key SET NOT NULL,
    ALTER COLUMN person_id SET NOT NULL,
    ALTER COLUMN language SET NOT NULL;

ALTER TABLE casts
    ALTER COLUMN movie_id SET NOT NULL,
    ALTER COLUMN person_id SET NOT NULL,
    ALTER COLUMN job_id SET NOT NULL,
    ALTER COLUMN role SET NOT NULL,
    ALTER COLUMN position SET NOT NULL;

ALTER TABLE jobs
    ALTER COLUMN name SET NOT NULL;

ALTER TABLE movie_categories
    ALTER COLUMN movie_id SET NOT NULL,
    ALTER COLUMN category_id SET NOT NULL;

ALTER TABLE movie_keywords
    ALTER COLUMN movie_id SET NOT NULL,
    ALTER COLUMN category_id SET NOT NULL;

ALTER TABLE categories
    -- parent_id nullable
    -- root_id nullable
    ALTER COLUMN name SET NOT NULL;

ALTER TABLE trailers
    ALTER COLUMN trailer_id SET NOT NULL,
    ALTER COLUMN key SET NOT NULL,
    ALTER COLUMN movie_id SET NOT NULL,
    ALTER COLUMN language SET NOT NULL,
    ALTER COLUMN source SET NOT NULL;

ALTER TABLE movie_links
    ALTER COLUMN source SET NOT NULL,
    ALTER COLUMN key SET NOT NULL,
    ALTER COLUMN movie_id SET NOT NULL,
    ALTER COLUMN language SET NOT NULL;

ALTER TABLE image_ids
    ALTER COLUMN object_id SET NOT NULL,
    ALTER COLUMN object_type SET NOT NULL,
    ALTER COLUMN image_version SET NOT NULL;

ALTER TABLE image_licenses
    ALTER COLUMN image_id SET NOT NULL,
    ALTER COLUMN source SET NOT NULL,
    ALTER COLUMN license_id SET NOT NULL,
    ALTER COLUMN author SET NOT NULL;

ALTER TABLE movie_references
    ALTER COLUMN movie_id SET NOT NULL,
    ALTER COLUMN referenced_id SET NOT NULL,
    ALTER COLUMN type SET NOT NULL;

COMMIT;