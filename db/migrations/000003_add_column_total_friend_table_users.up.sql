ALTER TABLE
    users
ADD
    COLUMN IF NOT EXISTS total_friend INTEGER DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_users_total_friend ON users (total_friend);