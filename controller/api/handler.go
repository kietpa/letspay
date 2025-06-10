package api

import (
	"context"
	"letspay/controller"
	"net/http"

	"github.com/gorilla/mux"
)

func (m *ApiModule) GetDisbursement(w http.ResponseWriter, r *http.Request) {
	response := controller.Data{}
	var err error
	// TODO: handle context
	ctx := context.Background()

	switch r.Method {
	case http.MethodGet:
		vars := mux.Vars(r)
		trxId := vars["referenceId"]
		InputParams := make(map[string]string)
		InputParams["referenceId"] = trxId
		response, err = m.disbursementAPI.GetDisbursement(ctx, InputParams)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Error")
		}
	default:
		// TODO: create error message/type constants
		respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	respondWithJSON(w, http.StatusOK, response)

}
