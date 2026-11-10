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

const easABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"recipient","type":"address"},{"indexed":true,"internalType":"address","name":"attester","type":"address"},{"indexed":false,"internalType":"bytes32","name":"uid","type":"bytes32"},{"indexed":true,"internalType":"bytes32","name":"schema","type":"bytes32"}],"name":"Attested","type":"event"},{"inputs":[{"components":[{"internalType":"bytes32","name":"schema","type":"bytes32"},{"components":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint64","name":"expirationTime","type":"uint64"},{"internalType":"bool","name":"revocable","type":"bool"},{"internalType":"bytes32","name":"refUID","type":"bytes32"},{"internalType":"bytes","name":"data","type":"bytes"},{"internalType":"uint256","name":"value","type":"uint256"}],"internalType":"struct AttestationRequestData","name":"data","type":"tuple"}],"internalType":"struct AttestationRequest","name":"request","type":"tuple"}],"name":"attest","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"payable","type":"function"}]`

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
