-- Admin User is created programmatically by setup_db.rs to ensure correct password hashing
-- (Removed manual insert)

-- Seed Customers
INSERT INTO customers (first_name, last_name, email, phone, address, city, loyalty_points) VALUES
('John', 'Doe', 'john@example.com', '555-0101', '123 Main St', 'New York', 100),
('Jane', 'Smith', 'jane@example.com', '555-0102', '456 Oak Ave', 'Los Angeles', 250),
('Alice', 'Johnson', 'alice@example.com', '555-0103', '789 Pine Ln', 'Chicago', 50)
ON CONFLICT (email) DO NOTHING;

-- Seed Client Profiles
INSERT INTO client_profiles (customer_id, date_of_birth, gender, newsletter_opt_in, account_status)
SELECT id, '1990-01-01', 'Male', true, 'Active' FROM customers WHERE email = 'john@example.com'
ON CONFLICT (customer_id) DO NOTHING;

INSERT INTO client_profiles (customer_id, date_of_birth, gender, newsletter_opt_in, account_status)
SELECT id, '1985-05-15', 'Female', false, 'Active' FROM customers WHERE email = 'jane@example.com'
ON CONFLICT (customer_id) DO NOTHING;

INSERT INTO client_profiles (customer_id, date_of_birth, gender, newsletter_opt_in, account_status)
SELECT id, '1992-11-20', 'Female', true, 'PendingDeletion' FROM customers WHERE email = 'alice@example.com'
ON CONFLICT (customer_id) DO NOTHING;
-- Update Alice to have a deletion request time
UPDATE client_profiles SET deletion_requested_at = NOW() - INTERVAL '2 days' WHERE customer_id = (SELECT id FROM customers WHERE email = 'alice@example.com');


-- Seed Items
INSERT INTO items (name, description, price, stock_quantity, category, sku, is_active) VALUES
('Modern Sofa', 'A comfortable modern grey sofa', 899.99, 10, 'Living Room', 'SOFA-MOD-001', true),
('Oak Dining Table', 'Solid oak dining table for 6', 450.00, 5, 'Dining Room', 'TABLE-OAK-001', true),
('Office Chair', 'Ergonomic office chair', 199.50, 20, 'Office', 'CHAIR-OFF-001', true)
ON CONFLICT (sku) DO NOTHING;

-- Seed Orders
INSERT INTO orders (customer_id, total_amount, status) 
SELECT id, 899.99, 'Delivered' FROM customers WHERE email = 'john@example.com'
LIMIT 1;

INSERT INTO orders (customer_id, total_amount, status) 
SELECT id, 649.50, 'Processing' FROM customers WHERE email = 'jane@example.com'
LIMIT 1;

-- Seed Payments
INSERT INTO payments (order_id, customer_id, amount, payment_method, payment_status, created_at)
SELECT o.id, o.customer_id, o.total_amount, 'Card', 'Completed', NOW() - INTERVAL '5 days'
FROM orders o 
JOIN customers c ON o.customer_id = c.id
WHERE c.email = 'john@example.com'
LIMIT 1;

INSERT INTO payments (order_id, customer_id, amount, payment_method, payment_status, created_at)
SELECT o.id, o.customer_id, o.total_amount, 'PayPal', 'Pending', NOW()
FROM orders o 
JOIN customers c ON o.customer_id = c.id
WHERE c.email = 'jane@example.com'
LIMIT 1;

