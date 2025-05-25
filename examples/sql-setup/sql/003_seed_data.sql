-- Insert sample users
INSERT INTO users (username, email) VALUES 
    ('admin', 'admin@example.com'),
    ('john_doe', 'john@example.com'),
    ('jane_smith', 'jane@example.com')
ON CONFLICT (username) DO NOTHING;

-- Insert sample products
INSERT INTO products (name, description, price, stock_quantity) VALUES 
    ('Laptop', 'High-performance laptop', 999.99, 10),
    ('Mouse', 'Wireless mouse', 29.99, 50),
    ('Keyboard', 'Mechanical keyboard', 79.99, 30),
    ('Monitor', '27-inch 4K monitor', 399.99, 15)
ON CONFLICT DO NOTHING;

-- Insert sample orders
INSERT INTO orders (user_id, total_amount, status) 
SELECT 
    u.id,
    999.99,
    'completed'
FROM users u 
WHERE u.username = 'john_doe'
AND NOT EXISTS (
    SELECT 1 FROM orders o WHERE o.user_id = u.id
);