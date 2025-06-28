Simple payment aggregator app that's still being worked on. Currently a monolith, will be divided into a microservice later for learning purposes.

**Stack:** Go, PostgreSQL, gorilla/mux, pgx

**API Docs:** [link](https://kietpa.github.io/projects/letspay/)

# Features
- Disbursements using Xendit
- Grafana logging using Loki & Promtail
- Simple user system w/JWT
- Scheduler to update pending disbursements

# How to Run
1. docker compose up -d
2. run migrations:
    goose -dir migrations postgres "user=letsuser password=letspassword dbname=letspay port=5372 sslmode=disable" up
3. add credentials in the env file
4. send requests at http://localhost:8080

### To test webhooks:
1. docker-compose app and expose http://localhost:8080 with ngrok, cloudflare tunneling, etc
2. fill the webhook url in the provider dashboards (eg https://asdf.com/callback/xendit)

# To do
- [x] struct validation
- [x] user system with JWT auth
- [x] Basic disbursement system (create, get status, webhook)
- [x] cron to check & update disbursement status
- [x] add basic openapi documentation
- [x] improved logging with grafana/elk
- [ ] handle race conditions (idempotency, redis, callback retries)
- [ ] add other provider (midtrans) and implement provider switching
- [ ] implement microservices ()
- [ ] improve readme (table of contents, etc)
- [ ] improved error handling (handle provider response, database errors, catch panics)
- [ ] handle context timeouts
