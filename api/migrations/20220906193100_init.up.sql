CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--bun:split

CREATE TABLE users (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	confirmed BOOL NOT NULL DEFAULT false,
	confirm_token TEXT,
	email TEXT UNIQUE NOT NULL, 
	password TEXT NOT NULL, 
	phone_number TEXT NOT NULL,
	first_name TEXT NOT NULL, 
	last_name TEXT NOT NULL, 
	address TEXT NOT NULL, 
	postal_code TEXT NOT NULL, 
	city TEXT NOT NULL, 
	country TEXT NOT NULL, 
	category TEXT NOT NULL, 
	customer TEXT NOT NULL, 
	identity_front TEXT, 
	identity_back TEXT, 
	address_proof TEXT
);

--bun:split

CREATE TABLE payments (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	stripe_event_id TEXT UNIQUE NOT NULL,
	shares INTEGER NOT NULL,
	user_id uuid NOT NULL REFERENCES users(id),
	created_at TIMESTAMPTZ NOT NULL
);
