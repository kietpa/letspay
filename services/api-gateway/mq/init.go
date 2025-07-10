package mq

import (
	"letspay/services/api-gateway/usecase"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitConsumers(
	conn *amqp.Connection,
	disbursementUsecase usecase.DisbursementUsecase,
) {
	go StartDisbursementCompletedConsumer(conn, disbursementUsecase.HandleDisbursementCompleted)
	go StartDisbursementFailedConsumer(conn, disbursementUsecase.HandleDisbursementFailed)
}
