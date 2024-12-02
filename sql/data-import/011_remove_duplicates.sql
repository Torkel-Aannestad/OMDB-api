BEGIN;
\echo ''
\echo '011_remove_duplicates'

WITH t as (
  select ctid, row_number() over (partition by id), * from movies
)
DELETE FROM movies WHERE ctid in (select ctid FROM t WHERE row_number > 1);

UPDATE people
SET aliases = (
  SELECT array_agg(DISTINCT alias) 
  FROM unnest(aliases) alias
) WHERE aliases IS NOT NULL;

WITH t as (
  select ctid, row_number() over (partition by movie_id, language), * from movie_links
  )
DELETE from movie_links where ctid in (select ctid FROM t WHERE row_number > 1);

-- DELETE FROM movie_references WHERE type IS NULL;
-- WITH t as (
--   select ctid, row_number() over (partition by movie_id, referenced_id, type), * from movie_references
--   )
-- DELETE from movie_references where ctid in (select ctid FROM t WHERE row_number > 1);

WITH t as (
  select ctid, row_number() over (partition by movie_id, person_id, job_id, role, "position"), * from casts
  )
DELETE from casts where ctid in (select ctid FROM t WHERE row_number > 1);

WITH t as (
  select ctid, row_number() over (partition by id), * from jobs
  )
DELETE from jobs where ctid in (select ctid FROM t WHERE row_number > 1);

COMMIT;