BEGIN;
\echo ''
\echo '001_kind_enum'

CREATE TYPE kind AS ENUM (
'movie',
'series',
'season',
'episode',
'movieseries'
);


COMMIT; 