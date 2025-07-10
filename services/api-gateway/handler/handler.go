package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"letspay/pkg/helper"
	"letspay/pkg/logger"
	"letspay/services/api-gateway/common/constants"
	"letspay/services/api-gateway/model"
	"letspay/services/api-gateway/mq"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ApiHandler struct {
	mqConn   *amqp.Connection
	validate *validator.Validate
}

func NewApiHandler(
	mqConn *amqp.Connection,
	validate *validator.Validate,
) *ApiHandler {
	return &ApiHandler{
		mqConn:   mqConn,
		validate: validate,
	}
}

func (a *ApiHandler) RequestDisbursement(w http.ResponseWriter, r *http.Request) {
	request := model.CreateDisbursementRequest{}
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Error(context.TODO(), fmt.Sprint("[Create Disbursement] unmarshal error: ", err))
		helper.RespondWithError(w, http.StatusBadRequest, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.INVALID_JSON_BODY,
		})
		return
	}

	if validationErrors := helper.ValidateStruct(request, *a.validate); len(validationErrors) > 0 {
		logger.Error(context.TODO(), fmt.Sprint("[Create Disbursement] validation error: ", validationErrors))
		helper.RespondWithError(w, http.StatusBadRequest, model.Error{
			Code:    http.StatusBadRequest,
			Message: constants.VALIDATION_ERROR,
			Errors:  validationErrors,
		})
		return
	}

	userid, _ := strconv.Atoi(r.Header.Get("X-User-ID"))

	err := mq.PublishDisbursementRequest(a.mqConn,
		model.DisbursementRequestEvent{
			UserId:            userid,
			UserReferenceId:   request.UserReferenceId,
			Amount:            request.Amount,
			BankCode:          request.BankCode,
			BankAccountNumber: request.BankAccountNumber,
			BankAccountName:   request.BankAccountName,
			Description:       request.Description,
		})
	if err != nil {
		logger.Error(context.TODO(), fmt.Sprint("[Create Disbursement] publish to queue error: ", err))
		helper.RespondWithError(w, http.StatusInternalServerError, model.Error{
			Code:    http.StatusInternalServerError,
			Message: constants.INTERNAL_ERROR_MESSAGE,
		})
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Disbursement request successfully received",
	})

}
