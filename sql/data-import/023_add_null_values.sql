BEGIN;
\echo ''
\echo '023_add_null_values'

--people
UPDATE people
SET
    name = COALESCE(name, 'Unknown'),      
    birthday = COALESCE(birthday, '1888-01-01'), 
    deathday = COALESCE(deathday, '1888-01-01'), 
    gender = COALESCE(gender, 99)
WHERE
    name IS NULL OR
    birthday IS NULL OR
    deathday IS NULL OR
    gender IS NULL;

-- movies
UPDATE movies
SET
    name = COALESCE(name, 'Unknown'),
    -- parent_id nullable
    date = COALESCE(date, '1888-01-01'),
    -- series_id nullable
    kind = COALESCE(kind, 'movie'),
    runtime = COALESCE(runtime, -1),
    budget = COALESCE(budget, -1),
    revenue = COALESCE(revenue, -1),
    homepage = COALESCE(homepage, ''),
    vote_average = COALESCE(vote_average, -1),
    votes_count = COALESCE(votes_count, -1),
    abstract = COALESCE(abstract, '')
WHERE
    name IS NULL
    OR parent_id IS NULL
    OR date IS NULL
    OR series_id IS NULL
    OR kind IS NULL
    OR runtime IS NULL
    OR budget IS NULL
    OR revenue IS NULL
    OR homepage IS NULL
    OR vote_average IS NULL
    OR vote_average IS NULL
    OR abstract IS NULL;


UPDATE image_licenses
SET
    source = COALESCE(source, ''),
    license_id = COALESCE(license_id, 0),
    author = COALESCE(author, '')
WHERE
    source IS NULL OR
    license_id IS NULL OR
    author IS NULL;

UPDATE image_ids
SET
    image_version = COALESCE(image_version, 0)
WHERE
    image_version IS NULL;

COMMIT;