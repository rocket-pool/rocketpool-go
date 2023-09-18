package node

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNodeStaking
type NodeStaking struct {
	*NodeStakingDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for RocketNodeStaking
type NodeStakingDetails struct {
	Version       uint8    `json:"version"`
	TotalRplStake *big.Int `json:"TotalRplStake"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NodeStaking contract binding
func NewNodeStaking(rp *rocketpool.RocketPool) (*NodeStaking, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketNodeStaking)
	if err != nil {
		return nil, fmt.Errorf("error getting node staking contract: %w", err)
	}

	return &NodeStaking{
		NodeStakingDetails: &NodeStakingDetails{},
		rp:                 rp,
		contract:           contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the version of the Node Staking contract
func (c *NodeStaking) GetVersion(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Version, "version")
}

// Get the total RPL staked in the network
func (c *NodeStaking) GetTotalRPLStake(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.TotalRplStake, "getTotalRPLStake")
}

// Get all basic details
func (c *NodeStaking) GetAllDetails(mc *batch.MultiCaller) {
	c.GetVersion(mc)
	c.GetTotalRPLStake(mc)
}

// =============
// === Utils ===
// =============

// Calculate total effective RPL stake for the subset of nodes provided
// NOTE: you will have to call this several times, iterating through subset ranges,
// to get the complete result since a single call for the complete node set may run out of gas
func (c *NodeStaking) CalculateTotalEffectiveRplStake(offset *big.Int, limit *big.Int, rplPrice *big.Int, opts *bind.CallOpts) (*big.Int, error) {
	return core.Call[*big.Int](c.contract, opts, "calculateTotalEffectiveRPLStake", offset, limit, rplPrice)
}
