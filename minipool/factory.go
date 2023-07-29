package minipool

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketMinipoolFactory
type MinipoolFactory struct {
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new MinipoolFactory contract binding
func NewMinipoolFactory(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*MinipoolFactory, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketMinipoolFactory", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool factory contract: %w", err)
	}

	return &MinipoolFactory{
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the address of a minipool based on the node address and a salt
func (c *MinipoolFactory) GetExpectedAddress(mc *multicall.MultiCaller, nodeAddress common.Address, salt *big.Int, address_Out *common.Address) {
	multicall.AddCall(mc, c.contract, address_Out, "getExpectedAddress", nodeAddress, salt)
}
