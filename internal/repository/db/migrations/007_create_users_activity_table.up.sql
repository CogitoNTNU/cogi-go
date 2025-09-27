CREATE TABLE IF NOT EXISTS users_activity (
    activity_id serial PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    activity TIMESTAMP[]
);