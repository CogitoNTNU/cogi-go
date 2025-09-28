CREATE TABLE IF NOT EXISTS applications (
    application_id serial PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE NOT NULL,

    project_id_1 INT REFERENCES projects(project_id) ON DELETE CASCADE,
    project_id_2 INT REFERENCES projects(project_id) ON DELETE CASCADE,
    project_id_3 INT REFERENCES projects(project_id) ON DELETE CASCADE,

    application_text TEXT NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    isMember BOOLEAN NOT NULL,

    CONSTRAINT check_unique_projects CHECK (
        project_id_1 != project_id_2 AND 
        project_id_1 != project_id_3 AND 
        project_id_2 != project_id_3
    )
);