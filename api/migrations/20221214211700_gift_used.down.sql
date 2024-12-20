ALTER TABLE gifts DROP CONSTRAINT gifts_claimed_by_user_id_foreign_key;
--migration:split
ALTER TABLE gifts DROP COLUMN claimed_by_user_id;
