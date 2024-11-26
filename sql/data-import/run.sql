\set base_path 'sql/data-import'

\i :base_path/000_drop-tables.sql
\i :base_path/001_kind_enum.sql
\i :base_path/002_schema_omdb.sql
\i :base_path/010_import_data.sql
\i :base_path/011_remove_duplicates.sql
-- VACUUM;
-- \i :base_path/020_add_primary_keys.sql
-- \i :base_path/021_clean_orphans.sql

-- \i :base_path/022_add_foreign_keys.sql
-- ANALYZE;
-- \i :base_path/031_purge_dirty_categories.sql

