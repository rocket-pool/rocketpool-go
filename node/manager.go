package node

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// Settings
const (
	nodeAddressBatchSize               = 50
	nodeDetailsBatchSize               = 20
	SmoothingPoolCountBatchSize uint64 = 2000
	NativeNodeDetailsBatchSize         = 10000
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNodeManager
type NodeManager struct {
	Details  NodeManagerDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for RocketNodeManager
type NodeManagerDetails struct {
	Version   uint8                  `json:"version"`
	NodeCount core.Parameter[uint64] `json:"nodeCount"`
}

// Count of nodes belonging to a timezone
type TimezoneCount struct {
	Timezone string   `abi:"timezone"`
	Count    *big.Int `abi:"count"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NodeManager contract binding
func NewNodeManager(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*NodeManager, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketNodeManager", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting node manager contract: %w", err)
	}

	return &NodeManager{
		Details:  NodeManagerDetails{},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the version of the Node Manager contract
func (c *NodeManager) GetVersion(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.Version, "version")
}

// Get the number of nodes in the network
func (c *NodeManager) GetNodeCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.NodeCount.RawValue, "getNodeCount")
}

// Get all basic details
func (c *NodeManager) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetVersion(mc)
	c.GetNodeCount(mc)
}

// Get the list of node addresses
func (c *NodeManager) GetNodeAddresses(nodeCount uint64, opts *bind.CallOpts) ([]*common.Address, error) {
	// Run the multicall query for each address
	addresses, err := rocketpool.BatchQuery[common.Address](c.rp,
		nodeCount,
		nodeAddressBatchSize,
		func(mc *multicall.MultiCaller, index uint64) (*common.Address, error) {
			address := new(common.Address)
			multicall.AddCall(mc, c.contract, address, "getNodeAt", big.NewInt(int64(index)))
			return address, nil
		},
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting node addresses: %w", err)
	}

	// Return
	return addresses, nil
}

// ===================
// === Sub-Getters ===
// ===================

// Get a node's details by its index and corresponding address
func (c *NodeManager) GetNodeAt(index uint64, address common.Address, staking *NodeStaking, opts *bind.CallOpts) (*Node, error) {
	// Create the node and get details via a multicall query
	node := NewNode(c, staking, index, address)
	err := c.rp.Query(func(mc *multicall.MultiCaller) {
		node.GetAllDetails(mc)
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting node %d (%s): %w", index, address.Hex(), err)
	}

	// Return
	return node, nil
}

// Get the details for all nodes
func (c *NodeManager) GetAllNodes(addresses []*common.Address, staking *NodeStaking, opts *bind.CallOpts) ([]*Node, error) {
	// Run the multicall query for each lot
	nodes, err := rocketpool.BatchQuery[Node](c.rp,
		uint64(len(addresses)),
		nodeDetailsBatchSize,
		func(mc *multicall.MultiCaller, index uint64) (*Node, error) {
			node := NewNode(c, staking, index, *addresses[index])
			node.GetAllDetails(mc)
			return node, nil
		},
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting all node details: %w", err)
	}

	// Return
	return nodes, nil
}

// =============
// === Utils ===
// =============

// Get a breakdown of the number of nodes per timezone for the subset of nodes provided
// NOTE: you will have to call this several times, iterating through subset ranges,
// to get the complete result since a single call for the complete node set may run out of gas
func (c *NodeManager) GetNodeCountPerTimezone(offset *big.Int, limit *big.Int, opts *bind.CallOpts) ([]TimezoneCount, error) {
	return core.CallUntyped[[]TimezoneCount](c.contract, opts, "getNodeCountPerTimezone", offset, limit)
}

// Get the number of nodes in the Smoothing Pool for the subset of nodes provided
// NOTE: you will have to call this several times, iterating through subset ranges,
// to get the complete result since a single call for the complete node set may run out of gas
func (c *NodeManager) GetSmoothingPoolRegisteredNodeCount(offset *big.Int, limit *big.Int, opts *bind.CallOpts) (uint64, error) {
	count, err := core.Call[*big.Int](c.contract, opts, "getSmoothingPoolRegisteredNodeCount", offset, limit)
	if err != nil {
		return 0, err
	}
	return count.Uint64(), nil
}
