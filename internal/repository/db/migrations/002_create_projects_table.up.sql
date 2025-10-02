CREATE TABLE IF NOT EXISTS projects(
   project_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   title VARCHAR (300) NOT NULL,
   github_url VARCHAR (300) NOT NULL,
   logo VARCHAR (300) NOT NULL,
   playable BOOLEAN NOT NULL,
   released BOOLEAN NOT NULL,
   active_project BOOLEAN NOT NULL,
   project_url VARCHAR (300)
);
