-- Enforce user ownership for all orders
DELETE FROM order_items
WHERE order_id IN (SELECT id FROM orders WHERE user_id IS NULL);

DELETE FROM orders
WHERE user_id IS NULL;

ALTER TABLE orders
ALTER COLUMN user_id SET NOT NULL;