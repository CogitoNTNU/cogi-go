CREATE TABLE IF NOT EXISTS users_activity (
    activity_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE NOT NULL,
    activity TIMESTAMP[]
);