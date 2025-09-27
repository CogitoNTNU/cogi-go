CREATE TABLE IF NOT EXISTS project_sponsor (
    project_id INT REFERENCES projects(project_id) ON DELETE CASCADE,
    sponsor_id INT REFERENCES sponsors(sponsor_id) ON DELETE CASCADE,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    PRIMARY KEY (project_id, sponsor_id)
);
