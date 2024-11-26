BEGIN;

WITH t as (
  select ctid, row_number() over (partition by id), * from movies
)
DELETE FROM movies WHERE ctid in (select ctid FROM t WHERE row_number > 1);

WITH t as (
  select ctid, row_number() over (partition by person_id, name), * from people_aliases
  )
DELETE from people_aliases where ctid in ( select ctid FROM t WHERE row_number > 1);

WITH t as (
  select ctid, row_number() over (partition by movie_id, language), * from movie_links
  )
DELETE from movie_links where ctid in (select ctid FROM t WHERE row_number > 1);

DELETE FROM movie_references WHERE type IS NULL;

WITH t as (
  select ctid, row_number() over (partition by movie_id, referenced_id, type), * from movie_references
  )
DELETE from movie_references where ctid in (select ctid FROM t WHERE row_number > 1);

WITH t as (
  select ctid, row_number() over (partition by movie_id, person_id, job_id, role, "position"), * from casts
  )
DELETE from casts where ctid in (select ctid FROM t WHERE row_number > 1);

WITH t as (
  select ctid, row_number() over (partition by id), * from jobs
  )
DELETE from jobs where ctid in (select ctid FROM t WHERE row_number > 1);

COMMIT;