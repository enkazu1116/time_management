package eas

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"time_management/domain/time_tracking/attestation"
)

const (
	RPCURLEnv     = "BASE_SEPOLIA_RPC_URL"
	PrivateKeyEnv = "EAS_PRIVATE_KEY"
	SchemaUIDEnv  = "EAS_SCHEMA_UID"
	BaseSepoliaID = int64(84532)
	EASAddress    = "0x4200000000000000000000000000000000000021"
	zeroAddress   = "0x0000000000000000000000000000000000000000"
)

var (
	ErrRPCURLRequired     = errors.New("eas: BASE_SEPOLIA_RPC_URL is required")
	ErrPrivateKeyRequired = errors.New("eas: EAS_PRIVATE_KEY is required")
	ErrPrivateKeyInvalid  = errors.New("eas: EAS_PRIVATE_KEY is invalid")
	ErrSchemaUIDRequired  = errors.New("eas: EAS_SCHEMA_UID is required")
	ErrSchemaUIDInvalid   = errors.New("eas: EAS_SCHEMA_UID is invalid")
	ErrRPCConnection      = errors.New("eas: Base Sepolia RPC connection failed")
	ErrWrongChain         = errors.New("eas: RPC is not Base Sepolia")
	ErrContractSetup      = errors.New("eas: contract setup failed")
	ErrTransaction        = errors.New("eas: attestation transaction failed")
	ErrReceipt            = errors.New("eas: attestation receipt failed")
	ErrAttestedEvent      = errors.New("eas: valid Attested event not found")
	ErrAttestationRead    = errors.New("eas: attestation read failed")
	ErrAttestationInvalid = errors.New("eas: attestation does not match configured identity")
	ErrAttestationRevoked = errors.New("eas: attestation is revoked")
)

type environmentConfig struct {
	rpcURL     string
	privateKey *ecdsa.PrivateKey
	schemaUID  [32]byte
}

// BaseSepolia is the production Chain implementation. It is intentionally
// fixed to Base Sepolia's chain ID and canonical EAS deployment.
type BaseSepolia struct {
	client        *ethclient.Client
	contract      *easContract
	schemaUID     [32]byte
	attester      common.Address
	privateKey    *ecdsa.PrivateKey
	transactionMu sync.Mutex
}

// NewBaseSepoliaFromEnv loads the RPC URL, signer key, and schema UID from
// environment variables. Parsing errors never include environment values.
func NewBaseSepoliaFromEnv(ctx context.Context) (*BaseSepolia, error) {
	config, err := loadEnvironment()
	if err != nil {
		return nil, err
	}

	client, err := ethclient.DialContext(ctx, config.rpcURL)
	if err != nil {
		return nil, ErrRPCConnection
	}
	chainID, err := client.ChainID(ctx)
	if err != nil {
		client.Close()
		return nil, ErrRPCConnection
	}
	if chainID.Cmp(big.NewInt(BaseSepoliaID)) != 0 {
		client.Close()
		return nil, ErrWrongChain
	}

	contract, err := newEASContract(common.HexToAddress(EASAddress), client)
	if err != nil {
		client.Close()
		return nil, ErrContractSetup
	}
	return &BaseSepolia{
		client:     client,
		contract:   contract,
		schemaUID:  config.schemaUID,
		attester:   crypto.PubkeyToAddress(config.privateKey.PublicKey),
		privateKey: config.privateKey,
	}, nil
}

func loadEnvironment() (environmentConfig, error) {
	var config environmentConfig
	config.rpcURL = os.Getenv(RPCURLEnv)
	if config.rpcURL == "" {
		return environmentConfig{}, ErrRPCURLRequired
	}

	privateKeyValue := os.Getenv(PrivateKeyEnv)
	if privateKeyValue == "" {
		return environmentConfig{}, ErrPrivateKeyRequired
	}
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyValue, "0x"))
	if err != nil {
		return environmentConfig{}, ErrPrivateKeyInvalid
	}
	config.privateKey = privateKey

	schemaValue := os.Getenv(SchemaUIDEnv)
	if schemaValue == "" {
		return environmentConfig{}, ErrSchemaUIDRequired
	}
	schemaUID, ok := parseBytes32(schemaValue)
	if !ok {
		return environmentConfig{}, ErrSchemaUIDInvalid
	}
	config.schemaUID = schemaUID
	return config, nil
}

func parseBytes32(value string) ([32]byte, bool) {
	value = strings.TrimPrefix(value, "0x")
	if len(value) != 64 {
		return [32]byte{}, false
	}
	decoded := common.FromHex("0x" + value)
	if len(decoded) != 32 {
		return [32]byte{}, false
	}
	var result [32]byte
	copy(result[:], decoded)
	return result, true
}

func (chain *BaseSepolia) Attest(ctx context.Context, payload attestation.Payload) ([32]byte, error) {
	// bind obtains the pending nonce when Transact is called. Serialize writes
	// for this signer so concurrent requests cannot select the same nonce.
	chain.transactionMu.Lock()
	defer chain.transactionMu.Unlock()

	data, err := encodeSchemaData(payload)
	if err != nil {
		return [32]byte{}, ErrTransaction
	}
	opts, err := bind.NewKeyedTransactorWithChainID(chain.privateKey, big.NewInt(BaseSepoliaID))
	if err != nil {
		return [32]byte{}, ErrTransaction
	}
	opts.Context = ctx

	transaction, err := chain.contract.attest(opts, chain.request(payload, data))
	if err != nil {
		return [32]byte{}, ErrTransaction
	}
	receipt, err := bind.WaitMined(ctx, chain.client, transaction)
	if err != nil || receipt.Status == 0 {
		return [32]byte{}, ErrReceipt
	}
	return chain.attestationUID(receipt.Logs)
}

func (chain *BaseSepolia) request(payload attestation.Payload, data []byte) attestationRequest {
	return attestationRequest{
		Schema: chain.schemaUID,
		Data: attestationRequestData{
			Recipient: common.HexToAddress(zeroAddress),
			Revocable: true,
			RefUID:    payload.RefUID,
			Data:      data,
			Value:     big.NewInt(0),
		},
	}
}

func encodeSchemaData(payload attestation.Payload) ([]byte, error) {
	bytes32Type, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, err
	}
	arguments := abi.Arguments{{Type: bytes32Type}, {Type: bytes32Type}}
	return arguments.Pack(payload.DayDigest, payload.SubjectCommitment)
}

func (chain *BaseSepolia) Read(ctx context.Context, uid [32]byte) (attestation.OnChainRecord, error) {
	record, err := chain.contract.getAttestation(&bind.CallOpts{Context: ctx}, uid)
	if err != nil {
		return attestation.OnChainRecord{}, ErrAttestationRead
	}
	return chain.validateRecord(uid, record)
}

func (chain *BaseSepolia) validateRecord(uid [32]byte, record attestationRecord) (attestation.OnChainRecord, error) {
	if record.UID != uid || record.UID == ([32]byte{}) ||
		record.Schema != chain.schemaUID ||
		record.Attester != chain.attester {
		return attestation.OnChainRecord{}, ErrAttestationInvalid
	}
	if record.RevocationTime != 0 {
		return attestation.OnChainRecord{}, ErrAttestationRevoked
	}
	dayDigest, subjectCommitment, err := decodeSchemaData(record.Data)
	if err != nil {
		return attestation.OnChainRecord{}, ErrAttestationInvalid
	}
	return attestation.OnChainRecord{
		DayDigest:         dayDigest,
		SubjectCommitment: subjectCommitment,
	}, nil
}

func decodeSchemaData(data []byte) ([32]byte, [32]byte, error) {
	bytes32Type, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return [32]byte{}, [32]byte{}, err
	}
	values, err := (abi.Arguments{{Type: bytes32Type}, {Type: bytes32Type}}).Unpack(data)
	if err != nil || len(values) != 2 {
		return [32]byte{}, [32]byte{}, ErrAttestationInvalid
	}
	digest, digestOK := values[0].([32]byte)
	commitment, commitmentOK := values[1].([32]byte)
	if !digestOK || !commitmentOK {
		return [32]byte{}, [32]byte{}, ErrAttestationInvalid
	}
	return digest, commitment, nil
}

func (chain *BaseSepolia) attestationUID(logs []*types.Log) ([32]byte, error) {
	event := chain.contract.abi.Events["Attested"]
	for _, eventLog := range logs {
		if eventLog.Address != common.HexToAddress(EASAddress) ||
			len(eventLog.Topics) != 4 ||
			eventLog.Topics[0] != event.ID ||
			common.BytesToAddress(eventLog.Topics[2].Bytes()) != chain.attester ||
			eventLog.Topics[3] != common.BytesToHash(chain.schemaUID[:]) {
			continue
		}
		values, err := event.Inputs.NonIndexed().Unpack(eventLog.Data)
		if err != nil || len(values) != 1 {
			continue
		}
		uid, ok := values[0].([32]byte)
		if ok {
			return uid, nil
		}
	}
	return [32]byte{}, ErrAttestedEvent
}

func (chain *BaseSepolia) Close() {
	chain.client.Close()
}

var _ attestation.Chain = (*BaseSepolia)(nil)
