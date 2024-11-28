\set base_path 'sql/data-import'

\i :base_path/000_drop-tables.sql
\i :base_path/001_kind_enum.sql
\i :base_path/002_schema_omdb.sql
\i :base_path/010_import_data.sql
\i :base_path/011_remove_duplicates.sql
VACUUM;
\i :base_path/020_add_primary_keys.sql
\i :base_path/021_clean_orphans.sql

\i :base_path/022_add_foreign_keys.sql
\i :base_path/023_add_null_values.sql
\i :base_path/024_add_not_null_constraints.sql
\i :base_path/025_add_identity_for_id.sql

ANALYZE;
\i :base_path/030_add_created_at_column.sql
\i :base_path/031_purge_dirty_categories.sql

