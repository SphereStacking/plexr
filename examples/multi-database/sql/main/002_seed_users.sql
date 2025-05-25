-- Seed initial users
INSERT INTO users (username, email) VALUES
    ('admin', 'admin@example.com'),
    ('test_user', 'test@example.com')
ON CONFLICT (username) DO NOTHING;

-- Add profiles for seeded users
INSERT INTO user_profiles (user_id, first_name, last_name, bio)
SELECT id, 'Admin', 'User', 'System administrator'
FROM users WHERE username = 'admin'
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO user_profiles (user_id, first_name, last_name, bio)
SELECT id, 'Test', 'User', 'Test account for development'
FROM users WHERE username = 'test_user'
ON CONFLICT (user_id) DO NOTHING;