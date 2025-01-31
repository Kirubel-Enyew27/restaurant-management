CREATE EXTENSION IF NOT EXISTS "uuid-ossp";  -- Ensure the uuid-ossp extension is available

CREATE TABLE IF NOT EXISTS order_items (
    order_item_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID REFERENCES orders(order_id) ON DELETE CASCADE,
    meal_id UUID REFERENCES meals(meal_id) ON DELETE CASCADE,
    quantity INT CHECK (quantity > 0),
    price DECIMAL(10,2) NOT NULL
);
