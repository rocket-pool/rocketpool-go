package tokens

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketTokenRETH
type TokenReth struct {
	*TokenRethDetails
	rp   *rocketpool.RocketPool
	reth *core.Contract
}

// Details for RocketTokenRETH
type TokenRethDetails struct {
	TotalSupply     *big.Int                       `json:"totalSupply"`
	ExchangeRate    core.Uint256Parameter[float64] `json:"exchangeRate"`
	TotalCollateral *big.Int                       `json:"totalCollateral"`
	CollateralRate  core.Uint256Parameter[float64] `json:"collateralRate"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new TokenReth contract binding
func NewTokenReth(rp *rocketpool.RocketPool) (*TokenReth, error) {
	// Create the contract
	reth, err := rp.GetContract(rocketpool.ContractName_RocketTokenRETH)
	if err != nil {
		return nil, fmt.Errorf("error getting rETH contract: %w", err)
	}

	return &TokenReth{
		TokenRethDetails: &TokenRethDetails{},
		rp:               rp,
		reth:             reth,
	}, nil
}

// =============
// === Calls ===
// =============

// === Core ERC-20 functions ===

// Get the rETH total supply
func (c *TokenReth) GetTotalSupply(mc *batch.MultiCaller) {
	core.AddCall(mc, c.reth, &c.TotalSupply, "totalSupply")
}

// Get the rETH balance of an address
func (c *TokenReth) GetBalance(mc *batch.MultiCaller, balance_Out **big.Int, address common.Address) {
	core.AddCall(mc, c.reth, balance_Out, "balanceOf", address)
}

// Get the rETH spending allowance of an address and spender
func (c *TokenReth) GetAllowance(mc *batch.MultiCaller, allowance_Out **big.Int, owner common.Address, spender common.Address) {
	core.AddCall(mc, c.reth, allowance_Out, "allowance", owner, spender)
}

// === rETH functions ===

// Get the ETH balance of the rETH contract
func (c *TokenReth) GetContractEthBalance(opts *bind.CallOpts) (*big.Int, error) {
	var blockNumber *big.Int
	if opts != nil {
		blockNumber = opts.BlockNumber
	}
	return c.rp.Client.BalanceAt(context.Background(), *c.reth.Address, blockNumber)
}

// Get the ETH value of an amount of rETH
func (c *TokenReth) GetEthValueOfReth(mc *batch.MultiCaller, value_Out **big.Int, rethAmount *big.Int) {
	core.AddCall(mc, c.reth, value_Out, "getEthValue", rethAmount)
}

// Get the rETH value of an amount of ETH
func (c *TokenReth) GetRethValueOfEth(mc *batch.MultiCaller, value_Out **big.Int, ethAmount *big.Int) {
	core.AddCall(mc, c.reth, value_Out, "getRethValue", ethAmount)
}

// Get the current ETH : rETH exchange rate
func (c *TokenReth) GetExchangeRate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.reth, &c.ExchangeRate.RawValue, "getExchangeRate")
}

// Get the total amount of ETH collateral available for rETH trades
func (c *TokenReth) GetTotalCollateral(mc *batch.MultiCaller) {
	core.AddCall(mc, c.reth, &c.TotalCollateral, "getTotalCollateral")
}

// Get the rETH collateralization rate
func (c *TokenReth) GetCollateralRate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.reth, &c.CollateralRate.RawValue, "getCollateralRate")
}

// ====================
// === Transactions ===
// ====================

// === Core ERC-20 functions ===

// Get info for approving rETH's usage by a spender
func (c *TokenReth) Approve(spender common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.reth, "approve", opts, spender, amount)
}

// Get info for transferring rETH
func (c *TokenReth) Transfer(to common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.reth, "transfer", opts, to, amount)
}

// Get info for transferring rETH from a sender
func (c *TokenReth) TransferFrom(from common.Address, to common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.reth, "transferFrom", opts, from, to, amount)
}

// === rETH functions ===

// Get info for burning rETH for ETH
func (c *TokenReth) Burn(amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.reth, "burn", opts, amount)
}
