package responses

import (
	"encoding/hex"

	"time_management/domain/time_tracking/attestation"
)

type AttestationResponse struct {
	Status         string `json:"status"`
	AttestationUID string `json:"attestation_uid,omitempty"`
}

func Attestation(result attestation.Result) AttestationResponse {
	response := AttestationResponse{Status: result.Status}
	if result.UID != ([32]byte{}) {
		response.AttestationUID = "0x" + hex.EncodeToString(result.UID[:])
	}
	return response
}
