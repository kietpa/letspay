-- +goose Up
insert into disbursements (user_id, reference_id, user_reference_id, status,provider_id,provider_reference_id,amount,bank_code,created_at,updated_at,bank_account_number,bank_account_name,description,failure_code)
values (1, 'john222', 'jim222', 'COMPLETED', 1, 'job222', 20000,'MANDIRI',now(), now(), 123123, 'John Smith', 'investasi', '');

-- +goose Down
DELETE FROM disbursements
WHERE reference_id = 'john222';