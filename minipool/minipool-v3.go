package minipool

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

const (
	minipoolV3EncodedAbi string = "eJztWltv2koQ/isVz3lq1SrqW3Jon5qeCJKch6pCY+8Aqyy71l7MQdH57x0bY2Mw2MAau0d9SsDjmW8uO7flx9sApJKrhXJm8HkKwuDNgMvIWfr4443+Zfgvsq1HFrUE8bSKcPB54Ojz+4+fBjcDCYvki0hjzInXvZLsjphKS89smfi/m9P5Slz6Zmk5/dnn9DMnSASOkLmQmOZ0GCMBSOQ1s5vVbh8LMKbRmALLVKtFIWPz+Byt4JoW+mLnqIcYKcNt+0Yi2tCV5FxiKKkY+gwnZ1B7DU9lQdyDABlWOaE1d/7D7ZxpWIJ41Cok67bvWKt+19jfGEueb6JqPMHK4paFYhCcgVX60QWvuCqkrekaKHiI4ZjPJFin8VyeH94XXNk6FQzBwkgpu8OSKK/s1R2ll3lo/6WRkZs4vXWG3qcHywOXPFKKjhQaC6+XHCn/kNRCXZS9fSMah9oFwRUQBZ47itBpTUCvm7Q3VnuBkKSuKMIi0O3n7AUuAtR+8natjmlEvFwWpZVqJbJuCyB0Nq0zZSi3LamUinqOKAlXa5XrULwT4IzLZ+ozhtxYzQNqhuhN5WxBmWiAD85CwAW3q7TPkRGsIBBbeKZOhpYrWRb0VqvVJDgY2gXKEEToBMH4Th3WeA66DLJeSjX3fc1ijstr65RYv2udtmHJLH0fhxNQfqjCkn7vCcg4qWtdwxDKeD0SP3dbm3rmp3PmxigRt3uUy6afaFyCZuZvKVYVbiiAZVmmqGgtWHbKCSZvy28ztNmAmlriaHyidIt3m3o6zB1eeYRvvQQsofuaqc86PjsEJUnYd3mhPwJmrxuo6AP8Qco8cWcMzSp9sVMGqupkdFUNMmRfsVdwRkiErId2elLRc/QCwvUGFfXu34DI5n1D9cBnGhKy7t1Yn4smB0eTklqlmaLjhDLezB1Ni1P+Qru1aS3mXqjwtS/RuIb0tB6veoHoKVnU5tFU48XrwUrnxL5VzgpQfXLlFrzuM90+srwf74En03WT7QRHgxIgS53s4TpA6TRZAnuePposHZRkmypfuSU8tseeNLoOaMahZv9/fN8/qV/4F7aONA49D9ENDX18x1vDoHapW9Iw2b1mm9h25tmoYt/jjblOb7tLF+2tSEkmgnZ4GwFm3u4mpfuTZPZWbd7s54oK2MKqy5Wq2J0Qatl5LYtRm+RxbaG/rYLhr92O6VSnXaRnxzZ7PVSS3OJCiuOL79n2yrESbIgCZ2C3JJ5wYbTHUOLyIoanX9dsxI2UEMjuIR2IMnI/t5P/Z6s9RzMN7M/PqA6YKf0pyQhD5HHju0C2FZAB7MznHm89Sg5sR0i6os/jshdr5y/TKRJNjH0D9pj9vrJvuGg+/UZMTENPdjcZahW+JlskpWF2MA3+XtPlzt2eQWu5nFUYddM+1rnrPIBN+kPIklVDer3OiQT+F40AQik="
)

// The decoded ABI for v2 minipools
var minipoolV3Abi *abi.ABI

// ===============
// === Structs ===
// ===============

type MinipoolV3 struct {
	*minipoolCommon
	Details MinipoolV3Details
	br      *core.Contract
}

type MinipoolV3Details struct {
	*minipoolCommonDetails
	IsVacant                     bool                      `json:"isVacant"`
	PreMigrationBalance          *big.Int                  `json:"preMigrationBalance"`
	HasUserDistributed           bool                      `json:"hasUserDistributed"`
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

// Create new minipool contract
func newMinipool_v3(rp *rocketpool.RocketPool, address common.Address) (*MinipoolV3, error) {
	var contract *core.Contract
	var err error
	if minipoolV3Abi == nil {
		// Get contract
		contract, err = rp.CreateMinipoolContractFromEncodedAbi(address, minipoolV3EncodedAbi)
	} else {
		contract, err = rp.CreateMinipoolContractFromAbi(address, minipoolV3Abi)
	}
	if err != nil {
		return nil, err
	} else if minipoolV3Abi == nil {
		minipoolV3Abi = contract.ABI
	}

	// Get the BondReducer binding
	br, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolBondReducer)
	if err != nil {
		return nil, fmt.Errorf("error creating minipool bond reducer: %w", err)
	}

	// Create the base binding
	base, err := newMinipoolCommonFromVersion(rp, contract, 3)
	if err != nil {
		return nil, fmt.Errorf("error creating minipool base: %w", err)
	}

	// Create and return
	return &MinipoolV3{
		minipoolCommon: base,
		Details: MinipoolV3Details{
			minipoolCommonDetails: &base.Details,
		},
		br: br,
	}, nil
}

// Get the minipool as a v3 minipool if it implements the required methods
func GetMinipoolAsV3(mp IMinipool) (*MinipoolV3, bool) {
	castedMp, ok := mp.(*MinipoolV3)
	if ok {
		return castedMp, true
	}
	return nil, false
}

// =============
// === Calls ===
// =============

// Get the basic details
func (c *MinipoolV3) QueryAllDetails(mc *batch.MultiCaller) {
	c.minipoolCommon.QueryAllDetails(mc)
	c.GetVacant(mc)
	c.GetPreMigrationBalance(mc)
	c.GetUserDistributed(mc)
	c.GetReduceBondCancelled(mc)
	c.GetReduceBondTime(mc)
	c.GetReduceBondValue(mc)
	c.GetLastBondReductionTime(mc)
	c.GetLastBondReductionPrevValue(mc)
	c.GetLastBondReductionPrevNodeFee(mc)
}

// === Minipool ===

// Check if this is a vacant minipool (pre-staking solo migration)
func (c *MinipoolV3) GetVacant(mc *batch.MultiCaller) {
	core.AddCall(mc, c.Contract, &c.Details.IsVacant, "getVacant")
}

// Get the node deposit balance of this minipool before its last bond reduction
func (c *MinipoolV3) GetPreMigrationBalance(mc *batch.MultiCaller) {
	core.AddCall(mc, c.Contract, &c.Details.PreMigrationBalance, "getPreMigrationBalance")
}

// Check if the minipool's balance has already been distributed by someone other than the node operator
func (c *MinipoolV3) GetUserDistributed(mc *batch.MultiCaller) {
	core.AddCall(mc, c.Contract, &c.Details.PreMigrationBalance, "getUserDistributed")
}

// === BondReducer ===

// Gets whether or not the bond reduction process for the minipool has already been cancelled
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolV3) GetReduceBondCancelled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.br, &c.Details.IsBondReduceCancelled, "getReduceBondCancelled", c.minipoolCommon.Details.Address)
}

// Gets the time at which the MP owner started the bond reduction process
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolV3) GetReduceBondTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.br, &c.Details.ReduceBondTime.RawValue, "getReduceBondTime", c.minipoolCommon.Details.Address)
}

// Gets the amount of ETH a minipool is reducing its bond to
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolV3) GetReduceBondValue(mc *batch.MultiCaller) {
	core.AddCall(mc, c.br, &c.Details.ReduceBondValue, "getReduceBondValue", c.minipoolCommon.Details.Address)
}

// Gets the timestamp at which the bond was last reduced
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolV3) GetLastBondReductionTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.br, &c.Details.LastBondReductionTime.RawValue, "getLastBondReductionTime", c.minipoolCommon.Details.Address)
}

// Gets the previous bond amount of the minipool prior to its last reduction
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolV3) GetLastBondReductionPrevValue(mc *batch.MultiCaller) {
	core.AddCall(mc, c.br, &c.Details.LastBondReductionPrevValue, "getLastBondReductionPrevValue", c.minipoolCommon.Details.Address)
}

// Gets the previous node fee (commission) of the minipool prior to its last reduction
// The output will be stored in details - note that the Address must already be set!
func (c *MinipoolV3) GetLastBondReductionPrevNodeFee(mc *batch.MultiCaller) {
	core.AddCall(mc, c.br, &c.Details.LastBondReductionPrevNodeFee.RawValue, "getLastBondReductionPrevNodeFee", c.minipoolCommon.Details.Address)
}

// ====================
// === Transactions ===
// ====================

// === Minipool ===

// Get info for distributing the minipool's ETH balance to the node operator and rETH staking pool
func (c *MinipoolV3) DistributeBalance(opts *bind.TransactOpts, rewardsOnly bool) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "distributeBalance", opts, rewardsOnly)
}

// Get info for reducing a minipool's bond
func (c *MinipoolV3) ReduceBondAmount(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "reduceBondAmount", opts)
}

// Get info for promoting a vacant minipool
func (c *MinipoolV3) Promote(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "promote", opts)
}

// === BondReducer ===

// Get info for beginning a minipool bond reduction
func (c *MinipoolV3) BeginReduceBondAmount(newBondAmount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.br, "beginReduceBondAmount", opts, c.minipoolCommon.Details.Address, newBondAmount)
}

// Get info for voting to cancel a minipool's bond reduction
func (c *MinipoolV3) VoteCancelReduction(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.br, "voteCancelReduction", opts, c.minipoolCommon.Details.Address)
}
