CREATE EXTENSION IF NOT EXISTS "uuid-ossp";  -- Ensure the uuid-ossp extension is available

CREATE TABLE IF NOT EXISTS meals (
    meal_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
