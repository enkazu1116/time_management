package response

import (
	"encoding/json"
	"log"
	"net/http"

	"time_management/infrastructure/util/messages"
)

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set(messages.ContentTypeHeader, messages.ContentTypeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	message := messages.Error(status, err.Error())
	WriteJSON(w, message.StatusCode, map[string]string{message.Key: message.Message})
}
