package handlers

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"time_management/domain/time_tracking/attestation"
	"time_management/infrastructure/http/requests"
	"time_management/infrastructure/http/responses"
	"time_management/infrastructure/util/parser"
	"time_management/infrastructure/util/response"
)

const AttestAdminTokenEnv = "ATTEST_ADMIN_TOKEN"

var (
	ErrAttestationDisabled = errors.New("attestation endpoint is disabled")
	ErrUnauthorized        = errors.New("invalid bearer token")
)

type attestationUsecase interface {
	Attest(context.Context, string, time.Time) (attestation.Result, error)
	Verify(context.Context, string, time.Time) (attestation.Result, error)
}

type AttestationHandler struct {
	usecase    attestationUsecase
	adminToken string
}

func NewAttestationHandler(usecase attestationUsecase) *AttestationHandler {
	return &AttestationHandler{usecase: usecase, adminToken: os.Getenv(AttestAdminTokenEnv)}
}

func (handler *AttestationHandler) Attest(w http.ResponseWriter, r *http.Request) {
	if handler.adminToken == "" {
		response.WriteError(w, http.StatusServiceUnavailable, ErrAttestationDisabled)
		return
	}
	if !constantTimeBearerMatch(r.Header.Get("Authorization"), handler.adminToken) {
		response.WriteError(w, http.StatusUnauthorized, ErrUnauthorized)
		return
	}
	var request requests.DailyAttestationRequest
	if !decode(w, r, &request) || !valid(w, request.Validate()) {
		return
	}
	date, err := parser.ParseDate(request.Date)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	result, err := handler.usecase.Attest(r.Context(), request.UserID, date)
	if err != nil {
		response.WriteError(w, http.StatusBadGateway, err)
		return
	}
	response.WriteJSON(w, http.StatusCreated, responses.Attestation(result))
}

func (handler *AttestationHandler) Verify(w http.ResponseWriter, r *http.Request) {
	request := requests.DailyAttestationRequest{
		UserID: r.URL.Query().Get("user_id"),
		Date:   r.URL.Query().Get("date"),
	}
	if !valid(w, request.Validate()) {
		return
	}
	date, err := parser.ParseDate(request.Date)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	result, err := handler.usecase.Verify(r.Context(), request.UserID, date)
	if err != nil {
		response.WriteError(w, http.StatusBadGateway, err)
		return
	}
	response.WriteJSON(w, http.StatusOK, responses.Attestation(result))
}

// This static token is intentionally only a demo-level write guard, not a
// general authentication system. Hashing both values before ConstantTimeCompare
// keeps the comparison fixed-length even when the supplied token length differs.
func constantTimeBearerMatch(header, expected string) bool {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return false
	}
	suppliedHash := sha256.Sum256([]byte(strings.TrimPrefix(header, prefix)))
	expectedHash := sha256.Sum256([]byte(expected))
	return subtle.ConstantTimeCompare(suppliedHash[:], expectedHash[:]) == 1
}
