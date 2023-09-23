package minipool

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
)

// Settings
const (
	minipoolBatchSize          int = 100
	minipoolPrelaunchBatchSize int = 250
)

// ===============
// === Structs ===
// ===============

// Binding for RocketMinipoolManager
type MinipoolManager struct {
	*MinipoolManagerDetails
	rp    *rocketpool.RocketPool
	mpMgr *core.Contract
	mq    *core.Contract
}

// Details for RocketMinipoolManager
type MinipoolManagerDetails struct {
	MinipoolCount          core.Uint256Parameter[uint64] `json:"minipoolCount"`
	StakingMinipoolCount   core.Uint256Parameter[uint64] `json:"stakingMinipoolCount"`
	FinalisedMinipoolCount core.Uint256Parameter[uint64] `json:"finalisedMinipoolCount"`
	ActiveMinipoolCount    core.Uint256Parameter[uint64] `json:"activeMinipoolCount"`
	VacantMinipoolCount    core.Uint256Parameter[uint64] `json:"vacantMinipoolCount"`

	TotalQueueLength       core.Uint256Parameter[uint64] `json:"totalQueueLength"`
	TotalQueueCapacity     *big.Int                      `json:"totalQueueCapacity"`
	EffectiveQueueCapacity *big.Int                      `json:"effectiveQueueCapacity"`
}

// The counts of minipools per status
type MinipoolCountsPerStatus struct {
	Initialized  *big.Int `abi:"initialisedCount"`
	Prelaunch    *big.Int `abi:"prelaunchCount"`
	Staking      *big.Int `abi:"stakingCount"`
	Withdrawable *big.Int `abi:"withdrawableCount"`
	Dissolved    *big.Int `abi:"dissolvedCount"`
}

// Minipools queue status details
type QueueDetails struct {
	Position int64
}

// ====================
// === Constructors ===
// ====================

// Creates a new MinipoolManager contract binding
func NewMinipoolManager(rp *rocketpool.RocketPool) (*MinipoolManager, error) {
	// Create the contracts
	mpMgr, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolManager)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool manager contract: %w", err)
	}
	mq, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolQueue)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool queue contract: %w", err)
	}

	return &MinipoolManager{
		MinipoolManagerDetails: &MinipoolManagerDetails{},
		rp:                     rp,
		mpMgr:                  mpMgr,
		mq:                     mq,
	}, nil
}

// =============
// === Calls ===
// =============

// === MinipoolManager ===

// Get the minipool count
func (c *MinipoolManager) GetMinipoolCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.MinipoolCount.RawValue, "getMinipoolCount")
}

// Get the number of staking minipools in the network
func (c *MinipoolManager) GetStakingMinipoolCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.StakingMinipoolCount.RawValue, "getStakingMinipoolCount")
}

// Get the number of finalised minipools in the network
func (c *MinipoolManager) GetFinalisedMinipoolCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.FinalisedMinipoolCount.RawValue, "getFinalisedMinipoolCount")
}

// Get the number of active minipools in the network
func (c *MinipoolManager) GetActiveMinipoolCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.ActiveMinipoolCount.RawValue, "getActiveMinipoolCount")
}

// Get the number of vacant minipools in the network
func (c *MinipoolManager) GetVacantMinipoolCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.VacantMinipoolCount.RawValue, "getVacantMinipoolCount")
}

// === MinipoolQueue ===

// Get the total length of the minipool queue
func (c *MinipoolManager) GetTotalQueueLength(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mq, &c.TotalQueueLength.RawValue, "getTotalLength")
}

// Get the total capacity of the minipool queue
func (c *MinipoolManager) GetTotalQueueCapacity(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mq, &c.TotalQueueCapacity, "getTotalCapacity")
}

// Get the total effective capacity of the minipool queue (used in node demand calculation)
func (c *MinipoolManager) GetEffectiveQueueCapacity(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mq, &c.EffectiveQueueCapacity, "getEffectiveCapacity")
}

// =================
// === Addresses ===
// =================

// Get a minipool address by index
func (c *MinipoolManager) GetMinipoolAddress(mc *batch.MultiCaller, address_Out *common.Address, index uint64) {
	core.AddCall(mc, c.mpMgr, address_Out, "getMinipoolAt", big.NewInt(int64(index)))
}

// Get a minipool address by pubkey
func (c *MinipoolManager) GetMinipoolAddressByPubkey(mc *batch.MultiCaller, address_Out *common.Address, pubkey types.ValidatorPubkey) {
	core.AddCall(mc, c.mpMgr, address_Out, "getMinipoolByPubkey", pubkey[:])
}

// Get a vacant minipool address by index
func (c *MinipoolManager) GetVacantMinipoolAddress(mc *batch.MultiCaller, address_Out *common.Address, index uint64) {
	core.AddCall(mc, c.mpMgr, address_Out, "getVacantMinipoolAt", big.NewInt(int64(index)))
}

// Get all minipool addresses in a standalone call.
// This will use an internal batched multicall invocation to retrieve all of them.
// Provide the value returned from GetMinipoolCount() in minipoolCount.
func (c *MinipoolManager) GetMinipoolAddresses(mc *batch.MultiCaller, minipoolCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, minipoolCount)

	// Run the multicall query for each address
	err := c.rp.BatchQuery(int(minipoolCount), c.rp.AddressBatchSize,
		func(mc *batch.MultiCaller, index int) error {
			c.GetMinipoolAddress(mc, &addresses[index], uint64(index))
			return nil
		}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool addresses: %w", err)
	}

	// Return
	return addresses, nil
}

// Get all prelaunch minipool addresses in a standalone call.
// This will use an internal batched multicall invocation to retrieve all of them.
// Provide the value returned from GetMinipoolCount() in minipoolCount.
func (c *MinipoolManager) GetPrelaunchMinipoolAddresses(minipoolCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, 0, minipoolCount)

	limit := big.NewInt(int64(minipoolPrelaunchBatchSize))
	for i := 0; i < int(minipoolCount); i += minipoolPrelaunchBatchSize {
		// Get a batch of addresses
		offset := big.NewInt(int64(i))
		newAddresses := new([]common.Address)
		if err := c.mpMgr.Call(opts, newAddresses, "getPrelaunchMinipools", offset, limit); err != nil {
			return []common.Address{}, fmt.Errorf("error getting prelaunch minipool addresses (offset %d, limit %d): %w", offset.Uint64(), limit.Uint64(), err)
		}
		addresses = append(addresses, *newAddresses...)
	}

	return addresses, nil
}

// Get all minipool addresses in a standalone call.
// This will use an internal batched multicall invocation to retrieve all of them.
// Provide the value returned from GetVacantMinipoolCount() in minipoolCount.
func (c *MinipoolManager) GetVacantMinipoolAddresses(mc *batch.MultiCaller, minipoolCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, minipoolCount)

	// Run the multicall query for each address
	err := c.rp.BatchQuery(int(minipoolCount), c.rp.AddressBatchSize,
		func(mc *batch.MultiCaller, index int) error {
			c.GetVacantMinipoolAddress(mc, &addresses[index], uint64(index))
			return nil
		}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting vacant minipool addresses: %w", err)
	}

	// Return
	return addresses, nil
}

// =============
// === Utils ===
// =============

// Get the minipool count by status
func (c *MinipoolManager) GetMinipoolCountPerStatus(minipoolCount uint64, opts *bind.CallOpts) (MinipoolCountsPerStatus, error) {
	minipoolCounts := MinipoolCountsPerStatus{
		Initialized:  big.NewInt(0),
		Prelaunch:    big.NewInt(0),
		Staking:      big.NewInt(0),
		Dissolved:    big.NewInt(0),
		Withdrawable: big.NewInt(0),
	}

	limit := big.NewInt(int64(minipoolPrelaunchBatchSize))
	for i := 0; i < int(minipoolCount); i += minipoolPrelaunchBatchSize {
		// Get a batch of counts
		offset := big.NewInt(int64(i))
		newMinipoolCounts := new(MinipoolCountsPerStatus)
		if err := c.mpMgr.Call(opts, newMinipoolCounts, "getMinipoolCountPerStatus", offset, limit); err != nil {
			return MinipoolCountsPerStatus{}, fmt.Errorf("error getting minipool counts: %w", err)
		}
		if newMinipoolCounts.Initialized != nil {
			minipoolCounts.Initialized.Add(minipoolCounts.Initialized, newMinipoolCounts.Initialized)
		}
		if newMinipoolCounts.Prelaunch != nil {
			minipoolCounts.Prelaunch.Add(minipoolCounts.Prelaunch, newMinipoolCounts.Prelaunch)
		}
		if newMinipoolCounts.Staking != nil {
			minipoolCounts.Staking.Add(minipoolCounts.Staking, newMinipoolCounts.Staking)
		}
		if newMinipoolCounts.Dissolved != nil {
			minipoolCounts.Dissolved.Add(minipoolCounts.Dissolved, newMinipoolCounts.Dissolved)
		}
		if newMinipoolCounts.Withdrawable != nil {
			minipoolCounts.Withdrawable.Add(minipoolCounts.Withdrawable, newMinipoolCounts.Withdrawable)
		}
	}
	return minipoolCounts, nil
}

// Get the 0x01-based withdrawal credentials for a minipool address (even if it doesn't exist yet)
func (c *MinipoolManager) GetMinipoolWithdrawalCredentials(mc *batch.MultiCaller, credentials_Out *common.Hash, address common.Address) {
	core.AddCall(mc, c.mpMgr, credentials_Out, "getMinipoolWithdrawalCredentials", address)
}

// Create a minipool binding from an explicit version number
func (c *MinipoolManager) NewMinipoolFromVersion(address common.Address, version uint8) (IMinipool, error) {
	switch version {
	case 1, 2:
		return newMinipool_v2(c.rp, address)
	case 3:
		return newMinipool_v3(c.rp, address)
	default:
		return nil, fmt.Errorf("unexpected minipool contract version [%d]", version)
	}
}

// Create a minipool binding from its address
func (c *MinipoolManager) CreateMinipoolFromAddress(address common.Address, includeDetails bool, opts *bind.CallOpts) (IMinipool, error) {
	// Get the minipool version
	var version uint8
	results, err := c.rp.FlexQuery(func(mc *batch.MultiCaller) error {
		return rocketpool.GetContractVersion(mc, &version, address)
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error querying minipool version: %w", err)
	}
	if !results[0] {
		// If it failed, this is a contract on Prater from before version() existed so it's v1
		version = 1
	}

	// Get the minipool
	minipool, err := c.NewMinipoolFromVersion(address, version)
	if err != nil {
		return nil, fmt.Errorf("error creating minipool: %w", err)
	}

	// Include the details if requested
	if includeDetails {
		err := c.rp.Query(func(mc *batch.MultiCaller) error {
			minipool.QueryAllDetails(mc)
			return nil
		}, opts)
		if err != nil {
			return nil, fmt.Errorf("error getting minipool %s details: %w", address.Hex(), err)
		}
	}

	return minipool, nil
}

// Create bindings for all minipools from the provided addresses in a standalone call.
// This will use an internal batched multicall invocation to build all of them quickly.
func (c *MinipoolManager) CreateMinipoolsFromAddresses(addresses []common.Address, includeDetails bool, opts *bind.CallOpts) ([]IMinipool, error) {
	minipoolCount := len(addresses)

	// Get the minipool versions
	versions := make([]uint8, minipoolCount)
	err := c.rp.FlexBatchQuery(int(minipoolCount), c.rp.ContractVersionBatchSize,
		func(mc *batch.MultiCaller, index int) error {
			return rocketpool.GetContractVersion(mc, &versions[index], addresses[index])
		},
		func(result bool, index int) error {
			if !result {
				// If it failed, this is a contract on Prater from before version() existed so it's v1
				versions[index] = 1
			}
			return nil
		}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool versions: %w", err)
	}

	// Create the minipools
	minipools := make([]IMinipool, minipoolCount)
	for i := 0; i < int(minipoolCount); i++ {
		minipool, err := c.NewMinipoolFromVersion(addresses[i], versions[i])
		if err != nil {
			return nil, fmt.Errorf("error creating minipool %d (%s): %w", i, addresses[i].Hex(), err)
		}
		minipools[i] = minipool
	}

	// Include the details if requested
	if includeDetails {
		err := c.rp.BatchQuery(int(minipoolCount), minipoolBatchSize, func(mc *batch.MultiCaller, index int) error {
			minipools[index].QueryAllDetails(mc)
			return nil
		}, opts)
		if err != nil {
			return nil, fmt.Errorf("error getting minipool details: %w", err)
		}
	}

	return minipools, nil
}

// =============
// === Utils ===
// =============

// Get the minipool at the specified position in queue (0-indexed).
func (c *MinipoolManager) GetMinipoolAtQueuePosition(mc *batch.MultiCaller, address_Out *common.Address, position uint64) {
	core.AddCall(mc, c.mq, address_Out, "getMinipoolAt", big.NewInt(int64(position)))
}
