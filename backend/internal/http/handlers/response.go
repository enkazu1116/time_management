package handlers

import (
	"encoding/json"
	"net/http"
)

/*
 * WriteJSON
 * Handler 共通の JSON レスポンス書き込み。
 */
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
