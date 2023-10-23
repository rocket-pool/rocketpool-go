package tokens

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketTokenRPL
type TokenRpl struct {
	// The RPL total supply
	TotalSupply *core.SimpleField[*big.Int]

	// The RPL inflation interval rate
	InflationIntervalRate *core.SimpleField[*big.Int]

	// The time that inflation started for the current interval
	InflationIntervalStartTime *core.FormattedUint256Field[time.Time]

	// === Internal fields ===
	rpl *core.Contract
	rp  *rocketpool.RocketPool
}

// ====================
// === Constructors ===
// ====================

// Creates a new TokenRpl contract binding
func NewTokenRpl(rp *rocketpool.RocketPool) (*TokenRpl, error) {
	// Create the contract
	rpl, err := rp.GetContract(rocketpool.ContractName_RocketTokenRPL)
	if err != nil {
		return nil, fmt.Errorf("error getting RPL contract: %w", err)
	}

	return &TokenRpl{
		TotalSupply:                core.NewSimpleField[*big.Int](rpl, "totalSupply"),
		InflationIntervalRate:      core.NewSimpleField[*big.Int](rpl, "getInflationIntervalRate"),
		InflationIntervalStartTime: core.NewFormattedUint256Field[time.Time](rpl, "getInflationIntervalStartTime"),

		rp:  rp,
		rpl: rpl,
	}, nil
}

// =============
// === Calls ===
// =============

// === Core ERC-20 functions ===

// Get the RPL balance of an address
func (c *TokenRpl) BalanceOf(mc *batch.MultiCaller, balance_Out **big.Int, address common.Address) {
	core.AddCall(mc, c.rpl, balance_Out, "balanceOf", address)
}

// Get the RPL spending allowance of an address and spender
func (c *TokenRpl) GetAllowance(mc *batch.MultiCaller, allowance_Out **big.Int, owner common.Address, spender common.Address) {
	core.AddCall(mc, c.rpl, allowance_Out, "allowance", owner, spender)
}

// ====================
// === Transactions ===
// ====================

// === Core ERC-20 functions ===

// Get info for approving RPL's usage by a spender
func (c *TokenRpl) Approve(spender common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.rpl, "approve", opts, spender, amount)
}

// Get info for transferring RPL
func (c *TokenRpl) Transfer(to common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.rpl, "transfer", opts, to, amount)
}

// Get info for transferring RPL from a sender
func (c *TokenRpl) TransferFrom(from common.Address, to common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.rpl, "transferFrom", opts, from, to, amount)
}

// === RPL functions ===

// Get info for minting new RPL tokens from inflation
func (c *TokenRpl) MintInflationRPL(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.rpl, "inflationMintTokens", opts)
}

// Get info for swapping fixed-supply RPL for new RPL tokens
func (c *TokenRpl) SwapFixedSupplyRplForRpl(amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.rpl, "swapTokens", opts, amount)
}
