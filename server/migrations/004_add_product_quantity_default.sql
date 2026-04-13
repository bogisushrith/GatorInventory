-- Ensure products.quantity exists with a safe default for both new and existing rows
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS quantity INTEGER;

UPDATE products
SET quantity = 0
WHERE quantity IS NULL;

ALTER TABLE products
    ALTER COLUMN quantity SET DEFAULT 0;

ALTER TABLE products
    ALTER COLUMN quantity SET NOT NULL;
