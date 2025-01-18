DROP TABLE activities;

-- Drop single-column indexes
DROP INDEX IF EXISTS idx_activities_activity_type;
DROP INDEX IF EXISTS idx_activities_done_at;
DROP INDEX IF EXISTS idx_activities_calories_burned;

-- Drop composite indexes
DROP INDEX IF EXISTS idx_activities_composite;
DROP INDEX IF EXISTS idx_activities_ordered;
