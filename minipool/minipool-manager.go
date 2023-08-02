package minipool

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
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
	Details  MinipoolManagerDetails
	rp       *rocketpool.RocketPool
	Contract *core.Contract
}

// Details for RocketMinipoolManager
type MinipoolManagerDetails struct {
	MinipoolCount          core.Parameter[uint64] `json:"minipoolCount"`
	StakingMinipoolCount   core.Parameter[uint64] `json:"stakingMinipoolCount"`
	FinalisedMinipoolCount core.Parameter[uint64] `json:"finalisedMinipoolCount"`
	ActiveMinipoolCount    core.Parameter[uint64] `json:"activeMinipoolCount"`
	VacantMinipoolCount    core.Parameter[uint64] `json:"vacantMinipoolCount"`
}

// The counts of minipools per status
type MinipoolCountsPerStatus struct {
	Initialized  *big.Int `abi:"initialisedCount"`
	Prelaunch    *big.Int `abi:"prelaunchCount"`
	Staking      *big.Int `abi:"stakingCount"`
	Withdrawable *big.Int `abi:"withdrawableCount"`
	Dissolved    *big.Int `abi:"dissolvedCount"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new MinipoolManager contract binding
func NewMinipoolManager(rp *rocketpool.RocketPool) (*MinipoolManager, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolManager)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool manager contract: %w", err)
	}

	return &MinipoolManager{
		Details:  MinipoolManagerDetails{},
		rp:       rp,
		Contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the minipool count
func (c *MinipoolManager) GetMinipoolCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.MinipoolCount.RawValue, "getMinipoolCount")
}

// Get the number of staking minipools in the network
func (c *MinipoolManager) GetStakingMinipoolCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.StakingMinipoolCount.RawValue, "getStakingMinipoolCount")
}

// Get the number of finalised minipools in the network
func (c *MinipoolManager) GetFinalisedMinipoolCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.FinalisedMinipoolCount.RawValue, "getFinalisedMinipoolCount")
}

// Get the number of active minipools in the network
func (c *MinipoolManager) GetActiveMinipoolCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.ActiveMinipoolCount.RawValue, "getActiveMinipoolCount")
}

// Get the number of vacant minipools in the network
func (c *MinipoolManager) GetVacantMinipoolCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.VacantMinipoolCount.RawValue, "getVacantMinipoolCount")
}

// Get all basic details
func (c *MinipoolManager) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetMinipoolCount(mc)
	c.GetStakingMinipoolCount(mc)
	c.GetFinalisedMinipoolCount(mc)
	c.GetActiveMinipoolCount(mc)
}

// =================
// === Addresses ===
// =================

// Get a minipool address by index
func (c *MinipoolManager) GetMinipoolAddress(mc *multicall.MultiCaller, address_Out *common.Address, index uint64) {
	multicall.AddCall(mc, c.Contract, address_Out, "getMinipoolAt", big.NewInt(int64(index)))
}

// Get a minipool address by pubkey
func (c *MinipoolManager) GetMinipoolAddressByPubkey(mc *multicall.MultiCaller, address_Out *common.Address, pubkey types.ValidatorPubkey) {
	multicall.AddCall(mc, c.Contract, address_Out, "getMinipoolByPubkey", pubkey[:])
}

// Get a vacant minipool address by index
func (c *MinipoolManager) GetVacantMinipoolAddress(mc *multicall.MultiCaller, address_Out *common.Address, index uint64) {
	multicall.AddCall(mc, c.Contract, address_Out, "getVacantMinipoolAt", big.NewInt(int64(index)))
}

// Get all minipool addresses in a standalone call.
// This will use an internal batched multicall invocation to retrieve all of them.
// Provide the value returned from GetMinipoolCount() in minipoolCount.
func (c *MinipoolManager) GetMinipoolAddresses(mc *multicall.MultiCaller, minipoolCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, minipoolCount)

	// Run the multicall query for each address
	err := c.rp.BatchQuery(int(minipoolCount), c.rp.AddressBatchSize,
		func(mc *multicall.MultiCaller, index int) error {
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
		if err := c.Contract.Call(opts, newAddresses, "getPrelaunchMinipools", offset, limit); err != nil {
			return []common.Address{}, fmt.Errorf("error getting prelaunch minipool addresses (offset %d, limit %d): %w", offset.Uint64(), limit.Uint64(), err)
		}
		addresses = append(addresses, *newAddresses...)
	}

	return addresses, nil
}

// Get all minipool addresses in a standalone call.
// This will use an internal batched multicall invocation to retrieve all of them.
// Provide the value returned from GetVacantMinipoolCount() in minipoolCount.
func (c *MinipoolManager) GetVacantMinipoolAddresses(mc *multicall.MultiCaller, minipoolCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, minipoolCount)

	// Run the multicall query for each address
	err := c.rp.BatchQuery(int(minipoolCount), c.rp.AddressBatchSize,
		func(mc *multicall.MultiCaller, index int) error {
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
		if err := c.Contract.Call(opts, newMinipoolCounts, "getMinipoolCountPerStatus", offset, limit); err != nil {
			return MinipoolCountsPerStatus{}, fmt.Errorf("Could not get minipool counts: %w", err)
		}
		if newMinipoolCounts != nil {
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
	}
	return minipoolCounts, nil
}
