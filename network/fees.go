package network

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	// Contract names
	NetworkFees_ContractName string = "rocketNetworkFees"

	// Calls
	networkFees_getNodeDemand      string = "getNodeDemand"
	networkFees_getNodeFee         string = "getNodeFee"
	networkFees_getNodeFeeByDemand string = "getNodeFeeByDemand"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNetworkFees
type NetworkFees struct {
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// Multicall details for network fees
type NetworkFeesDetails struct {
	// Raw parameters
	NodeDemand         *big.Int `json:"nodeDemand"`
	NodeFeeRaw         *big.Int `json:"nodeFeeRaw"`
	NodeFeeByDemandRaw *big.Int `json:"nodeFeeByDemandRaw"`

	// Formatted parameters
	NodeFee         float64 `json:"nodeFee"`
	NodeFeeByDemand float64 `json:"nodeFeeByDemand"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NetworkBalances contract binding
func NewNetworkFees(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*NetworkFees, error) {
	// Create the contract
	contract, err := rp.GetContract(NetworkFees_ContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting network fees contract: %w", err)
	}

	return &NetworkFees{
		rp:       rp,
		contract: contract,
	}, nil
}

// ===================
// === Raw Getters ===
// ===================

// Get the current network node demand in ETH
func (c *NetworkFees) GetNodeDemand(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkFees_getNodeDemand)
}

// Get the current network node commission rate
func (c *NetworkFees) GetNodeFeeRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkFees_getNodeFee)
}

// Get the network node fee for a node demand value
func (c *NetworkFees) GetNodeFeeByDemandRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkFees_getNodeFeeByDemand)
}

// =========================
// === Formatted Getters ===
// =========================

// Get the current network node commission rate
func (c *NetworkFees) GetNodeFee(opts *bind.CallOpts) (float64, error) {
	raw, err := c.GetNodeFeeRaw(opts)
	if err != nil {
		return 0, err
	}
	return eth.WeiToEth(raw), nil
}

// Get the network node fee for a node demand value
func (c *NetworkFees) GetNodeFeeByDemand(opts *bind.CallOpts) (float64, error) {
	raw, err := c.GetNodeFeeByDemandRaw(opts)
	if err != nil {
		return 0, err
	}
	return eth.WeiToEth(raw), nil
}

// =================
// === Multicall ===
// =================

// Add queries to a multicall batcher
func (c *NetworkFees) AddMulticallQueries(mc *multicall.MultiCaller, details *NetworkFeesDetails) {
	mc.AddCall(c.contract, &details.NodeDemand, networkFees_getNodeDemand)
	mc.AddCall(c.contract, &details.NodeFeeRaw, networkFees_getNodeFee)
	mc.AddCall(c.contract, &details.NodeFeeByDemandRaw, networkFees_getNodeFeeByDemand)
}

// Postprocess the multicalled data to get the formatted parameters
func (c *NetworkFees) PostprocessAfterMulticall(details *NetworkFeesDetails) {
	details.NodeFee = eth.WeiToEth(details.NodeFeeRaw)
	details.NodeFeeByDemand = eth.WeiToEth(details.NodeFeeByDemandRaw)
}
