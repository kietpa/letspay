package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"letspay/services/payment/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishDisbursementCompleted(conn *amqp.Connection, payload model.DisbursementCompletedEvent) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel error: %w", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("disbursement.events", "topic", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("exchange declare error: %w", err)
	}

	body, _ := json.Marshal(payload)
	return ch.PublishWithContext(context.TODO(),
		"disbursement.events", "disbursement.completed", false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func PublishDisbursementFailed(conn *amqp.Connection, payload model.DisbursementFailedEvent) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel error: %w", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("disbursement.events", "topic", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("exchange declare error: %w", err)
	}

	body, _ := json.Marshal(payload)
	return ch.PublishWithContext(context.TODO(),
		"disbursement.events", "disbursement.failed", false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
