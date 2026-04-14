DROP TYPE IF EXISTS event_type;

CREATE TYPE event_type AS ENUM ('workshop','news', 'hackathon', 'meeting', 'other');

CREATE TABLE IF NOT EXISTS event (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(100) NOT NULL,
	start_at TIMESTAMP NOT NULL,
	end_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	type event_type NOT NULL,
	location VARCHAR(100) NOT NULL,
	description VARCHAR(300) NOT NULL,
	content JSONB NOT NULL,
	max_atendees INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS event_registration (
    event_registration_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID REFERENCES event(event_id) ON DELETE CASCADE NOT NULL,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE NOT null,
    created_at TIMESTAMP NOT NULL
    
);
