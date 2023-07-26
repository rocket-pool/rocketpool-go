package rocketpool

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Information of a candidate transaction
type TransactionInfo struct {
	Data     string  `json:"data"`
	Nonce    uint64  `json:"nonce"`
	GasInfo  GasInfo `json:"gasInfo"`
	SimError string  `json:"simError"`
}

// Create a new serializable TransactionInfo wrapper
func NewTransactionInfo(contract *Contract, method string, opts *bind.TransactOpts, parameters ...interface{}) (*TransactionInfo, error) {
	// Create the input data
	input, err := contract.ABI.Pack(method, parameters...)
	if err != nil {
		return nil, fmt.Errorf("error packing input data: %w", err)
	}

	// Get the gas estimate
	gasInfo, simErr := contract.GetTransactionGasInfo(opts, method, parameters)

	// Serialize the data
	dataString := hex.EncodeToString(input)

	// Create the info wrapper
	txInfo := &TransactionInfo{
		Data:     dataString,
		Nonce:    opts.Nonce.Uint64(),
		GasInfo:  gasInfo,
		SimError: "",
	}
	if simErr != nil {
		txInfo.SimError = simErr.Error()
	}

	return txInfo, nil
}
