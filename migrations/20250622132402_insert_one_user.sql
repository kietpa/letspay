-- +goose Up
insert into users (name, email, password)
values ('john', 'john@gmail.com', 'this is a sample user')

-- +goose Down
DELETE FROM users
WHERE email = 'john@gmail.com';
