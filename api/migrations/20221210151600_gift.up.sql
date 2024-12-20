CREATE TABLE gifts (
  id TEXT NOT NULL DEFAULT gen_random_uuid(),
  code TEXT NOT NULL,

  CONSTRAINT gifts_primary_key PRIMARY KEY (id),
  CONSTRAINT gifts_code_unique UNIQUE (code)
);
--migration:split
ALTER TABLE payments ADD COLUMN gift_id TEXT;
--migration:split
ALTER TABLE payments ADD CONSTRAINT payments_gift_id_foreign_key FOREIGN KEY (gift_id) REFERENCES gifts (id);
--migration:split
ALTER TABLE payments ADD CONSTRAINT payments_gift_id_unique UNIQUE (gift_id);
