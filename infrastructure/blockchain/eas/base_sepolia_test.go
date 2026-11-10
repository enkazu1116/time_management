package eas

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"time_management/domain/time_tracking/attestation"
)

func TestLoadEnvironmentDoesNotLeakPrivateKey(t *testing.T) {
	const privateKey = "private-key-value-that-must-not-leak"
	t.Setenv(RPCURLEnv, "https://example.invalid")
	t.Setenv(PrivateKeyEnv, privateKey)
	t.Setenv(SchemaUIDEnv, "0x"+strings.Repeat("01", 32))

	var logs bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(originalOutput) })

	_, err := loadEnvironment()
	if !errors.Is(err, ErrPrivateKeyInvalid) {
		t.Fatalf("loadEnvironment() error = %v, want %v", err, ErrPrivateKeyInvalid)
	}
	if strings.Contains(err.Error(), privateKey) {
		t.Fatal("loadEnvironment() returned the private key in its error")
	}
	if strings.Contains(logs.String(), privateKey) {
		t.Fatal("loadEnvironment() logged the private key")
	}
}

func TestEncodeSchemaDataUsesTwoBytes32Fields(t *testing.T) {
	payload := attestation.Payload{
		DayDigest:         [32]byte{1, 2},
		SubjectCommitment: [32]byte{3, 4},
	}
	got, err := encodeSchemaData(payload)
	if err != nil {
		t.Fatal(err)
	}
	want := append(payload.DayDigest[:], payload.SubjectCommitment[:]...)
	if !bytes.Equal(got, want) {
		t.Fatalf("encodeSchemaData() = %x, want %x", got, want)
	}
}

func TestRequestCarriesSchemaDataAndReferenceUID(t *testing.T) {
	schema := [32]byte{1}
	refUID := [32]byte{2}
	payload := attestation.Payload{RefUID: refUID}
	data := []byte{3, 4}
	chain := &BaseSepolia{schemaUID: schema}

	request := chain.request(payload, data)
	if request.Schema != schema || request.Data.RefUID != refUID {
		t.Fatalf("request did not preserve schema/refUID: %#v", request)
	}
	if request.Data.Recipient != common.HexToAddress(zeroAddress) ||
		!request.Data.Revocable ||
		request.Data.ExpirationTime != 0 ||
		request.Data.Value.Sign() != 0 ||
		!bytes.Equal(request.Data.Data, data) {
		t.Fatalf("request has unexpected EAS fields: %#v", request.Data)
	}
}

func TestAttestationUIDAcceptsOnlyExpectedContractAttesterAndSchema(t *testing.T) {
	contract, err := newEASContract(common.HexToAddress(EASAddress), nil)
	if err != nil {
		t.Fatal(err)
	}
	attester := common.HexToAddress("0x0000000000000000000000000000000000001234")
	schema := [32]byte{1}
	uid := [32]byte{2}
	event := contract.abi.Events["Attested"]
	eventLog := &types.Log{
		Address: common.HexToAddress(EASAddress),
		Topics: []common.Hash{
			event.ID,
			common.BytesToHash(common.HexToAddress(zeroAddress).Bytes()),
			common.BytesToHash(attester.Bytes()),
			common.BytesToHash(schema[:]),
		},
		Data: uid[:],
	}
	chain := &BaseSepolia{contract: contract, schemaUID: schema, attester: attester}

	got, err := chain.attestationUID([]*types.Log{eventLog})
	if err != nil {
		t.Fatal(err)
	}
	if got != uid {
		t.Fatalf("attestationUID() = %x, want %x", got, uid)
	}

	eventLog.Topics[2] = common.HexToHash("0xffff")
	if _, err := chain.attestationUID([]*types.Log{eventLog}); !errors.Is(err, ErrAttestedEvent) {
		t.Fatalf("attestationUID() error = %v, want %v", err, ErrAttestedEvent)
	}
}
