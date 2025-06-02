-- +goose Up
CREATE TABLE disbursement (
    id SERIAL PRIMARY KEY,
    user_id int,
    reference_id varchar(255),
    user_reference_id varchar(255),
    status varchar(30),
    provider_id int,
    provider_reference_id varchar(255),
    amount float,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    bank_account_number varchar(50),
    bank_account_name varchar(255),
    description TEXT
);

-- +goose Down
DROP TABLE disbursement;
