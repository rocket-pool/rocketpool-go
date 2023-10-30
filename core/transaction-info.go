package core

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const (
	GasLimitMultiplier    float64 = 1.5
	MaxGasLimit           uint64  = 30000000
	NethermindRevertRegex string  = "Reverted 0x(?P<message>[0-9a-fA-F]+).*"

	gasSimErrorPrefix string = "error estimating gas needed"
)

// Information of a candidate transaction
type TransactionInfo struct {
	Data     []byte         `json:"data"`
	To       common.Address `json:"to"`
	Value    *big.Int       `json:"value"`
	GasInfo  GasInfo        `json:"gasInfo"`
	SimError string         `json:"simError"`
}

// Create a new serializable TransactionInfo wrapper
func NewTransactionInfo(contract *Contract, method string, opts *bind.TransactOpts, parameters ...interface{}) (*TransactionInfo, error) {
	// Create the data data
	data, err := contract.ABI.Pack(method, parameters...)
	if err != nil {
		return nil, fmt.Errorf("error packing input data: %w", err)
	}

	// Get the gas estimate
	estGasLimit, safeGasLimit, simErr := estimateGasLimit(contract.Client, *contract.Address, opts, data)
	if simErr != nil && !strings.HasPrefix(simErr.Error(), gasSimErrorPrefix) {
		return nil, simErr
	}

	// Create the info wrapper
	var value *big.Int
	if opts != nil {
		value = opts.Value
	}
	txInfo := &TransactionInfo{
		Data:  data,
		To:    *contract.Address,
		Value: value,
		GasInfo: GasInfo{
			EstGasLimit:  estGasLimit,
			SafeGasLimit: safeGasLimit,
		},
		SimError: "",
	}
	if simErr != nil {
		txInfo.SimError = simErr.Error()
	}

	return txInfo, nil
}

// Create a new serializable TransactionInfo from raw data
func NewTransactionInfoRaw(ec ExecutionClient, to common.Address, data []byte, opts *bind.TransactOpts) (*TransactionInfo, error) {
	// Get the gas estimate
	estGasLimit, safeGasLimit, simErr := estimateGasLimit(ec, to, opts, data)
	if simErr != nil && !strings.HasPrefix(simErr.Error(), gasSimErrorPrefix) {
		return nil, simErr
	}

	// Create the info wrapper
	var value *big.Int
	if opts != nil {
		value = opts.Value
	}
	txInfo := &TransactionInfo{
		Data:  data,
		To:    to,
		Value: value,
		GasInfo: GasInfo{
			EstGasLimit:  estGasLimit,
			SafeGasLimit: safeGasLimit,
		},
		SimError: "",
	}
	if simErr != nil {
		txInfo.SimError = simErr.Error()
	}

	return txInfo, nil
}

// Estimate the expected and safe gas limits for a contract transaction
func estimateGasLimit(client ExecutionClient, to common.Address, opts *bind.TransactOpts, input []byte) (uint64, uint64, error) {

	// Estimate gas limit
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:      opts.From,
		To:        &to,
		GasFeeCap: big.NewInt(0),
		GasTipCap: big.NewInt(0),
		Value:     opts.Value,
		Data:      input,
	})

	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", gasSimErrorPrefix, normalizeErrorMessage(err))
	}

	// Pad and return gas limit
	safeGasLimit := uint64(float64(gasLimit) * GasLimitMultiplier)
	if gasLimit > MaxGasLimit {
		return 0, 0, fmt.Errorf("estimated gas of %d is greater than the max gas limit of %d", gasLimit, MaxGasLimit)
	}
	if safeGasLimit > MaxGasLimit {
		safeGasLimit = MaxGasLimit
	}
	return gasLimit, safeGasLimit, nil

}

// Normalize error messages so they're all in ASCII format
func normalizeErrorMessage(err error) error {
	if err == nil {
		return err
	}

	// Get the message in hex format, if it exists
	reg := regexp.MustCompile(NethermindRevertRegex)
	matches := reg.FindStringSubmatch(err.Error())
	if matches == nil {
		return err
	}
	messageIndex := reg.SubexpIndex("message")
	if messageIndex == -1 {
		return err
	}
	message := matches[messageIndex]

	// Convert the hex message to ASCII
	bytes, err2 := hex.DecodeString(message)
	if err2 != nil {
		return err // Return the original error if decoding failed somehow
	}

	return fmt.Errorf("Reverted: %s", string(bytes))
}
