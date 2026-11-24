package http

import (
	"net/http"

	"time_management/backend/internal/infrastrcture/http/handlers"

	"github.com/gorilla/mux"
)

/*
 * NewRouter
 * 各 Domain Handlerのルーティングを実施し、Registerに登録
 */
func NewRouter(userHandler *handlers.UserHandler) http.Handler {
	router := mux.NewRouter()

	userHandler.Register(router)

	return router
}
