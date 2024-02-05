package tokens

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/eth-utils/eth"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketTokenRETH
type TokenReth struct {
	// The rETH total supply
	TotalSupply *core.SimpleField[*big.Int]

	// The current ETH : rETH exchange rate
	ExchangeRate *core.FormattedUint256Field[float64]

	// The total amount of ETH collateral available for rETH trades
	TotalCollateral *core.SimpleField[*big.Int]

	// The rETH collateralization rate
	CollateralRate *core.FormattedUint256Field[float64]

	// === Internal fields ===
	rp    *rocketpool.RocketPool
	reth  *core.Contract
	txMgr *eth.TransactionManager
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
		TotalSupply:     core.NewSimpleField[*big.Int](reth, "totalSupply"),
		ExchangeRate:    core.NewFormattedUint256Field[float64](reth, "getExchangeRate"),
		TotalCollateral: core.NewSimpleField[*big.Int](reth, "getTotalCollateral"),
		CollateralRate:  core.NewFormattedUint256Field[float64](reth, "getCollateralRate"),

		rp:    rp,
		reth:  reth,
		txMgr: rp.GetTransactionManager(),
	}, nil
}

// =============
// === Calls ===
// =============

// === Core ERC-20 functions ===

// Get the rETH balance of an address
func (c *TokenReth) BalanceOf(mc *batch.MultiCaller, balance_Out **big.Int, address common.Address) {
	core.AddCall(mc, c.reth, balance_Out, "balanceOf", address)
}

// Get the rETH spending allowance of an address and spender
func (c *TokenReth) GetAllowance(mc *batch.MultiCaller, allowance_Out **big.Int, owner common.Address, spender common.Address) {
	core.AddCall(mc, c.reth, allowance_Out, "allowance", owner, spender)
}

// === rETH functions ===

// Get the ETH value of an amount of rETH
func (c *TokenReth) GetEthValueOfReth(mc *batch.MultiCaller, value_Out **big.Int, rethAmount *big.Int) {
	core.AddCall(mc, c.reth, value_Out, "getEthValue", rethAmount)
}

// Get the rETH value of an amount of ETH
func (c *TokenReth) GetRethValueOfEth(mc *batch.MultiCaller, value_Out **big.Int, ethAmount *big.Int) {
	core.AddCall(mc, c.reth, value_Out, "getRethValue", ethAmount)
}

// ====================
// === Transactions ===
// ====================

// === Core ERC-20 functions ===

// Get info for approving rETH's usage by a spender
func (c *TokenReth) Approve(spender common.Address, amount *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.reth.Contract, "approve", opts, spender, amount)
}

// Get info for transferring rETH
func (c *TokenReth) Transfer(to common.Address, amount *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.reth.Contract, "transfer", opts, to, amount)
}

// Get info for transferring rETH from a sender
func (c *TokenReth) TransferFrom(from common.Address, to common.Address, amount *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.reth.Contract, "transferFrom", opts, from, to, amount)
}

// === rETH functions ===

// Get info for burning rETH for ETH
func (c *TokenReth) Burn(amount *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.reth.Contract, "burn", opts, amount)
}
