package protocol

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

const (
	// Contract names
	DaoProtocol_ContractName string = "rocketDAOProtocol"

	// Transactions
	daoProtocol_bootstrapSettingBool    string = "bootstrapSettingBool"
	daoProtocol_bootstrapSettingUint    string = "bootstrapSettingUint"
	daoProtocol_bootstrapSettingAddress string = "bootstrapSettingAddress"
	daoProtocol_bootstrapSettingClaimer string = "bootstrapSettingClaimer"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProtocol
type DaoProtocol struct {
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProtocol contract binding
func NewDaoProtocol(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*DaoProtocol, error) {
	// Create the contract
	contract, err := rp.GetContract(DaoProtocol_ContractName, opts)
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
func (c *DaoProtocol) BootstrapBool(contractName string, settingPath string, value bool, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, daoProtocol_bootstrapSettingBool, opts, contractName, settingPath, value)
}

// Get info for bootstrapping a uint256 setting
func (c *DaoProtocol) BootstrapUint(contractName string, settingPath string, value *big.Int, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, daoProtocol_bootstrapSettingUint, opts, contractName, settingPath, value)
}

// Get info for bootstrapping an address setting
func (c *DaoProtocol) BootstrapAddress(contractName string, settingPath string, value common.Address, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, daoProtocol_bootstrapSettingAddress, opts, contractName, settingPath, value)
}

// Get info for bootstrapping a rewards claimer
func (c *DaoProtocol) BootstrapClaimer(contractName string, amount float64, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, daoProtocol_bootstrapSettingClaimer, opts, contractName, eth.EthToWei(amount))
}
