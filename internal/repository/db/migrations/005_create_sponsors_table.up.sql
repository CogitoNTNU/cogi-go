DROP TYPE IF EXISTS sponsor_level;
CREATE TYPE sponsor_level AS ENUM ('gold', 'silver', 'bronze');

CREATE TABLE IF NOT EXISTS sponsors (
    sponsor_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR (100) NOT NULL,
    logo VARCHAR (300) NOT NULL,
    website VARCHAR (300) NOT NULL,
    description TEXT,
    level sponsor_level NOT NULL
);