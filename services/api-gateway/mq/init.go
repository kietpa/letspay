package mq

import amqp "github.com/rabbitmq/amqp091-go"

func InitConsumers(conn *amqp.Connection) {
	StartDisbursementCompletedConsumer(conn)
}
