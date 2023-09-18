package network

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNetworkPenalties
type NetworkPenalties struct {
	rp *rocketpool.RocketPool
	np *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new NetworkPenalties contract binding
func NewNetworkPenalties(rp *rocketpool.RocketPool) (*NetworkPenalties, error) {
	// Create the contract
	np, err := rp.GetContract(rocketpool.ContractName_RocketNetworkPenalties)
	if err != nil {
		return nil, fmt.Errorf("error getting network penalties contract: %w", err)
	}

	return &NetworkPenalties{
		rp: rp,
		np: np,
	}, nil
}

// ====================
// === Transactions ===
// ====================

// Get info for minipool penalty submission
func (c *NetworkPenalties) SubmitPenalty(minipoolAddress common.Address, block *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.np, "submitPenalty", opts, minipoolAddress, block)
}
