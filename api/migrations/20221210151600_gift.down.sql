ALTER TABLE payments DROP CONSTRAINT payments_gift_id_unique;
--migration:split
ALTER TABLE payments DROP CONSTRAINT payments_gift_id_foreign_key;
--migration:split
ALTER TABLE payments DROP COLUMN gift_id;
--migration:split
DROP TABLE gifts;
