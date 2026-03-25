-- Ensure users.email exists and remains optional
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS email VARCHAR(255);

-- Ensure users.role exists and defaults to user
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role VARCHAR(20);

ALTER TABLE users
    ALTER COLUMN role SET DEFAULT 'user';

UPDATE users SET role = 'user' WHERE role IS NULL OR role = '';

-- Ensure at least one admin user exists
-- password: Test1234
INSERT INTO users (username, email, password, role)
VALUES (
    'admin',
    'admin@inventory.com',
    '$2a$10$BvZM.DEP8RSRg/yurAqlVuPlXpamUZLDWO0TQSsJPPt0lMyUtX6OW',
    'admin'
)
ON CONFLICT (username) DO NOTHING;
