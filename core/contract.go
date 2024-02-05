package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/eth-utils/eth"
)

type CallReturnType interface {
	*big.Int | uint8 | bool | string | common.Address | common.Hash | beacon.ValidatorPubkey | []byte
}

// Contract type wraps go-ethereum bound contract
type Contract struct {
	*eth.Contract
	Version uint8
}

// Call a contract method
func (c *Contract) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	results := make([]interface{}, 1)
	results[0] = result
	return c.ContractImpl.Call(opts, &results, method, params...)
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
