package core

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Transaction settings
const ()

type CallReturnType interface {
	*big.Int | uint8 | bool | string | common.Address | common.Hash
}

type FormattedType interface {
	time.Time | uint64 | float64 | time.Duration
}

// Contract type wraps go-ethereum bound contract
type Contract struct {
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

// Calls a contract method without type safety on the return type (useful for custom structs)
func CallUntyped[retType any](contract *Contract, opts *bind.CallOpts, method string, params ...interface{}) (retType, error) {
	// Set up the return capture
	result := new(retType)
	results := make([]interface{}, 1)
	results[0] = result

	// Run the function
	err := contract.Call(opts, &results, method, params...)
	return *result, err
}

// Calls a contract method that returns a slice
func CallForSlice[retType CallReturnType](contract *Contract, opts *bind.CallOpts, method string, params ...interface{}) ([]retType, error) {
	// Set up the return capture
	result := new([]retType)
	results := make([]interface{}, 1)
	results[0] = result

	// Run the function
	err := contract.Call(opts, &results, method, params...)
	return *result, err
}

// Calls a contract method for a parameter
func CallForParameter[fType FormattedType](contract *Contract, opts *bind.CallOpts, method string, params ...interface{}) (Parameter[fType], error) {
	// Set up the return capture
	result := new(*big.Int)
	results := make([]interface{}, 1)
	results[0] = result

	// Run the function
	var param Parameter[fType]
	err := contract.Call(opts, &results, method, params...)
	if err != nil {
		return param, err
	}

	// Wrap and return
	param.RawValue = *result
	return param, err
}
