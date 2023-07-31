package tokens

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketTokenRETH
type TokenReth struct {
	Details  TokenRethDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for RocketTokenRETH
type TokenRethDetails struct {
	TotalSupply     *big.Int                `json:"totalSupply"`
	ExchangeRate    core.Parameter[float64] `json:"exchangeRate"`
	TotalCollateral *big.Int                `json:"totalCollateral"`
	CollateralRate  core.Parameter[float64] `json:"collateralRate"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new TokenReth contract binding
func NewTokenReth(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*TokenReth, error) {
	// Create the contract
	contract, err := rp.GetContract(, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting rETH contract: %w", err)
	}

	return &TokenReth{
		Details:  TokenRethDetails{},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// === Core ERC-20 functions ===

// Get the rETH total supply
func (c *TokenReth) GetTotalSupply(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.TotalSupply, "totalSupply")
}

// Get the rETH balance of an address
func (c *TokenReth) GetBalance(mc *multicall.MultiCaller, balance_Out **big.Int, address common.Address) {
	multicall.AddCall(mc, c.contract, balance_Out, "balanceOf", address)
}

// Get the rETH spending allowance of an address and spender
func (c *TokenReth) GetAllowance(mc *multicall.MultiCaller, allowance_Out **big.Int, owner common.Address, spender common.Address) {
	multicall.AddCall(mc, c.contract, allowance_Out, "allowance", owner, spender)
}

// === rETH functions ===

// Get the ETH balance of the rETH contract
func (c *TokenReth) GetContractEthBalance(opts *bind.CallOpts) (*big.Int, error) {
	var blockNumber *big.Int
	if opts != nil {
		blockNumber = opts.BlockNumber
	}
	return c.rp.Client.BalanceAt(context.Background(), *c.contract.Address, blockNumber)
}

// Get the ETH value of an amount of rETH
func (c *TokenReth) GetEthValueOfReth(mc *multicall.MultiCaller, value_Out **big.Int, rethAmount *big.Int) {
	multicall.AddCall(mc, c.contract, value_Out, "getEthValue", rethAmount)
}

// Get the rETH value of an amount of ETH
func (c *TokenReth) GetRethValueOfEth(mc *multicall.MultiCaller, value_Out **big.Int, ethAmount *big.Int) {
	multicall.AddCall(mc, c.contract, value_Out, "getRethValue", ethAmount)
}

// Get the current ETH : rETH exchange rate
func (c *TokenReth) GetExchangeRate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.ExchangeRate.RawValue, "getExchangeRate")
}

// Get the total amount of ETH collateral available for rETH trades
func (c *TokenReth) GetTotalCollateral(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.TotalCollateral, "getTotalCollateral")
}

// Get the rETH collateralization rate
func (c *TokenReth) GetCollateralRate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.CollateralRate.RawValue, "getCollateralRate")
}

// ====================
// === Transactions ===
// ====================

// === Core ERC-20 functions ===

// Get info for approving rETH's usage by a spender
func (c *TokenReth) Approve(spender common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "approve", opts, spender, amount)
}

// Get info for transferring rETH
func (c *TokenReth) Transfer(to common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "transfer", opts, to, amount)
}

// Get info for transferring rETH from a sender
func (c *TokenReth) TransferFrom(from common.Address, to common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "transferFrom", opts, from, to, amount)
}

// === rETH functions ===

// Get info for burning rETH for ETH
func (c *TokenReth) Burn(amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "burn", opts, amount)
}
