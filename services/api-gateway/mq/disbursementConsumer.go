package mq

import (
	"encoding/json"
	"letspay/services/api-gateway/model"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartDisbursementRequestConsumer(conn *amqp.Connection, callback func(model.DisbursementRequestEvent)) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel error: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("disbursement.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange declare error: %v", err)
	}

	q, err := ch.QueueDeclare("disbursement.requested.queue", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue declare error: %v", err)
	}

	err = ch.QueueBind(q.Name, "disbursement.requested", "disbursement.events", false, nil)
	if err != nil {
		log.Fatalf("queue bind error: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %v", err)
	}

	log.Println("📥 Waiting for disbursement.requested events...")
	for msg := range msgs {
		var event model.DisbursementRequestEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("bad message: %s", msg.Body)
			continue
		}
		callback(event)
	}
}

func StartDisbursementCompletedConsumer(conn *amqp.Connection, callback func(model.DisbursementCompletedEvent)) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel error: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("disbursement.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange declare error: %v", err)
	}

	q, err := ch.QueueDeclare("disbursement.completed.queue", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue declare error: %v", err)
	}

	err = ch.QueueBind(q.Name, "disbursement.completed", "disbursement.events", false, nil)
	if err != nil {
		log.Fatalf("queue bind error: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %v", err)
	}

	log.Println("📥 Waiting for disbursement.completed events...")
	for msg := range msgs {
		var event model.DisbursementCompletedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("bad message: %s", msg.Body)
			continue
		}
		callback(event)
	}
}

func StartDisbursementFailedConsumer(conn *amqp.Connection, callback func(model.DisbursementFailedEvent)) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel error: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("disbursement.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange declare error: %v", err)
	}

	q, err := ch.QueueDeclare("disbursement.failed.queue", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue declare error: %v", err)
	}

	err = ch.QueueBind(q.Name, "disbursement.failed", "disbursement.events", false, nil)
	if err != nil {
		log.Fatalf("queue bind error: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %v", err)
	}

	log.Println("📥 Waiting for disbursement.failed events...")
	for msg := range msgs {
		var event model.DisbursementFailedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("bad message: %s", msg.Body)
			continue
		}
		callback(event)
	}
}

func StartDisbursementGetConsumer(conn *amqp.Connection, callback func(model.DisbursementGetEvent)) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel error: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("disbursement.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange declare error: %v", err)
	}

	q, err := ch.QueueDeclare("disbursement.get.queue", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue declare error: %v", err)
	}

	err = ch.QueueBind(q.Name, "disbursement.get", "disbursement.events", false, nil)
	if err != nil {
		log.Fatalf("queue bind error: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %v", err)
	}

	log.Println("📥 Waiting for disbursement.get events...")
	for msg := range msgs {
		var event model.DisbursementGetEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("bad message: %s", msg.Body)
			continue
		}
		callback(event)
	}
}
