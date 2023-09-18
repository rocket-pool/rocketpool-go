package network

import (
	"fmt"
	"math/big"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNetworkFees
type NetworkFees struct {
	*NetworkFeesDetails
	rp *rocketpool.RocketPool
	nf *core.Contract
}

// Details for network fees
type NetworkFeesDetails struct {
	NodeDemand      *big.Int                `json:"nodeDemand"`
	NodeFee         core.Parameter[float64] `json:"nodeFee"`
	NodeFeeByDemand core.Parameter[float64] `json:"nodeFeeByDemand"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NetworkBalances contract binding
func NewNetworkFees(rp *rocketpool.RocketPool) (*NetworkFees, error) {
	// Create the contract
	nf, err := rp.GetContract(rocketpool.ContractName_RocketNetworkFees)
	if err != nil {
		return nil, fmt.Errorf("error getting network fees contract: %w", err)
	}

	return &NetworkFees{
		NetworkFeesDetails: &NetworkFeesDetails{},
		rp:                 rp,
		nf:                 nf,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the current network node demand in ETH
func (c *NetworkFees) GetNodeDemand(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nf, &c.NodeDemand, "getNodeDemand")
}

// Get the current network node commission rate
func (c *NetworkFees) GetNodeFee(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nf, &c.NodeFee.RawValue, "getNodeFee")
}

// Get the network node fee for a node demand value
func (c *NetworkFees) GetNodeFeeByDemand(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nf, &c.NodeFeeByDemand.RawValue, "getNodeFeeByDemand")
}

// Get all basic details
func (c *NetworkFees) GetAllDetails(mc *batch.MultiCaller) {
	c.GetNodeDemand(mc)
	c.GetNodeFee(mc)
	c.GetNodeFeeByDemand(mc)
}
