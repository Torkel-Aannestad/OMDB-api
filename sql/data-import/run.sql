\i 0001_import_data.sql
\i 0002_remove_duplicates.sql
\i 0003_clean_orphans.sql
VACUUM;

ANALYZE;
\i 0004_default_english_cateogry.sql
\i 0005_purge_dirty_categories.sql


