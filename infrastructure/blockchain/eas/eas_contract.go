// Code generated from the official EAS ABI. DO NOT EDIT.

package eas

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const easABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"recipient","type":"address"},{"indexed":true,"internalType":"address","name":"attester","type":"address"},{"indexed":false,"internalType":"bytes32","name":"uid","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"schema","type":"bytes32"}],"name":"Attested","type":"event"},{"inputs":[{"components":[{"internalType":"bytes32","name":"schema","type":"bytes32"},{"components":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint64","name":"expirationTime","type":"uint64"},{"internalType":"bool","name":"revocable","type":"bool"},{"internalType":"bytes32","name":"refUID","type":"bytes32"},{"internalType":"bytes","name":"data","type":"bytes"},{"internalType":"uint256","name":"value","type":"uint256"}],"internalType":"struct AttestationRequestData","name":"data","type":"tuple"}],"internalType":"struct AttestationRequest","name":"request","type":"tuple"}],"name":"attest","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"uid","type":"bytes32"}],"name":"getAttestation","outputs":[{"components":[{"internalType":"bytes32","name":"uid","type":"bytes32"},{"internalType":"bytes32","name":"schema","type":"bytes32"},{"internalType":"uint64","name":"time","type":"uint64"},{"internalType":"uint64","name":"expirationTime","type":"uint64"},{"internalType":"uint64","name":"revocationTime","type":"uint64"},{"internalType":"bytes32","name":"refUID","type":"bytes32"},{"internalType":"address","name":"recipient","type":"address"},{"internalType":"address","name":"attester","type":"address"},{"internalType":"bool","name":"revocable","type":"bool"},{"internalType":"bytes","name":"data","type":"bytes"}],"internalType":"struct Attestation","name":"","type":"tuple"}],"stateMutability":"view","type":"function"}]`

type attestationRequestData struct {
	Recipient      common.Address
	ExpirationTime uint64
	Revocable      bool
	RefUID         [32]byte
	Data           []byte
	Value          *big.Int
}

type attestationRequest struct {
	Schema [32]byte
	Data   attestationRequestData
}

type attestationRecord struct {
	UID            [32]byte
	Schema         [32]byte
	Time           uint64
	ExpirationTime uint64
	RevocationTime uint64
	RefUID         [32]byte
	Recipient      common.Address
	Attester       common.Address
	Revocable      bool
	Data           []byte
}

type easContract struct {
	contract *bind.BoundContract
	abi      abi.ABI
}

func newEASContract(address common.Address, backend bind.ContractBackend) (*easContract, error) {
	parsed, err := abi.JSON(strings.NewReader(easABI))
	if err != nil {
		return nil, err
	}
	return &easContract{
		contract: bind.NewBoundContract(address, parsed, backend, backend, backend),
		abi:      parsed,
	}, nil
}

func (contract *easContract) attest(opts *bind.TransactOpts, request attestationRequest) (*types.Transaction, error) {
	return contract.contract.Transact(opts, "attest", request)
}

func (contract *easContract) getAttestation(opts *bind.CallOpts, uid [32]byte) (attestationRecord, error) {
	var out []interface{}
	err := contract.contract.Call(opts, &out, "getAttestation", uid)
	if err != nil {
		return attestationRecord{}, err
	}
	return *abi.ConvertType(out[0], new(attestationRecord)).(*attestationRecord), nil
}
