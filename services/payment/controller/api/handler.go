package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"letspay/pkg/helper"
	"letspay/pkg/logger"
	"letspay/pkg/util"
	"letspay/services/payment/common/constants"
	"letspay/services/payment/controller"
	"letspay/services/payment/model"
	"net/http"

	"github.com/gorilla/mux"
)

var notAllowedError = model.Error{
	Code:    http.StatusMethodNotAllowed,
	Message: constants.METHOD_NOT_ALLOWED_MESSAGE,
}

// @Summary		Get Disbursement
// @Description	Get disbursement transaction by reference ID
// @Tags			Disbursement
// @Accept			json
// @Produce		json
// @Param			ID				path		int		true	"reference id"
// @Param			Authorization	header		string	true	"JWT Token"
// @Success		200				{object}	[]model.DisbursementDetail
// @Failure		400				{object}	model.Error
// @Failure		401				{object}	model.Error
// @Failure		500				{object}	model.Error
// @Router			/disbursement/{referenceid} [Get]
func (m *ApiModule) GetDisbursement(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	err := model.Error{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.PROCESS_ID, util.GenerateRandomHex())

	switch r.Method {
	case http.MethodGet:
		vars := mux.Vars(r)
		trxId := vars["referenceId"]
		InputParams := make(map[string]string)
		InputParams["referenceId"] = trxId

		logger.Info(ctx, fmt.Sprintf("[Get Disbursement] request received from %s", helper.GetIP(r)))

		response, err = m.disbursementApi.GetDisbursement(ctx, InputParams)
		if err.Code != 0 {
			controller.RespondWithError(w, err.Code, err)
			return
		}
	default:
		logger.Error(ctx, fmt.Sprintf("[Get Disbursement] method not allowed"))
		controller.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	controller.RespondWithJSON(w, http.StatusOK, response)

}

// @Summary		Create Disbursement
// @Description	Create a disbursement transaction
// @Tags			Disbursement
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string							true	"JWT Token"
// @Param			data			body		model.CreateDisbursementRequest	true	"The input body"
// @Success		200				{object}	model.DisbursementDetail
// @Failure		400				{object}	model.Error
// @Failure		401				{object}	model.Error
// @Failure		500				{object}	model.Error
// @Router			/disbursement [Post]
func (m *ApiModule) CreateDisbursement(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	err := model.Error{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.PROCESS_ID, util.GenerateRandomHex())

	switch r.Method {
	case http.MethodPost:
		param := make(map[string]string)
		body, _ := io.ReadAll(r.Body)
		param[constants.JSON_BODY] = string(body)
		param[constants.USER_ID] = r.Header.Get(constants.X_USER_ID)

		logger.Info(ctx, fmt.Sprintf("[Create Disbursement] request received from %s", helper.GetIP(r)))
		response, err = m.disbursementApi.CreateDisbursement(ctx, param)
		if err.Code != 0 {
			controller.RespondWithError(w, err.Code, err)
			return
		}
	default:
		logger.Error(ctx, fmt.Sprintf("[Create Disbursement] method not allowed"))
		controller.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	controller.RespondWithJSON(w, http.StatusOK, response)
}

func (m *ApiModule) CallbackDisbursement(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	err := model.Error{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.PROCESS_ID, util.GenerateRandomHex())

	switch r.Method {
	case http.MethodPost:
		param := make(map[string]string)
		vars := mux.Vars(r)
		param[constants.PROVIDER] = vars[constants.PROVIDER]

		body, _ := io.ReadAll(r.Body)
		param[constants.JSON_BODY] = string(body)
		encodedHeaders, _ := json.Marshal(r.Header)
		param[constants.REQUEST_HEADERS] = string(encodedHeaders)

		logger.Info(ctx, fmt.Sprintf("[Callback Disbursement] request received from %s", helper.GetIP(r)))

		response, err = m.disbursementApi.CallbackDisbursement(ctx, param)
		if err.Code != 0 {
			controller.RespondWithError(w, err.Code, err)
			return
		}
	default:
		logger.Error(ctx, fmt.Sprintf("[Callback Disbursement] method not allowed"))
		controller.RespondWithError(w, http.StatusMethodNotAllowed, notAllowedError)
		return
	}

	controller.RespondWithJSON(w, http.StatusOK, response)
}
