-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    password varchar(255),
    created_at TIMESTAMP DEFAULT NOW(),
    webhook varchar(255)
);

-- +goose Down
DROP TABLE users;
