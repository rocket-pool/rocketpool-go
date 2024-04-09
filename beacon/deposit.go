package beacon

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/v2/core"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for Beacon Deposit
type BeaconDeposit struct {
	// The deposit root for new deposits
	DepositRoot *core.SimpleField[common.Hash]

	// === Internal fields ===
	cd    *core.Contract
	txMgr *eth.TransactionManager
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

		cd:    cd,
		txMgr: rp.GetTransactionManager(),
	}, nil
}

// ====================
// === Transactions ===
// ====================

// Deposit to the Beacon contract, creating a new validator
func (c *BeaconDeposit) Deposit(opts *bind.TransactOpts, pubkey beacon.ValidatorPubkey, withdrawalCredentials common.Hash, signature beacon.ValidatorSignature, depositDataRoot common.Hash) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.cd.Contract, "deposit", opts, pubkey, withdrawalCredentials, signature, depositDataRoot)
}
