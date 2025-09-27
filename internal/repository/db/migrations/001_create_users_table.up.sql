DROP TYPE IF EXISTS gender_type;
CREATE TYPE gender_type AS ENUM ('male', 'female', 'other');

CREATE TABLE IF NOT EXISTS users(
   user_id serial PRIMARY KEY,
   first_name VARCHAR (50) NOT NULL,
   last_name VARCHAR (50) NOT NULL,
   description TEXT,
   nickname VARCHAR (300) UNIQUE NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL,
   tlf_number VARCHAR(15),
   gender gender_type NOT NULL,
   github_url VARCHAR (100),
   linkedin_url VARCHAR (100),
   kaggle_url VARCHAR (100),
   huggingface_url VARCHAR (100),
   password_hash VARCHAR (300) NOT NULL,
   avatar VARCHAR (300),
   image_permission TIMESTAMP,
   allergy_info TEXT[]
);