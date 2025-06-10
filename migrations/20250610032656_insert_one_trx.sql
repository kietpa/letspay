-- +goose Up
insert into disbursement (user_id, reference_id, user_reference_id, status,provider_id,provider_reference_id,amount,created_at,updated_at,bank_account_number,bank_account_name,description)
values (1, 'john222', 'jim222', 'COMPLETED', 1, 'job222', 20000,now(), now(), 123123, 'John Smith', 'investasi');

--

