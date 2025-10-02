CREATE TABLE IF NOT EXISTS project_sponsor (
    project_id UUID REFERENCES projects(project_id) ON DELETE CASCADE NOT NULL,
    sponsor_id UUID REFERENCES sponsors(sponsor_id) ON DELETE CASCADE NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    PRIMARY KEY (project_id, sponsor_id)
);
