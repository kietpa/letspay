package rabbitmq

import (
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func Connect(url string) *amqp091.Connection {
	var conn *amqp091.Connection
	var err error

	// add some delay because amqp takes some time to start
	time.Sleep(10 * time.Second)
	for i := 0; i < 5; i++ {
		conn, err = amqp091.Dial(url)
		if err == nil {
			return conn
		}
		log.Printf("🐇 Failed to connect to RabbitMQ: %v. Retrying...", err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("❌ Failed to connect to RabbitMQ after retries: %v", err)
	return nil
}
