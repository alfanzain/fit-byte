CREATE TABLE IF NOT EXISTS activities (
    id SERIAL PRIMARY KEY,
    activity_type VARCHAR(255) NOT NULL,
    done_at TIMESTAMP NOT NULL,
    duration_in_minutes INT,
    calories_burned INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Single-column indexes
CREATE INDEX IF NOT EXISTS idx_activities_activity_type ON activities (activity_type);
CREATE INDEX IF NOT EXISTS idx_activities_done_at ON activities (done_at);
CREATE INDEX IF NOT EXISTS idx_activities_calories_burned ON activities (calories_burned);
CREATE INDEX IF NOT EXISTS idx_activities_user_id ON activities (user_id);

-- Composite indexes for common filter combinations
CREATE INDEX IF NOT EXISTS idx_activities_composite ON activities (activity_type, done_at, calories_burned);
CREATE INDEX IF NOT EXISTS idx_activities_ordered ON activities (activity_type, done_at);
