package mq

import (
	"encoding/json"
	"letspay/services/payment/model"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeDisbursementRequest(conn *amqp.Connection, callback func(model.DisbursementRequestEvent)) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("channel error:", err)
	}

	err = ch.ExchangeDeclare(
		"disbursement.events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("exchange error:", err)
	}

	queue, err := ch.QueueDeclare(
		"disbursement.requested.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("queue declare error:", err)
	}

	err = ch.QueueBind(
		queue.Name,
		"disbursement.requested", // key
		"disbursement.events",    // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatal("queue bind error:", err)
	}

	messages, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("consume error:", err)
	}

	log.Println("📥 Waiting for disbursement.requested events...")
	for msg := range messages {
		var event model.DisbursementRequestEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("bad message: %s", msg.Body)
			continue
		}
		callback(event)
	}
}
