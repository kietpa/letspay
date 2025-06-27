Simple payment aggregator app that's still being worked on. Currently a monolith, will be divided into a microservice later for learning purposes.

Stack: Go, PostgreSQL, gorilla/mux, pgx

# Features
- Disbursements using Xendit
- Simple user system
- Scheduler to update pending disbursements

# How to Start
1. add credentials in the env file
2. docker compose up -d
3. run migrations:
    goose -dir migrations postgres "user=letsuser password=letspassword dbname=letspay port=5372 sslmode=disable" up

### To test webhooks:
1. deploy app and expose it (eg: port forwarding) 
2. fill the webhook url in the provider dashboards

# API Documentation
(WIP)

# To do
- [x] struct validation
- [x] user system with JWT auth
- [x] Basic disbursement system (create, get status, webhook)
- [x] cron to check & update disbursement status
- [x] add basic openapi documentation
- [ ] improved logging with grafana/elk
- [ ] handle race conditions (idempotency, redis, callback retries)
- [ ] add other provider (midtrans) and implement provider switching
- [ ] improve readme (table of contents, etc)
- [ ] improved error handling (handle provider response, database errors, catch panics)
- [ ] handle context timeouts
- [ ] split into microservices?
