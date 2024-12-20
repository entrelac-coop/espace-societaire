ALTER TABLE gifts ADD COLUMN claimed_by_user_id UUID;
--migration:split
ALTER TABLE gifts ADD CONSTRAINT gifts_claimed_by_user_id_foreign_key FOREIGN KEY (claimed_by_user_id) REFERENCES users (id);
