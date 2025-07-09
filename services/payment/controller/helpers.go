package controller

import (
	"encoding/json"
	"letspay/services/payment/model"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error encoding JSON: %v", err)
		}
	}
}

func RespondWithError(w http.ResponseWriter, statusCode int, errors model.Error) {
	RespondWithJSON(w, statusCode, errors)
}
