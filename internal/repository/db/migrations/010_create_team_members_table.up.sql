CREATE TABLE IF NOT EXISTS team_members (
    team_id INT REFERENCES teams(team_id) ON DELETE CASCADE NOT NULL,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE NOT NULL,
    role_id INT REFERENCES roles(role_id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (team_id, user_id, role_id)
)