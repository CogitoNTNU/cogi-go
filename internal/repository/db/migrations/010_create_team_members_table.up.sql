CREATE TABLE IF NOT EXISTS team_members (
    team_id UUID REFERENCES teams(team_id) ON DELETE CASCADE NOT NULL,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE NOT NULL,
    role_id UUID REFERENCES roles(role_id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (team_id, user_id, role_id)
)