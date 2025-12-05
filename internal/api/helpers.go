package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andrey-samosuk/answer-questions/internal/entity"
)

func sendJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
	}
}

func sendError(w http.ResponseWriter, statusCode int, errorMessage string) {
	sendJSON(w, statusCode, entity.ErrorResponse{Error: errorMessage})
}

func sendCustomError(w http.ResponseWriter, err error, logMessage string) {
	if customErr, ok := err.(entity.CustomError); ok {
		log.Printf("%s: %v", logMessage, err)
		sendError(w, customErr.Code, customErr.Message)
		return
	}

	log.Printf("%s: %v", logMessage, err)
	sendError(w, http.StatusInternalServerError, entity.ErrDatabaseQuery.Message)
}
