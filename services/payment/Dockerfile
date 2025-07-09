# inside services/user-service/Dockerfile
FROM golang:1.23-alpine

WORKDIR /app
COPY ../../ .
WORKDIR /app/services/payment
RUN go build -o payment-service .
EXPOSE 8080
CMD ["./payment-service"]
