CREATE TABLE IF NOT EXISTS applications (
    application_id serial PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    project_id INT REFERENCES projects(project_id) ON DELETE CASCADE,
    application_text TEXT NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);