Simple payment aggregator app that's still being worked on. The app uses microservices for learning purposes.

**Stack:** Go, PostgreSQL, gorilla/mux, pgx

**API Docs:** [link](https://kietpa.github.io/projects/letspay/)

# Features
- Monorepo microservice architecture
- Disbursements using Xendit & Midtrans
- Grafana logging using Loki & Promtail
- Redis for idempotency & caching
- Simple user system w/JWT
- Scheduler to update pending disbursements

# How to Run
1. add provider credentials in the .sample.env file in the payments service
2. rename the .sample.env files to .env (in root and each service)
3. run docker on your machine
4. make run
5. make migrate-up
6. send requests at http://localhost:8080

### To test webhooks:
1. run app and expose http://localhost:8080 with ngrok, cloudflare tunneling, etc
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
- [x] add other provider (midtrans) and implement provider switching
- [x] split into microservices
- [ ] implement rabbitmq/kafka
- [ ] handle race conditions using go routines & channels
- [ ] add rate limiter & test it
- [ ] improve caching
- [ ] improved error handling (handle provider response, database errors, catch panics)
- [ ] handle context timeouts
- [ ] add user balance management
