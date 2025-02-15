-- Drop indexes
DROP INDEX IF EXISTS idx_word_review_items_activity_id;
DROP INDEX IF EXISTS idx_word_review_items_word_id;
DROP INDEX IF EXISTS idx_study_sessions_activity_id;
DROP INDEX IF EXISTS idx_study_sessions_group_id;
DROP INDEX IF EXISTS idx_word_groups_group_id;
DROP INDEX IF EXISTS idx_word_groups_word_id;
DROP INDEX IF EXISTS idx_words_level;

-- Drop tables
DROP TABLE IF EXISTS word_review_items;
DROP TABLE IF EXISTS study_sessions;
DROP TABLE IF EXISTS study_activity_groups;
DROP TABLE IF EXISTS study_activities;
DROP TABLE IF EXISTS word_groups;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS words;
