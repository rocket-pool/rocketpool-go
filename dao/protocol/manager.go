package protocol

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProtocol
type ProtocolDaoManager struct {
	// Settings for the Protocol DAO
	Settings *ProtocolDaoSettings

	// === Internal fields ===
	rp *rocketpool.RocketPool
	dp *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new ProtocolDaoManager contract binding
func NewProtocolDaoManager(rp *rocketpool.RocketPool) (*ProtocolDaoManager, error) {
	// Create the contracts
	dp, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocol)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO manager contract: %w", err)
	}

	pdaoMgr := &ProtocolDaoManager{
		rp: rp,
		dp: dp,
	}
	settings, err := newProtocolDaoSettings(pdaoMgr)
	if err != nil {
		return nil, fmt.Errorf("error creating Protocol DAO settings binding: %w", err)
	}
	pdaoMgr.Settings = settings
	return pdaoMgr, nil
}

// ====================
// === Transactions ===
// ====================

// Get info for bootstrapping a bool setting
func (c *ProtocolDaoManager) BootstrapBool(contractName rocketpool.ContractName, settingPath string, value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingBool", opts, contractName, settingPath, value)
}

// Get info for bootstrapping a uint256 setting
func (c *ProtocolDaoManager) BootstrapUint(contractName rocketpool.ContractName, settingPath string, value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingUint", opts, contractName, settingPath, value)
}

// Get info for bootstrapping an address setting
func (c *ProtocolDaoManager) BootstrapAddress(contractName rocketpool.ContractName, settingPath string, value common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingAddress", opts, contractName, settingPath, value)
}

// Get info for bootstrapping a rewards claimer
func (c *ProtocolDaoManager) BootstrapClaimer(contractName rocketpool.ContractName, amount float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingClaimer", opts, contractName, eth.EthToWei(amount))
}
