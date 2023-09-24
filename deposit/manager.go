package deposit

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDepositPool
type DepositPoolManager struct {
	// The deposit pool balance
	Balance *core.SimpleField[*big.Int]

	// The deposit pool balance provided by pool stakers
	UserBalance *core.SimpleField[*big.Int]

	// The excess deposit pool balance
	ExcessBalance *core.SimpleField[*big.Int]

	// === Internal fields ===
	rp *rocketpool.RocketPool
	dp *core.Contract
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
		Balance:       core.NewSimpleField[*big.Int](dp, "getBalance"),
		UserBalance:   core.NewSimpleField[*big.Int](dp, "getUserBalance"),
		ExcessBalance: core.NewSimpleField[*big.Int](dp, "getExcessBalance"),

		rp: rp,
		dp: dp,
	}, nil
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
