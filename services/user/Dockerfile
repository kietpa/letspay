# inside services/user-service/Dockerfile
FROM golang:1.23-alpine

WORKDIR /app
COPY ../../ .
WORKDIR /app/services/user
RUN go build -o user-service .
EXPOSE 8080
CMD ["./user-service"]
