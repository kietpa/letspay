# migration

goose -dir migrations create create_users_table sql
goose -dir migrations postgres "postgres://myuser:mypassword@localhost:5432/letspay" up -v
goose -dir migrations postgres "postgres://myuser:mypassword@localhost:5432/letspay" down -v

# commits

setup DB config, goose migrations

# tables

user
user_login
disbursement_transactions
    id
    transaction_id unique
    user_transaction_id
    bank_code

change agent to provider

