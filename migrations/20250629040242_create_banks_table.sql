-- +goose Up
CREATE TABLE banks (
    id SERIAL PRIMARY KEY,
    bank_name varchar(255),
    bank_code varchar(100),
    providers varchar(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE banks;