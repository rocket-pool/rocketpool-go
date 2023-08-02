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
	nodeTimezoneBatchSize       int = 1000
	smoothingPoolCountBatchSize int = 1000
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
func NewNodeManager(rp *rocketpool.RocketPool) (*NodeManager, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketNodeManager)
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

// =================
// === Addresses ===
// =================

// Get a node address by index
func (c *NodeManager) GetNodeAddress(mc *multicall.MultiCaller, address_Out *common.Address, index uint64) {
	multicall.AddCall(mc, c.contract, address_Out, "getNodeAt", big.NewInt(int64(index)))
}

// Get all minipool addresses in a standalone call.
// This will use an internal batched multicall invocation to retrieve all of them.
// Provide the value returned from GetNodeCount() in nodeCount.
func (c *NodeManager) GetNodeAddresses(nodeCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, nodeCount)

	// Run the multicall query for each address
	err := c.rp.BatchQuery(int(nodeCount), c.rp.AddressBatchSize,
		func(mc *multicall.MultiCaller, index int) error {
			c.GetNodeAddress(mc, &addresses[index], uint64(index))
			return nil
		}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting node addresses: %w", err)
	}

	// Return
	return addresses, nil
}

// =============
// === Utils ===
// =============

// Get a breakdown of the number of nodes per timezone for the subset of nodes provided.
// Provide the value returned from GetNodeCount() in nodeCount.
func (c *NodeManager) GetNodeCountPerTimezone(nodeCount uint64, opts *bind.CallOpts) (map[string]uint64, error) {
	timezoneCountMap := map[string]uint64{}

	limit := big.NewInt(int64(nodeTimezoneBatchSize))
	for i := 0; i < int(nodeCount); i += nodeTimezoneBatchSize {
		// Get a batch of timezone counts
		offset := big.NewInt(int64(i))
		timezoneCounts := new([]TimezoneCount)
		if err := c.contract.Call(opts, timezoneCounts, "getNodeCountPerTimezone", offset, limit); err != nil {
			return nil, fmt.Errorf("error getting node counts per timezone (offset %d, limit %d): %w", offset.Uint64(), limit.Uint64(), err)
		}
		for _, countWrapper := range *timezoneCounts {
			timezoneCountMap[countWrapper.Timezone] += countWrapper.Count.Uint64()
		}
	}

	return timezoneCountMap, nil
}

// Get the number of nodes in the Smoothing Pool for the subset of nodes provided.
// Provide the value returned from GetNodeCount() in nodeCount.
func (c *NodeManager) GetSmoothingPoolRegisteredNodeCount(nodeCount uint64, opts *bind.CallOpts) (uint64, error) {
	total := uint64(0)

	limit := big.NewInt(int64(smoothingPoolCountBatchSize))
	for i := 0; i < int(nodeCount); i += smoothingPoolCountBatchSize {
		// Get an SP count from the batch
		offset := big.NewInt(int64(i))
		count := new(*big.Int)
		if err := c.contract.Call(opts, count, "getSmoothingPoolRegisteredNodeCount", offset, limit); err != nil {
			return 0, fmt.Errorf("error getting smoothing pool registration count (offset %d, limit %d): %w", offset.Uint64(), limit.Uint64(), err)
		}
		total += (*count).Uint64()
	}

	return total, nil
}
