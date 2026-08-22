package handlers

import (
	"net/http"

	"time_management/infrastructure/util/messages"
	"time_management/infrastructure/util/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	message := messages.OK()
	response.WriteJSON(w, message.StatusCode, map[string]string{message.Key: message.Message})
}
