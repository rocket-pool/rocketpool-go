package minipool

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// Settings
const (
	MinipoolPrelaunchBatchSize     = 250
	MinipoolAddressBatchSize       = 50
	MinipoolDetailsBatchSize       = 20
	NativeMinipoolDetailsBatchSize = 1000
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

// ====================
// === Transactions ===
// ====================

// Get info for creating a new lot
func (c *AuctionManager) CreateLot(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "createLot", opts)
}

// ===================
// === Sub-Getters ===
// ===================

// Get a minipool with details
func (c *MinipoolManager) GetMinipool(index uint64, opts *bind.CallOpts) (Minipool, error) {

	// Decide how to do this - just get the empty binding, the address and version, or what?
	// Or do the original thing - add a raw getter for the address and then convenience methods to wrap it all up so people can do what they want

	// Create the lot and get details via a multicall query
	lot := NewAuctionLot(c, index)
	err := c.rp.Query(func(mc *multicall.MultiCaller) {
		if bidder != nil {
			lot.GetAllDetailsWithBidAmount(mc, *bidder)
		} else {
			lot.GetAllDetails(mc)
		}
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting lot: %w", err)
	}

	// Return
	return lot, nil
}
