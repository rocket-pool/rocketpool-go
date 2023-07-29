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
type DaoProtocol struct {
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProtocol contract binding
func NewDaoProtocol(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*DaoProtocol, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketDAOProtocol", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO protocol contract: %w", err)
	}

	return &DaoProtocol{
		rp:       rp,
		contract: contract,
	}, nil
}

// ====================
// === Transactions ===
// ====================

// Get info for bootstrapping a bool setting
func (c *DaoProtocol) BootstrapBool(contractName string, settingPath string, value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "bootstrapSettingBool", opts, contractName, settingPath, value)
}

// Get info for bootstrapping a uint256 setting
func (c *DaoProtocol) BootstrapUint(contractName string, settingPath string, value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "bootstrapSettingUint", opts, contractName, settingPath, value)
}

// Get info for bootstrapping an address setting
func (c *DaoProtocol) BootstrapAddress(contractName string, settingPath string, value common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "bootstrapSettingAddress", opts, contractName, settingPath, value)
}

// Get info for bootstrapping a rewards claimer
func (c *DaoProtocol) BootstrapClaimer(contractName string, amount float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "bootstrapSettingClaimer", opts, contractName, eth.EthToWei(amount))
}
