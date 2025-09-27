CREATE TABLE IF NOT EXISTS roles(
    role_id serial PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL
)