Simple payment aggregator app that's still being worked on. Currently a monolith, will be divided into a microservice later for learning purposes.

**Stack:** Go, PostgreSQL, gorilla/mux, pgx

**API Docs:** [link](https://kietpa.github.io/projects/letspay/)

# Features
- Disbursements using Xendit & Midtrans
- Grafana logging using Loki & Promtail
- Redis for idempotency & caching
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
2. fill the webhook url in the provider dashboards (eg: https://asdf.com/callback/xendit)

note: midtrans API key for disbursements (IRIS) needs business registration

# To do
- [x] struct validation
- [x] user system with JWT auth
- [x] Basic disbursement system (create, get status, webhook)
- [x] cron to check & update disbursement status
- [x] add basic openapi documentation
- [x] improved logging with grafana/elk
- [x] use redis to handle idempotency on xendit's callback retries
- [ ] add other provider (midtrans) and implement provider switching
- [ ] handle race conditions using go routines & channels
- [ ] implement microservices ()
- [ ] improve readme (table of contents, etc)
- [ ] improved error handling (handle provider response, database errors, catch panics)
- [ ] handle context timeouts
