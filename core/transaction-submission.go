package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Information for submitting a candidate transaction to the network
type TransactionSubmission struct {
	TxInfo   *TransactionInfo `json:"txInfo"`
	GasLimit uint64           `json:"gasLimit"`
}

// Create a transaction from serialized info, signs it, and submits it to the network if requested in opts.
// Note the value in opts is not used; set it in the value argument instead.
func ExecuteTransaction(client ExecutionClient, data []byte, to common.Address, value *big.Int, opts *bind.TransactOpts) (*types.Transaction, error) {
	// Create a "dummy" contract for the Geth API with no ABI since we don't need it for this
	contract := bind.NewBoundContract(to, abi.ABI{}, client, client, client)

	newOpts := &bind.TransactOpts{
		// Copy the original fields
		From:      opts.From,
		Nonce:     opts.Nonce,
		Signer:    opts.Signer,
		GasPrice:  opts.GasPrice,
		GasFeeCap: opts.GasFeeCap,
		GasTipCap: opts.GasTipCap,
		GasLimit:  opts.GasLimit,
		Context:   opts.Context,
		NoSend:    opts.NoSend,

		// Overwrite the value
		Value: value,
	}

	return contract.RawTransact(newOpts, data)
}

// Create a transaction submission directly from serialized info (and the error provided by the transaction info constructor),
// using the SafeGasLimit as the GasLimit for the submission automatically.
func CreateTxSubmissionFromInfo(txInfo *TransactionInfo, err error) (*TransactionSubmission, error) {
	if err != nil {
		return nil, err
	}
	return &TransactionSubmission{
		TxInfo:   txInfo,
		GasLimit: txInfo.GasInfo.SafeGasLimit,
	}, nil
}
