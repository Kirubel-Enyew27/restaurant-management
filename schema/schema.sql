CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE meals (
    meal_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE drinks (
    drink_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    total_price NUMERIC(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL, 
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_meals (
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    meal_id INT REFERENCES meals(meal_id),
    quantity INT NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (order_id, meal_id)
);

CREATE TABLE order_drinks (
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    drink_id INT REFERENCES drinks(drink_id),
    quantity INT NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (order_id, drink_id)
);


