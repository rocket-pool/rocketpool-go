package network

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNetworkPenalties
type NetworkPenalties struct {
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new NetworkPenalties contract binding
func NewNetworkPenalties(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*NetworkPenalties, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketNetworkPenalties", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting network penalties contract: %w", err)
	}

	return &NetworkPenalties{
		rp:       rp,
		contract: contract,
	}, nil
}

// ====================
// === Transactions ===
// ====================

// Get info for minipool penalty submission
func (c *NetworkPenalties) SubmitPenalty(minipoolAddress common.Address, block *big.Int, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, "submitPenalty", opts, minipoolAddress, block)
}
