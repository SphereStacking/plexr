-- Create aggregation tables for analytics
CREATE TABLE IF NOT EXISTS daily_event_summary (
    date DATE NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    count BIGINT DEFAULT 0,
    unique_users INTEGER DEFAULT 0,
    PRIMARY KEY (date, event_type)
);

CREATE TABLE IF NOT EXISTS user_activity_summary (
    user_id INTEGER NOT NULL,
    date DATE NOT NULL,
    total_events INTEGER DEFAULT 0,
    page_views INTEGER DEFAULT 0,
    last_seen_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (user_id, date)
);

-- Create indexes for summary tables
CREATE INDEX IF NOT EXISTS idx_daily_event_summary_date ON daily_event_summary(date);
CREATE INDEX IF NOT EXISTS idx_user_activity_summary_date ON user_activity_summary(date);