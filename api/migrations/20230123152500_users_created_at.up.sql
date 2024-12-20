ALTER TABLE users ADD COLUMN created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

--bun:split

UPDATE users SET created_at = null;
