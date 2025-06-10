package api

import (
	"context"
	"letspay/common/constants"
	"letspay/controller"
	"letspay/model"
	"net/http"

	"github.com/gorilla/mux"
)

func (m *ApiModule) GetDisbursement(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	err := model.Error{}
	// TODO: handle context
	ctx := context.Background()

	switch r.Method {
	case http.MethodGet:
		vars := mux.Vars(r)
		trxId := vars["referenceId"]
		InputParams := make(map[string]string)
		InputParams["referenceId"] = trxId
		response, err = m.disbursementAPI.GetDisbursement(ctx, InputParams)
		if err.Code != 0 {
			respondWithError(w, err.Code, err.Message)
			return
		}
	default:
		respondWithError(w, http.StatusMethodNotAllowed, constants.METHOD_NOT_ALLOWED_MESSAGE)
		return
	}

	respondWithJSON(w, http.StatusOK, response)

}
