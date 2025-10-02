CREATE TABLE IF NOT EXISTS achievements(
    achievement_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(100) NOT NULL,
    icon_url VARCHAR(300) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_achievements(
    user_achievement_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE NOT NULL,
    achievement_id UUID REFERENCES achievements(achievement_id) ON DELETE CASCADE NOT NULL,
    date_earned TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);