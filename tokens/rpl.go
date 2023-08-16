package tokens

import (
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

// Binding for RocketTokenRPL
type TokenRpl struct {
	Details  TokenRplDetails
	Contract *core.Contract
	rp       *rocketpool.RocketPool
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
func NewTokenRpl(rp *rocketpool.RocketPool) (*TokenRpl, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketTokenRPL)
	if err != nil {
		return nil, fmt.Errorf("error getting RPL contract: %w", err)
	}

	return &TokenRpl{
		Details:  TokenRplDetails{},
		Contract: contract,
		rp:       rp,
	}, nil
}

// =============
// === Calls ===
// =============

// === Core ERC-20 functions ===

// Get the RPL total supply
func (c *TokenRpl) GetTotalSupply(mc *batch.MultiCaller) {
	core.AddCall(mc, c.Contract, &c.Details.TotalSupply, "totalSupply")
}

// Get the RPL balance of an address
func (c *TokenRpl) GetBalance(mc *batch.MultiCaller, balance_Out **big.Int, address common.Address) {
	core.AddCall(mc, c.Contract, balance_Out, "balanceOf", address)
}

// Get the RPL spending allowance of an address and spender
func (c *TokenRpl) GetAllowance(mc *batch.MultiCaller, allowance_Out **big.Int, owner common.Address, spender common.Address) {
	core.AddCall(mc, c.Contract, allowance_Out, "allowance", owner, spender)
}

// === RPL functions ===

// Get the RPL inflation interval rate
func (c *TokenRpl) GetInflationIntervalRate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.Contract, &c.Details.InflationIntervalRate, "getInflationIntervalRate")
}

// ====================
// === Transactions ===
// ====================

// === Core ERC-20 functions ===

// Get info for approving RPL's usage by a spender
func (c *TokenRpl) Approve(spender common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "approve", opts, spender, amount)
}

// Get info for transferring RPL
func (c *TokenRpl) Transfer(to common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "transfer", opts, to, amount)
}

// Get info for transferring RPL from a sender
func (c *TokenRpl) TransferFrom(from common.Address, to common.Address, amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "transferFrom", opts, from, to, amount)
}

// === RPL functions ===

// Get info for minting new RPL tokens from inflation
func (c *TokenRpl) MintInflationRPL(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "inflationMintTokens", opts)
}

// Get info for swapping fixed-supply RPL for new RPL tokens
func (c *TokenRpl) SwapFixedSupplyRplForRpl(amount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "swapTokens", opts, amount)
}
