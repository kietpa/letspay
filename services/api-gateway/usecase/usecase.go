package usecase

import "letspay/services/api-gateway/model"

type (
	DisbursementUsecase interface {
		HandleDisbursementCompleted(model.DisbursementCompletedEvent)
		HandleDisbursementFailed(model.DisbursementFailedEvent)
	}
)
