-- Add order status support for dashboard filtering and badges
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS status VARCHAR(20);

UPDATE orders
SET status = 'pending'
WHERE status IS NULL OR status = '';

ALTER TABLE orders
    ALTER COLUMN status SET DEFAULT 'pending';

ALTER TABLE orders
    ALTER COLUMN status SET NOT NULL;