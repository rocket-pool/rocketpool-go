package tokens

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/eth"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketTokenRPLFixedSupply
type TokenRplFixedSupply struct {
	// The fixed-supply RPL total supply
	TotalSupply *core.SimpleField[*big.Int]

	// === Internal fields ===
	rp    *rocketpool.RocketPool
	fsrpl *core.Contract
	txMgr *eth.TransactionManager
}

// Details for RocketTokenRPLFixedSupply
type TokenRplFixedSupplyDetails struct {
}

// ====================
// === Constructors ===
// ====================

// Creates a new TokenRplFixedSupply contract binding
func NewTokenRplFixedSupply(rp *rocketpool.RocketPool) (*TokenRplFixedSupply, error) {
	// Create the contract
	fsrpl, err := rp.GetContract(rocketpool.ContractName_RocketTokenRPLFixedSupply)
	if err != nil {
		return nil, fmt.Errorf("error getting RPL fixed supply contract: %w", err)
	}

	return &TokenRplFixedSupply{
		TotalSupply: core.NewSimpleField[*big.Int](fsrpl, "totalSupply"),

		rp:    rp,
		fsrpl: fsrpl,
		txMgr: rp.GetTransactionManager(),
	}, nil
}

// =============
// === Calls ===
// =============

// === Core ERC-20 functions ===

// Get the fixed-supply RPL balance of an address
func (c *TokenRplFixedSupply) BalanceOf(mc *batch.MultiCaller, balance_Out **big.Int, address common.Address) {
	core.AddCall(mc, c.fsrpl, balance_Out, "balanceOf", address)
}

// Get the fixed-supply RPL spending allowance of an address and spender
func (c *TokenRplFixedSupply) GetAllowance(mc *batch.MultiCaller, allowance_Out **big.Int, owner common.Address, spender common.Address) {
	core.AddCall(mc, c.fsrpl, allowance_Out, "allowance", owner, spender)
}

// ====================
// === Transactions ===
// ====================

// === Core ERC-20 functions ===

// Get info for approving fixed-supply RPL's usage by a spender
func (c *TokenRplFixedSupply) Approve(spender common.Address, amount *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.fsrpl.Contract, "approve", opts, spender, amount)
}

// Get info for transferring fixed-supply RPL
func (c *TokenRplFixedSupply) Transfer(to common.Address, amount *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.fsrpl.Contract, "transfer", opts, to, amount)
}

// Get info for transferring fixed-supply RPL from a sender
func (c *TokenRplFixedSupply) TransferFrom(from common.Address, to common.Address, amount *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.fsrpl.Contract, "transferFrom", opts, from, to, amount)
}
