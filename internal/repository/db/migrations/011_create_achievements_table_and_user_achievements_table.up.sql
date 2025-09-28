CREATE TABLE IF NOT EXISTS achievements(
    achievement_id serial PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    icon_url VARCHAR(300) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_achievements(
    user_achievement_id serial PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE NOT NULL,
    achievement_id INT REFERENCES achievements(achievement_id) ON DELETE CASCADE NOT NULL,
    date_earned TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);