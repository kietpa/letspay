# inside services/user-service/Dockerfile
FROM golang:1.23-alpine

WORKDIR /app
COPY ../../ .
WORKDIR /app/services/api-gateway
RUN go build -o api-gateway .
EXPOSE 8080
CMD ["./api-gateway"]
