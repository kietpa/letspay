package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"letspay/services/api-gateway/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishDisbursementRequest(conn *amqp.Connection, payload model.DisbursementRequestEvent) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel error: %w", err)
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
		return fmt.Errorf("exchange declare error: %w", err)
	}

	body, _ := json.Marshal(payload)
	return ch.PublishWithContext(context.TODO(),
		"disbursement.events",
		"disbursement.requested",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
