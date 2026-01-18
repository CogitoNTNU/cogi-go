CREATE TABLE IF NOT EXISTS temp_applications(
    temp_application_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(100) NOT NULL UNIQUE,
    phone_number INT NOT NULL UNIQUE,
    projects TEXT[] NOT NULL,
    application_text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);