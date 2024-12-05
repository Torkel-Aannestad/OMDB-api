BEGIN;
\echo ''
\echo '029_add_alter_columns'

ALTER TABLE movies
    ADD version integer NOT NULL DEFAULT 1,             
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE people
    ADD version integer NOT NULL DEFAULT 1,
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE categories 
    ADD version integer NOT NULL DEFAULT 1,            
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE image_ids RENAME TO images;
ALTER TABLE images          
    RENAME COLUMN image_version TO version;
ALTER TABLE images
    ALTER COLUMN version SET DEFAULT 1,
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
UPDATE images
    SET version = 1
    WHERE version = '0';   

ALTER TABLE image_licenses    
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE jobs
    ADD version integer NOT NULL DEFAULT 1,                 
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE movie_categories   
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE movie_keywords     
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE movie_links    
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE people_links   
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE trailers           
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE casts
    ADD version integer NOT NULL DEFAULT 1,    
    ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();          

-- ALTER TABLE movie_references   
--     ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
--     ADD modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
COMMIT;