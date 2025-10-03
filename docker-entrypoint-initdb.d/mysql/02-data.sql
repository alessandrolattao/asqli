-- Insert sample data for SQLAI testing

-- Insert users
INSERT INTO users (username, email, first_name, last_name, age, country, is_active) VALUES
('john_doe', 'john.doe@example.com', 'John', 'Doe', 28, 'USA', TRUE),
('jane_smith', 'jane.smith@example.com', 'Jane', 'Smith', 32, 'UK', TRUE),
('bob_wilson', 'bob.wilson@example.com', 'Bob', 'Wilson', 45, 'Canada', TRUE),
('alice_brown', 'alice.brown@example.com', 'Alice', 'Brown', 25, 'Australia', TRUE),
('charlie_davis', 'charlie.davis@example.com', 'Charlie', 'Davis', 38, 'USA', FALSE),
('emma_johnson', 'emma.johnson@example.com', 'Emma', 'Johnson', 29, 'UK', TRUE),
('michael_lee', 'michael.lee@example.com', 'Michael', 'Lee', 41, 'Singapore', TRUE),
('sophia_garcia', 'sophia.garcia@example.com', 'Sophia', 'Garcia', 27, 'Spain', TRUE),
('james_martin', 'james.martin@example.com', 'James', 'Martin', 35, 'France', TRUE),
('olivia_white', 'olivia.white@example.com', 'Olivia', 'White', 31, 'Germany', FALSE);

-- Insert products
INSERT INTO products (name, description, price, stock, category, brand) VALUES
('Laptop Pro 15', 'High-performance laptop with 16GB RAM and 512GB SSD', 1299.99, 50, 'Electronics', 'TechBrand'),
('Wireless Mouse', 'Ergonomic wireless mouse with USB receiver', 29.99, 200, 'Electronics', 'TechBrand'),
('Mechanical Keyboard', 'RGB backlit mechanical keyboard with cherry switches', 89.99, 120, 'Electronics', 'TechBrand'),
('USB-C Hub', '7-in-1 USB-C hub with HDMI and SD card reader', 49.99, 80, 'Electronics', 'AccessoryHub'),
('Monitor 27"', '4K UHD monitor with HDR support', 399.99, 35, 'Electronics', 'DisplayCo'),
('Desk Chair', 'Ergonomic office chair with lumbar support', 249.99, 25, 'Furniture', 'OfficePro'),
('Standing Desk', 'Electric height-adjustable standing desk', 599.99, 15, 'Furniture', 'OfficePro'),
('Webcam HD', '1080p webcam with built-in microphone', 79.99, 60, 'Electronics', 'TechBrand'),
('Headphones', 'Noise-cancelling over-ear headphones', 199.99, 45, 'Electronics', 'AudioMax'),
('Desk Lamp', 'LED desk lamp with adjustable brightness', 39.99, 100, 'Furniture', 'LightCo'),
('Phone Stand', 'Adjustable aluminum phone stand', 19.99, 150, 'Accessories', 'AccessoryHub'),
('Cable Organizer', 'Set of 10 cable clips for desk management', 12.99, 300, 'Accessories', 'AccessoryHub'),
('External SSD 1TB', 'Portable SSD with USB 3.2 Gen 2', 149.99, 70, 'Electronics', 'TechBrand'),
('Laptop Sleeve', 'Padded laptop sleeve for 15" laptops', 24.99, 180, 'Accessories', 'BagCo'),
('Bluetooth Speaker', 'Portable waterproof Bluetooth speaker', 59.99, 90, 'Electronics', 'AudioMax');

-- Insert orders
INSERT INTO orders (user_id, total_amount, status, shipping_address, order_date) VALUES
(1, 1379.98, 'delivered', '123 Main St, New York, NY 10001', '2025-09-15 10:30:00'),
(1, 49.99, 'shipped', '123 Main St, New York, NY 10001', '2025-09-28 14:20:00'),
(2, 689.97, 'processing', '456 Oak Ave, London, UK SW1A 1AA', '2025-09-29 09:15:00'),
(3, 249.99, 'delivered', '789 Maple Dr, Toronto, ON M5H 2N2', '2025-09-10 16:45:00'),
(4, 1299.99, 'cancelled', '321 Beach Rd, Sydney, NSW 2000', '2025-09-20 11:00:00'),
(6, 159.98, 'delivered', '654 High St, Manchester, UK M1 1AD', '2025-09-25 13:30:00'),
(7, 599.99, 'shipped', '987 River Rd, Singapore 018956', '2025-09-30 15:00:00'),
(8, 89.99, 'pending', '147 Plaza St, Madrid, Spain 28001', '2025-10-01 10:00:00'),
(9, 449.98, 'processing', '258 Boulevard, Paris, France 75001', '2025-10-01 12:30:00'),
(1, 199.99, 'pending', '123 Main St, New York, NY 10001', '2025-10-02 09:00:00');

-- Insert order items
INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES
-- Order 1 (user 1)
(1, 1, 1, 1299.99),
(1, 2, 1, 29.99),
(1, 3, 1, 49.99),
-- Order 2 (user 1)
(2, 4, 1, 49.99),
-- Order 3 (user 2)
(3, 5, 1, 399.99),
(3, 2, 1, 29.99),
(3, 8, 1, 79.99),
(3, 11, 9, 19.99),
-- Order 4 (user 3)
(4, 6, 1, 249.99),
-- Order 5 (user 4) - cancelled
(5, 1, 1, 1299.99),
-- Order 6 (user 6)
(6, 9, 1, 199.99),
(6, 12, 1, 12.99),
(6, 14, 2, 24.99),
-- Order 7 (user 7)
(7, 7, 1, 599.99),
-- Order 8 (user 8)
(8, 3, 1, 89.99),
-- Order 9 (user 9)
(9, 5, 1, 399.99),
(9, 4, 1, 49.99),
-- Order 10 (user 1)
(10, 9, 1, 199.99);

-- Insert reviews
INSERT INTO reviews (product_id, user_id, rating, title, comment) VALUES
(1, 1, 5, 'Excellent laptop!', 'Best laptop I have ever owned. Fast, reliable, and great build quality.'),
(1, 2, 4, 'Very good but expensive', 'Great performance but a bit pricey for my budget.'),
(2, 1, 5, 'Perfect mouse', 'Comfortable and responsive. Works flawlessly.'),
(3, 3, 4, 'Good keyboard', 'Nice mechanical feel, but the RGB could be brighter.'),
(5, 2, 5, 'Amazing display', 'Colors are vibrant and the 4K resolution is stunning.'),
(6, 3, 5, 'Best office chair', 'My back pain is gone after using this chair. Highly recommended!'),
(8, 6, 3, 'Decent webcam', 'Image quality is okay but struggles in low light.'),
(9, 1, 5, 'Outstanding sound quality', 'Noise cancellation is incredible. Perfect for work and travel.'),
(9, 6, 5, 'Worth every penny', 'Best headphones in this price range.'),
(10, 7, 4, 'Great lamp', 'Brightness levels are perfect. Wish it had a warmer color option.'),
(13, 2, 5, 'Super fast SSD', 'Transfer speeds are amazing. Very reliable.'),
(15, 8, 4, 'Good speaker', 'Sound quality is excellent for its size. Battery lasts long too.');
