DROP TYPE IF EXISTS semester;
DROP TYPE IF EXISTS group_type;

CREATE TYPE semester AS ENUM ('spring', 'autumn');
CREATE TYPE group_type AS ENUM ('marketing', 'board', 'social');

CREATE TABLE IF NOT EXISTS teams (
    team_id serial PRIMARY KEY,
    semester semester NOT NULL,
    year INT NOT NULL
);

CREATE TABLE IF NOT EXISTS project_teams (
    project_id INT REFERENCES projects(project_id) ON DELETE CASCADE
) INHERITS (teams);

CREATE TABLE IF NOT EXISTS administrasjon_teams (
    group_type group_type NOT NULL
) INHERITS (teams);