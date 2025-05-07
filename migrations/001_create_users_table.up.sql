CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    points INT DEFAULT 0,
    referrer_id INTEGER REFERENCES users(id)
);