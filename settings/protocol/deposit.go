package protocol

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	depositSettingsContractName string = "rocketDAOProtocolSettingsDeposit"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProtocolSettingsDeposit
type DaoProtocolSettingsDeposit struct {
	Details             DaoProtocolSettingsDepositDetails
	rp                  *rocketpool.RocketPool
	contract            *core.Contract
	daoProtocolContract *protocol.DaoProtocol
}

// Details for RocketDAOProtocolSettingsDeposit
type DaoProtocolSettingsDepositDetails struct {
	IsDepositingEnabled          bool                   `json:"isDepositingEnabled"`
	AreDepositAssignmentsEnabled bool                   `json:"areDepositAssignmentsEnabled"`
	MinimumDeposit               *big.Int               `json:"minimumDeposit"`
	MaximumDepositPoolSize       *big.Int               `json:"maximumDepositPoolSize"`
	MaximumAssignmentsPerDeposit core.Parameter[uint64] `json:"maximumAssignmentsPerDeposit"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProtocolSettingsDeposit contract binding
func NewDaoProtocolSettingsDeposit(rp *rocketpool.RocketPool, daoProtocolContract *protocol.DaoProtocol, opts *bind.CallOpts) (*DaoProtocolSettingsDeposit, error) {
	// Create the contract
	contract, err := rp.GetContract(depositSettingsContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO protocol settings deposit contract: %w", err)
	}

	return &DaoProtocolSettingsDeposit{
		Details:             DaoProtocolSettingsDepositDetails{},
		rp:                  rp,
		contract:            contract,
		daoProtocolContract: daoProtocolContract,
	}, nil
}

// =============
// === Calls ===
// =============

// Check if deposits are currently enabled
func (c *DaoProtocolSettingsDeposit) GetDepositEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.IsDepositingEnabled, "getDepositEnabled")
}

// Check if deposit assignments are currently enabled
func (c *DaoProtocolSettingsDeposit) GetAssignDepositsEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.AreDepositAssignmentsEnabled, "getAssignDepositsEnabled")
}

// Get the minimum deposit to the deposit pool
func (c *DaoProtocolSettingsDeposit) GetMinimumDeposit(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.MinimumDeposit, "getMinimumDeposit")
}

// Get the maximum size of the deposit pool
func (c *DaoProtocolSettingsDeposit) GetMaximumDepositPoolSize(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.MaximumDepositPoolSize, "getMaximumDepositPoolSize")
}

// Get the maximum assignments per deposit transaction
func (c *DaoProtocolSettingsDeposit) GetMaximumDepositAssignments(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.MaximumAssignmentsPerDeposit.RawValue, "getMaximumDepositAssignments")
}

// Get all basic details
func (c *DaoProtocolSettingsDeposit) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetDepositEnabled(mc)
	c.GetAssignDepositsEnabled(mc)
	c.GetMinimumDeposit(mc)
	c.GetMaximumDepositPoolSize(mc)
	c.GetMaximumDepositAssignments(mc)
}

// ====================
// === Transactions ===
// ====================

// Set the deposit enabled flag
func (c *DaoProtocolSettingsDeposit) BootstrapDepositEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapBool(depositSettingsContractName, "deposit.enabled", value, opts)
}

// Set the deposit assignments enabled flag
func (c *DaoProtocolSettingsDeposit) BootstrapAssignDepositsEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapBool(depositSettingsContractName, "deposit.assign.enabled", value, opts)
}

// Set the minimum deposit amount
func (c *DaoProtocolSettingsDeposit) BootstrapMinimumDeposit(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(depositSettingsContractName, "deposit.minimum", value, opts)
}

// Set the maximum deposit pool size
func (c *DaoProtocolSettingsDeposit) BootstrapMaximumDepositPoolSize(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(depositSettingsContractName, "deposit.pool.maximum", value, opts)
}

// Set the max assignments per deposit
func (c *DaoProtocolSettingsDeposit) BootstrapMaximumDepositAssignments(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(depositSettingsContractName, "deposit.assign.maximum", big.NewInt(int64(value)), opts)
}
