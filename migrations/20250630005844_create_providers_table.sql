-- +goose Up
CREATE TABLE providers (
    id SERIAL PRIMARY KEY,
    provider_name varchar(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE providers;