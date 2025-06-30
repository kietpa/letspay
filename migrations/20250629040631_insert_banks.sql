-- +goose Up
insert into banks (bank_name, bank_code, providers)
values ('Bank Central Asia', 'BCA', '2,1'),
('Bank Mandiri', 'MANDIRI', '2,1'),
('Bank Rakyat Indonesia', 'BRI', '2,1'),
('Bank Nasional Indonesia', 'BNI', '2,1');

-- +goose Down
TRUNCATE TABLE banks;