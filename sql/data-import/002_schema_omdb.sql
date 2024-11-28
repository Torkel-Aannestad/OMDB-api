BEGIN;
\echo ''
\echo '002_schema_omdb'

CREATE TABLE IF NOT EXISTS movies(
    id bigint, name text, parent_id bigint, date date, series_id bigint, kind kind, -- from all_*
	runtime int, budget DOUBLE PRECISION, revenue DOUBLE PRECISION, homepage text, -- from movie_details
	vote_average DOUBLE PRECISION, votes_count bigint, -- from votes
    abstract text -- from movie_abstracts_en
);
CREATE TABLE IF NOT EXISTS people (id bigint, name text, birthday date, deathday date, gender int, aliases text[]);
--  CREATE TABLE IF NOT EXISTS people_aliases (person_id bigint, name text);
CREATE TABLE IF NOT EXISTS people_links (source text, key text, person_id bigint, language text);
CREATE TABLE IF NOT EXISTS casts (movie_id bigint, person_id bigint, job_id bigint, role text, position int);
CREATE TABLE IF NOT EXISTS jobs (id bigint, name text);
CREATE TABLE IF NOT EXISTS movie_categories (movie_id bigint, category_id bigint);
CREATE TABLE IF NOT EXISTS movie_keywords (movie_id bigint, category_id bigint);
CREATE TABLE IF NOT EXISTS categories (id bigint, parent_id bigint, root_id bigint, name text);
CREATE TABLE IF NOT EXISTS trailers (trailer_id bigint, key text, movie_id bigint, language text, source text);
COMMENT ON TABLE trailers is 'Youtube/Vimeo Trailer';
CREATE TABLE IF NOT EXISTS movie_links (source text, key text, movie_id bigint, language text);
CREATE TABLE IF NOT EXISTS image_ids (id bigint, object_id bigint, object_type text, image_version int);
CREATE TABLE IF NOT EXISTS image_licenses (image_id bigint, source text, license_id bigint, author text);
CREATE TABLE IF NOT EXISTS movie_references (movie_id bigint, referenced_id bigint, type text);

COMMIT;