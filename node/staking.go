package node

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNodeStaking
type NodeStaking struct {
	Details  NodeStakingDetails
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
		Details:  NodeStakingDetails{},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the version of the Node Staking contract
func (c *NodeStaking) GetVersion(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.Version, "version")
}

// Get the total RPL staked in the network
func (c *NodeStaking) GetTotalRPLStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.TotalRplStake, "getTotalRPLStake")
}

// Get all basic details
func (c *NodeStaking) GetAllDetails(mc *multicall.MultiCaller) {
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
