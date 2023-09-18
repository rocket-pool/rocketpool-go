package deposit

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDepositPool
type DepositPoolManager struct {
	*DepositPoolManagerDetails
	rp *rocketpool.RocketPool
	dp *core.Contract
}

// Details for RocketDepositPool
type DepositPoolManagerDetails struct {
	Balance       *big.Int `json:"balance"`
	UserBalance   *big.Int `json:"userBalance"`
	ExcessBalance *big.Int `json:"excessBalance"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DepositPool contract binding
func NewDepositPoolManager(rp *rocketpool.RocketPool) (*DepositPoolManager, error) {
	// Create the contract
	dp, err := rp.GetContract(rocketpool.ContractName_RocketDepositPool)
	if err != nil {
		return nil, fmt.Errorf("error getting deposit pool contract: %w", err)
	}

	return &DepositPoolManager{
		DepositPoolManagerDetails: &DepositPoolManagerDetails{},
		rp:                        rp,
		dp:                        dp,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the deposit pool balance
func (c *DepositPoolManager) GetBalance(mc *batch.MultiCaller) {
	core.AddCall(mc, c.dp, &c.Balance, "getBalance")
}

// Get the deposit pool balance provided by pool stakers
func (c *DepositPoolManager) GetUserBalance(mc *batch.MultiCaller) {
	core.AddCall(mc, c.dp, &c.Balance, "getUserBalance")
}

// Get the excess deposit pool balance
func (c *DepositPoolManager) GetExcessBalance(mc *batch.MultiCaller) {
	core.AddCall(mc, c.dp, &c.Balance, "getExcessBalance")
}

// Get all basic details
func (c *DepositPoolManager) GetAllDetails(mc *batch.MultiCaller) {
	c.GetBalance(mc)
	c.GetUserBalance(mc)
	c.GetExcessBalance(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for making a deposit
func (c *DepositPoolManager) Deposit(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "deposit", opts)
}

// Get info for assigning deposits
func (c *DepositPoolManager) AssignDeposits(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "assignDeposits", opts)
}
