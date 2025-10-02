DROP TYPE IF EXISTS gender_type;
CREATE TYPE gender_type AS ENUM ('male', 'female', 'non_binary', 'other');

CREATE TABLE IF NOT EXISTS users(
   user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   first_name VARCHAR (50) NOT NULL,
   last_name VARCHAR (50) NOT NULL,
   description TEXT NOT NULL,
   nickname VARCHAR (300) UNIQUE NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL,
   phone VARCHAR(15) UNIQUE NOT NULL,
   gender gender_type NOT NULL,
   github_url VARCHAR (100),
   linkedin_url VARCHAR (100),
   kaggle_url VARCHAR (100),
   huggingface_url VARCHAR (100),
   password VARCHAR (300) NOT NULL,
   avatar VARCHAR (300),
   image_permission TIMESTAMP,
   food_preference TEXT[]
);