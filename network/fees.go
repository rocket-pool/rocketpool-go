package network

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNetworkFees
type NetworkFees struct {
	Details  NetworkFeesDetails
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// Details for network fees
type NetworkFeesDetails struct {
	NodeDemand      *big.Int                      `json:"nodeDemand"`
	NodeFee         rocketpool.Parameter[float64] `json:"nodeFee"`
	NodeFeeByDemand rocketpool.Parameter[float64] `json:"nodeFeeByDemand"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NetworkBalances contract binding
func NewNetworkFees(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*NetworkFees, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketNetworkFees", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting network fees contract: %w", err)
	}

	return &NetworkFees{
		Details:  NetworkFeesDetails{},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the current network node demand in ETH
func (c *NetworkFees) GetNodeDemand(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.NodeDemand, "getNodeDemand")
}

// Get the current network node commission rate
func (c *NetworkFees) GetNodeFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.NodeFee.RawValue, "getNodeFee")
}

// Get the network node fee for a node demand value
func (c *NetworkFees) GetNodeFeeByDemand(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.NodeFeeByDemand.RawValue, "getNodeFeeByDemand")
}

// Get all basic details
func (c *NetworkFees) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetNodeDemand(mc)
	c.GetNodeFee(mc)
	c.GetNodeFeeByDemand(mc)
}
