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
5. deploy app on ngrok (so providers can communicate)

# To do
- [x] struct validation
- [x] user system with JWT auth
- [x] Basic disbursement system (create, get status, webhook)
- [x] cron to check & update disbursement status
- [ ] handle race conditions (idempotency, redis, callback retries)
- [ ] add other provider (midtrans) and implement switching
- [ ] improved logging with grafana/elk
- [ ] add api documentation & improve readme (table of contents, etc)
- [ ] improved error handling (handle provider response, database errors, catch panics)
- [ ] handle context timeouts
- [ ] split into microservices?
