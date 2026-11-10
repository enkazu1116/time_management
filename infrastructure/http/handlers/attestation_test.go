package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"time_management/domain/time_tracking/attestation"
)

type stubAttestationUsecase struct {
	attestResult attestation.Result
	verifyResult attestation.Result
	attestCalls  int
}

func (stub *stubAttestationUsecase) Attest(context.Context, string, time.Time) (attestation.Result, error) {
	stub.attestCalls++
	return stub.attestResult, nil
}

func (stub *stubAttestationUsecase) Verify(context.Context, string, time.Time) (attestation.Result, error) {
	return stub.verifyResult, nil
}

func TestAttestIsDisabledWhenAdminTokenIsUnset(t *testing.T) {
	t.Setenv(AttestAdminTokenEnv, "")
	usecase := &stubAttestationUsecase{}
	handler := NewAttestationHandler(usecase)
	request := httptest.NewRequest(http.MethodPost, "/attestations", strings.NewReader(`{}`))
	response := httptest.NewRecorder()

	handler.Attest(response, request)

	if response.Code != http.StatusServiceUnavailable || usecase.attestCalls != 0 {
		t.Fatalf("Attest() status/calls = %d/%d, want 503/0", response.Code, usecase.attestCalls)
	}
}

func TestAttestRequiresMatchingBearerToken(t *testing.T) {
	t.Setenv(AttestAdminTokenEnv, "demo-admin-token")
	usecase := &stubAttestationUsecase{}
	handler := NewAttestationHandler(usecase)
	request := httptest.NewRequest(http.MethodPost, "/attestations", strings.NewReader(`{}`))
	request.Header.Set("Authorization", "Bearer wrong-token")
	response := httptest.NewRecorder()

	handler.Attest(response, request)

	if response.Code != http.StatusUnauthorized || usecase.attestCalls != 0 {
		t.Fatalf("Attest() status/calls = %d/%d, want 401/0", response.Code, usecase.attestCalls)
	}
}

func TestAttestAcceptsValidBearerToken(t *testing.T) {
	t.Setenv(AttestAdminTokenEnv, "demo-admin-token")
	usecase := &stubAttestationUsecase{
		attestResult: attestation.Result{Status: attestation.StatusMatch, UID: [32]byte{1}},
	}
	handler := NewAttestationHandler(usecase)
	request := httptest.NewRequest(
		http.MethodPost,
		"/attestations",
		strings.NewReader(`{"user_id":"00000000-0000-4000-8000-000000000101","date":"2026-11-10"}`),
	)
	request.Header.Set("Authorization", "Bearer demo-admin-token")
	response := httptest.NewRecorder()

	handler.Attest(response, request)

	if response.Code != http.StatusCreated || usecase.attestCalls != 1 {
		t.Fatalf("Attest() status/calls = %d/%d, want 201/1; body=%s", response.Code, usecase.attestCalls, response.Body)
	}
}

func TestVerifyReturnsNotAttestedWithoutWriteToken(t *testing.T) {
	usecase := &stubAttestationUsecase{
		verifyResult: attestation.Result{Status: attestation.StatusNotAttested},
	}
	handler := NewAttestationHandler(usecase)
	request := httptest.NewRequest(
		http.MethodGet,
		"/attestations/verify?user_id=00000000-0000-4000-8000-000000000101&date=2026-11-10",
		nil,
	)
	response := httptest.NewRecorder()

	handler.Verify(response, request)

	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"not-attested"`) {
		t.Fatalf("Verify() status/body = %d/%s", response.Code, response.Body)
	}
}
