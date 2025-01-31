CREATE EXTENSION IF NOT EXISTS "uuid-ossp";  -- Ensure the uuid-ossp extension is available

CREATE TABLE IF NOT EXISTS orders (
    order_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    order_status VARCHAR(50) CHECK (order_status IN ('Pending', 'Completed', 'Cancelled')),
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
