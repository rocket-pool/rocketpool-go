package beacon

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
)

// ===============
// === Structs ===
// ===============

// Binding for Beacon Deposit
type BeaconDeposit struct {
	// The deposit root for new deposits
	DepositRoot *core.SimpleField[common.Hash]

	// === Internal fields ===
	cd *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new Beacon Deposit contract binding
func NewBeaconDeposit(rp *rocketpool.RocketPool) (*BeaconDeposit, error) {
	// Create the contract
	cd, err := rp.GetContract(rocketpool.ContractName_CasperDeposit)
	if err != nil {
		return nil, fmt.Errorf("error getting Beacon deposit contract: %w", err)
	}

	return &BeaconDeposit{
		DepositRoot: core.NewSimpleField[common.Hash](cd, "get_deposit_root"),

		cd: cd,
	}, nil
}

// ====================
// === Transactions ===
// ====================

// Deposit to the Beacon contract, creating a new validator
func (c *BeaconDeposit) Deposit(opts *bind.TransactOpts, pubkey types.ValidatorPubkey, withdrawalCredentials common.Hash, signature types.ValidatorSignature, depositDataRoot common.Hash) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.cd, "deposit", opts, pubkey, withdrawalCredentials, signature, depositDataRoot)
}
