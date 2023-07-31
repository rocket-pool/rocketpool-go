package deposit

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDepositPool
type DepositPool struct {
	Details  DepositPoolDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for RocketDepositPool
type DepositPoolDetails struct {
	Balance       *big.Int `json:"balance"`
	UserBalance   *big.Int `json:"userBalance"`
	ExcessBalance *big.Int `json:"excessBalance"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DepositPool contract binding
func NewDepositPool(rp *rocketpool.RocketPool) (*DepositPool, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDepositPool)
	if err != nil {
		return nil, fmt.Errorf("error getting deposit pool contract: %w", err)
	}

	return &DepositPool{
		Details:  DepositPoolDetails{},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the deposit pool balance
func (c *DepositPool) GetBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.Balance, "getBalance")
}

// Get the deposit pool balance provided by pool stakers
func (c *DepositPool) GetUserBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.Balance, "getUserBalance")
}

// Get the excess deposit pool balance
func (c *DepositPool) GetExcessBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.Balance, "getExcessBalance")
}

// Get all basic details
func (c *DepositPool) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetBalance(mc)
	c.GetUserBalance(mc)
	c.GetExcessBalance(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for making a deposit
func (c *DepositPool) Deposit(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "deposit", opts)
}

// Get info for assigning deposits
func (c *DepositPool) AssignDeposits(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "assignDeposits", opts)
}
