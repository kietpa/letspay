package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func Connect(url string) *amqp091.Connection {
	conn, err := amqp091.Dial(url)
	failOnError(err, "fail to connect to amqp") // "amqp://guest:guest@localhost:5672/"

	return conn
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
