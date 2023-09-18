package beacon

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
)

// ===============
// === Structs ===
// ===============

// Binding for Beacon Deposit
type BeaconDeposit struct {
	*BeaconDepositDetails
	cd *core.Contract
}

// Details for Beacon Deposit
type BeaconDepositDetails struct {
	DepositRoot common.Hash `json:"depositRoot"`
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
		BeaconDepositDetails: &BeaconDepositDetails{},
		cd:                   cd,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the deposit root for new deposits
func (c *BeaconDeposit) GetDepositRoot(mc *batch.MultiCaller) {
	core.AddCall(mc, c.cd, &c.DepositRoot, "get_deposit_root")
}

// ====================
// === Transactions ===
// ====================

// Deposit to the Beacon contract, creating a new validator
func (c *BeaconDeposit) Deposit(opts *bind.TransactOpts, pubkey types.ValidatorPubkey, withdrawalCredentials common.Hash, signature types.ValidatorSignature, depositDataRoot common.Hash) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.cd, "deposit", opts, pubkey, withdrawalCredentials, signature, depositDataRoot)
}
