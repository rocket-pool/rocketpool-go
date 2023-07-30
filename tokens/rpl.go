package tokens

import (
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

// Binding for RocketTokenRPL
type TokenRpl struct {
	Details  TokenRplDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for RocketTokenRPL
type TokenRplDetails struct {
	TotalSupply           *big.Int `json:"totalSupply"`
	InflationIntervalRate *big.Int `json:"inflationIntervalRate"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new TokenRpl contract binding
func NewTokenRpl(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*TokenRpl, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketTokenRPL", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting RPL contract: %w", err)
	}

	return &TokenRpl{
		Details:  TokenRplDetails{},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// === Core ERC-20 functions ===

// Get the RPL total supply
func (c *TokenRpl) GetTotalSupply(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.TotalSupply, "totalSupply")
}

// Get the RPL balance of an address
func (c *TokenRpl) GetBalance(mc *multicall.MultiCaller, balance_Out **big.Int, address common.Address) {
	multicall.AddCall(mc, c.contract, balance_Out, "balanceOf", address)
}

// Get the RPL spending allowance of an address and spender
func (c *TokenRpl) GetAllowance(mc *multicall.MultiCaller, allowance_Out **big.Int, owner common.Address, spender common.Address) {
	multicall.AddCall(mc, c.contract, allowance_Out, "allowance", owner, spender)
}

// === RPL functions ===

// Get the RPL inflation interval rate
func (c *TokenRpl) GetInflationIntervalRate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.InflationIntervalRate, "getInflationIntervalRate")
}

// ====================
// === Transactions ===
// ====================

// === Core ERC-20 functions ===

// Get info for approving RPL's usage by a spender
func (c *TokenRpl) Approve(spender common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "approve", opts, spender, amount)
}

// Get info for transferring RPL
func (c *TokenRpl) Transfer(to common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "transfer", opts, to, amount)
}

// Get info for transferring RPL from a sender
func (c *TokenRpl) TransferFrom(from common.Address, to common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "transferFrom", opts, from, to, amount)
}

// === RPL functions ===

// Get info for minting new RPL tokens from inflation
func (c *TokenRpl) MintInflationRPL(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "inflationMintTokens", opts)
}

// Get info for swapping fixed-supply RPL for new RPL tokens
func (c *TokenRpl) SwapFixedSupplyRplForRpl(amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "swapTokens", opts, amount)
}
