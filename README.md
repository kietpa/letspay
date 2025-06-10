# Tech stack
- Golang
- Postgresql
- gorilla/mux

# tables
user
user_login
disbursement

# how to start
1. docker compose up -d
2. run migrations:
    goose -dir migrations postgres "user=letsuser password=letspassword dbname=letspay port=5372 sslmode=disable" up
change port if needed
3. add credentials in the env file (WIP)
4. run the app (WIP)