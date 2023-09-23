package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/types"
)

type CallReturnType interface {
	*big.Int | uint8 | bool | string | common.Address | common.Hash | types.ValidatorPubkey | []byte
}

type FormattedUint8Type interface {
	types.MinipoolStatus | types.MinipoolDeposit | types.ProposalState
}

// Contract type wraps go-ethereum bound contract
type Contract struct {
	Name     string
	Contract *bind.BoundContract
	Address  *common.Address
	ABI      *abi.ABI
	Version  uint8
	Client   ExecutionClient
}

// Response for gas limits from network and from user request
type GasInfo struct {
	EstGasLimit  uint64 `json:"estGasLimit"`
	SafeGasLimit uint64 `json:"safeGasLimit"`
}

// Call a contract method
func (c *Contract) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	results := make([]interface{}, 1)
	results[0] = result
	return c.Contract.Call(opts, &results, method, params...)
}

// Calls a contract method
func Call[retType CallReturnType](contract *Contract, opts *bind.CallOpts, method string, params ...interface{}) (retType, error) {
	// Set up the return capture
	result := new(retType)
	results := make([]interface{}, 1)
	results[0] = result

	// Run the function
	err := contract.Call(opts, &results, method, params...)
	return *result, err
}
