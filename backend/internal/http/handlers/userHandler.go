package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"time_management/backend/domain/user/dto"
	userusecases "time_management/backend/domain/user/usecases"
	"time_management/backend/internal/http/responses"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	usecase *userusecases.UserUsecase
}

func NewUserHandler(usecase *userusecases.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) Register(r *mux.Router) {
	r.HandleFunc("/users", h.create).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", h.getByID).Methods(http.MethodGet)
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.usecase.CreateUser(r.Context(), req)
	if err != nil {
		responses.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.WriteJSON(w, http.StatusCreated, res)
}

func (h *UserHandler) getByID(w http.ResponseWriter, r *http.Request) {
	res, err := h.usecase.GetUserByID(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			responses.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		responses.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	responses.WriteJSON(w, http.StatusOK, res)
}
