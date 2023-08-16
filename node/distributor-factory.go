package node

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNodeDistributorFactory
type NodeDistributorFactory struct {
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new NodeDistributorFactory contract binding
func NewNodeDistributorFactory(rp *rocketpool.RocketPool) (*NodeDistributorFactory, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketNodeDistributorFactory)
	if err != nil {
		return nil, fmt.Errorf("error getting node distributor factory contract: %w", err)
	}

	return &NodeDistributorFactory{
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Gets the deterministic address for a node's reward distributor contract
func (c *NodeDistributorFactory) GetDistributorAddress(mc *batch.MultiCaller, nodeAddress common.Address, address_Out *common.Address) {
	core.AddCall(mc, c.contract, address_Out, "getProxyAddress", nodeAddress)
}

// ===================
// === Sub-Getters ===
// ===================

// Get a node's distributor with details
func (c *NodeDistributorFactory) GetNodeDistributor(nodeAddress common.Address, distributorAddress common.Address, opts *bind.CallOpts) (*NodeDistributor, error) {
	// Create the distributor
	distributor, err := NewNodeDistributor(c.rp, nodeAddress, distributorAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating node distributor binding for node %s at %s: %w", nodeAddress.Hex(), distributorAddress.Hex(), err)
	}

	// Get details via a multicall query
	err = c.rp.Query(func(mc *batch.MultiCaller) error {
		distributor.GetAllDetails(mc)
		return nil
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting node distributor for node %s at %s: %w", nodeAddress.Hex(), distributorAddress.Hex(), err)
	}

	// Return
	return distributor, nil
}
