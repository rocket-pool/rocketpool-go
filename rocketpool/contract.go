package rocketpool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Transaction settings
const ()

// Contract type wraps go-ethereum bound contract
type Contract struct {
	Contract *bind.BoundContract
	Address  *common.Address
	ABI      *abi.ABI
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
func Call[retType *big.Int | uint8 | bool | string](contract *Contract, opts *bind.CallOpts, method string, params ...interface{}) (retType, error) {
	// Set up the return capture
	result := new(retType)
	results := make([]interface{}, 1)
	results[0] = result

	// Run the function
	err := contract.Call(opts, &results, method, params...)
	return *result, err
}

// Calls a contract method for a parameter
func CallForParameter[FormattedType time.Time | uint64 | float64](contract *Contract, opts *bind.CallOpts, method string, params ...interface{}) (Parameter[FormattedType], error) {
	// Set up the return capture
	result := new(*big.Int)
	results := make([]interface{}, 1)
	results[0] = result

	// Run the function
	var param Parameter[FormattedType]
	err := contract.Call(opts, &results, method, params...)
	if err != nil {
		return param, err
	}

	// Wrap and return
	param.RawValue = *result
	return param, err
}

// Transact on a contract method and wait for a receipt
func (c *Contract) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {

	// Estimate gas limit
	if opts.GasLimit == 0 {
		input, err := c.ABI.Pack(method, params...)
		if err != nil {
			return nil, fmt.Errorf("Could not encode input data: %w", err)
		}
		_, safeGasLimit, err := estimateGasLimit(c.Client, *c.Address, opts, input)
		if err != nil {
			return nil, err
		}
		opts.GasLimit = safeGasLimit
	}

	// Send transaction
	tx, err := c.Contract.Transact(opts, method, params...)
	if err != nil {
		return nil, normalizeErrorMessage(err)
	}

	return tx, nil

}

// Get gas limit for a transfer call
func (c *Contract) GetTransferGasInfo(opts *bind.TransactOpts) (GasInfo, error) {

	response := GasInfo{}

	// Estimate gas limit
	estGasLimit, safeGasLimit, err := estimateGasLimit(c.Client, *c.Address, opts, []byte{})
	if err != nil {
		return response, fmt.Errorf("Error getting transfer gas info: could not estimate gas limit: %w", err)
	}
	response.EstGasLimit = estGasLimit
	response.SafeGasLimit = safeGasLimit

	return response, nil
}

// Transfer ETH to a contract and wait for a receipt
func (c *Contract) Transfer(opts *bind.TransactOpts) (common.Hash, error) {

	// Estimate gas limit
	if opts.GasLimit == 0 {
		_, safeGasLimit, err := estimateGasLimit(c.Client, *c.Address, opts, []byte{})
		if err != nil {
			return common.Hash{}, err
		}
		opts.GasLimit = safeGasLimit
	}

	// Send transaction
	tx, err := c.Contract.Transfer(opts)
	if err != nil {
		return common.Hash{}, normalizeErrorMessage(err)
	}

	return tx.Hash(), nil

}

// Wait for a transaction to be mined and get a tx receipt
func (c *Contract) getTransactionReceipt(tx *types.Transaction) (*types.Receipt, error) {

	// Wait for transaction to be mined
	txReceipt, err := bind.WaitMined(context.Background(), c.Client, tx)
	if err != nil {
		return nil, err
	}

	// Check transaction status
	if txReceipt.Status == 0 {
		return txReceipt, errors.New("Transaction failed with status 0")
	}

	// Return
	return txReceipt, nil

}

// Get contract events from a transaction
// eventPrototype must be an event struct type
// Returns a slice of untyped values; assert returned events to event struct type
func (c *Contract) GetTransactionEvents(txReceipt *types.Receipt, eventName string, eventPrototype interface{}) ([]interface{}, error) {

	// Get event type
	eventType := reflect.TypeOf(eventPrototype)
	if eventType.Kind() != reflect.Struct {
		return nil, errors.New("Invalid event type")
	}

	// Get ABI event
	abiEvent, ok := c.ABI.Events[eventName]
	if !ok {
		return nil, fmt.Errorf("Event '%s' does not exist on contract", eventName)
	}

	// Process transaction receipt logs
	events := make([]interface{}, 0)
	for _, log := range txReceipt.Logs {

		// Check log address matches contract address
		if !bytes.Equal(log.Address.Bytes(), c.Address.Bytes()) {
			continue
		}

		// Check log first topic matches event ID
		if len(log.Topics) == 0 || !bytes.Equal(log.Topics[0].Bytes(), abiEvent.ID.Bytes()) {
			continue
		}

		// Unpack event
		event := reflect.New(eventType)
		if err := c.Contract.UnpackLog(event.Interface(), eventName, *log); err != nil {
			return nil, fmt.Errorf("Could not unpack event data: %w", err)
		}
		events = append(events, reflect.Indirect(event).Interface())

	}

	// Return events
	return events, nil

}
