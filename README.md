Simple payment aggregator app that's still being worked on. Currently a monolith, will be divided into a microservice later for learning purposes.

# Features
- Disbursements using Xendit
- Stack: Go, PostgreSQL, gorilla/mux, pgx

# how to start
1. docker compose up -d
2. run migrations:
    goose -dir migrations postgres "user=letsuser password=letspassword dbname=letspay port=5372 sslmode=disable" up
change port if needed
3. add credentials in the env file
4. run the app

# To do
- [x] struct validation
- [ ] user system with JWT auth
- [ ] Basic disbursement system
- [ ] improved logging with grafana/elk
- [ ] improved error handling (handle provider response, database errors, catch panics)
- [ ] handle timeouts
- [ ] handle race conditions
- [ ] scheduler system to check & update disbursement status
- [ ] add api documentation & improve readme (table of contents, etc)