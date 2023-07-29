package minipool

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketMinipoolBondReducer
type MinipoolBondReducer struct {
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for RocketMinipoolBondReducer for a specific minipool
type MinipoolBondReducerDetails struct {
	Address                      common.Address            `json:"address"`
	IsBondReduceCancelled        bool                      `json:"isBondReduceCancelled"`
	ReduceBondTime               core.Parameter[time.Time] `json:"reduceBondTime"`
	ReduceBondValue              *big.Int                  `json:"reduceBondValue"`
	LastBondReductionTime        core.Parameter[time.Time] `json:"lastBondReductionTime"`
	LastBondReductionPrevValue   *big.Int                  `json:"lastBondReductionPrevValue"`
	LastBondReductionPrevNodeFee core.Parameter[float64]   `json:"lastBondReductionPrevNodeFee"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new MinipoolBondReducer contract binding
func NewMinipoolBondReducer(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*MinipoolBondReducer, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketMinipoolBondReducer", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool bond reducer contract: %w", err)
	}

	return &MinipoolBondReducer{
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Gets whether or not the bond reduction process for the minipool has already been cancelled
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolBondReducer) GetReduceBondCancelled(mc *multicall.MultiCaller, details *MinipoolBondReducerDetails) {
	multicall.AddCall(mc, c.contract, &details.IsBondReduceCancelled, "getReduceBondCancelled", details.Address)
}

// Gets the time at which the MP owner started the bond reduction process
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolBondReducer) GetReduceBondTime(mc *multicall.MultiCaller, details *MinipoolBondReducerDetails) {
	multicall.AddCall(mc, c.contract, &details.ReduceBondTime.RawValue, "getReduceBondTime", details.Address)
}

// Gets the amount of ETH a minipool is reducing its bond to
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolBondReducer) GetReduceBondValue(mc *multicall.MultiCaller, details *MinipoolBondReducerDetails) {
	multicall.AddCall(mc, c.contract, &details.ReduceBondValue, "getReduceBondValue", details.Address)
}

// Gets the timestamp at which the bond was last reduced
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolBondReducer) GetLastBondReductionTime(mc *multicall.MultiCaller, details *MinipoolBondReducerDetails) {
	multicall.AddCall(mc, c.contract, &details.LastBondReductionTime.RawValue, "getLastBondReductionTime", details.Address)
}

// Gets the previous bond amount of the minipool prior to its last reduction
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolBondReducer) GetLastBondReductionPrevValue(mc *multicall.MultiCaller, details *MinipoolBondReducerDetails) {
	multicall.AddCall(mc, c.contract, &details.LastBondReductionPrevValue, "getLastBondReductionPrevValue", details.Address)
}

// Gets the previous node fee (commission) of the minipool prior to its last reduction
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolBondReducer) GetLastBondReductionPrevNodeFee(mc *multicall.MultiCaller, details *MinipoolBondReducerDetails) {
	multicall.AddCall(mc, c.contract, &details.LastBondReductionPrevNodeFee.RawValue, "getLastBondReductionPrevNodeFee", details.Address)
}

// Get all basic details
func (c *MinipoolBondReducer) GetAllDetails(mc *multicall.MultiCaller, details *MinipoolBondReducerDetails) {
	c.GetReduceBondCancelled(mc, details)
	c.GetReduceBondTime(mc, details)
	c.GetReduceBondValue(mc, details)
	c.GetLastBondReductionTime(mc, details)
	c.GetLastBondReductionPrevValue(mc, details)
	c.GetLastBondReductionPrevNodeFee(mc, details)
}

// ====================
// === Transactions ===
// ====================

// Get info for beginning a minipool bond reduction
func (c *MinipoolBondReducer) BeginReduceBondAmount(minipoolAddress common.Address, newBondAmount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "beginReduceBondAmount", opts, minipoolAddress, newBondAmount)
}

// Get info for voting to cancel a minipool's bond reduction
func (c *MinipoolBondReducer) VoteCancelReduction(minipoolAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "voteCancelReduction", opts, minipoolAddress)
}
