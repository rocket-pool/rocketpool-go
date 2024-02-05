package tests

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nodeset-org/eth-utils/eth"
)

// An account containing an address and a transactor for it
type Account struct {
	Address    common.Address
	Transactor *bind.TransactOpts
}

// Get an account by index
func CreateAccountFromPrivateKey(privateKeyHex string, chainID *big.Int) (*Account, error) {
	// Get private key data
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// Get private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	// Get the account transactor
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("error creating transactor: %w", err)
	}
	opts.Context = context.Background()
	opts.GasFeeCap = eth.GweiToWei(10)
	opts.GasTipCap = eth.GweiToWei(2)
	//opts.GasPrice = eth.GweiToWei(10)

	// Return account
	return &Account{
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		Transactor: opts,
	}, nil
}
