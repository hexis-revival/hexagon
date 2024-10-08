DROP INDEX idx_relationships_target_id;
DROP INDEX idx_relationships_user_id;

DROP INDEX idx_stats_user_id;
DROP INDEX idx_users_name;
DROP INDEX idx_users_id;

DELETE FROM relationships;
DROP TABLE relationships;

DELETE FROM stats;
DROP TABLE stats;

DELETE FROM users;
DROP TABLE users;