BEGIN;
\echo ''
\echo '040_add_indexes'

CREATE INDEX IF NOT EXISTS movies_name_idx ON movies USING GIN (to_tsvector('simple', name));

COMMIT;