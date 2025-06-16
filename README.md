Simple payment aggregator app that's still being worked on. Currently a monolith, will be divided into a microservice later (for learning purposes).

# Features
- Disbursements using Xendit & Midtrans

# Tech stack
- Golang
- Postgresql
- gorilla/mux

# how to start
1. docker compose up -d
2. run migrations:
    goose -dir migrations postgres "user=letsuser password=letspassword dbname=letspay port=5372 sslmode=disable" up
change port if needed
3. add credentials in the env file (WIP)
4. run the app (WIP)

# To do
- add api documentation & improve readme (table of contents, etc)
- 

# tables
user
user_login
disbursement