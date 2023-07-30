package protocol

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	nodeSettingsContractName string = "rocketDAOProtocolSettingsNode"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProtocolSettingsNode
type DaoProtocolSettingsNode struct {
	Details             DaoProtocolSettingsNodeDetails
	rp                  *rocketpool.RocketPool
	contract            *core.Contract
	daoProtocolContract *protocol.DaoProtocol
}

// Details for RocketDAOProtocolSettingsNode
type DaoProtocolSettingsNodeDetails struct {
	IsRegistrationEnabled     bool                    `json:"isRegistrationEnabled"`
	IsDepositingEnabled       bool                    `json:"isDepositingEnabled"`
	AreVacantMinipoolsEnabled bool                    `json:"areVacantMinipoolsEnabled"`
	MinimumPerMinipoolStake   core.Parameter[float64] `json:"minimumPerMinipoolStake"`
	MaximumPerMinipoolStake   core.Parameter[float64] `json:"maximumPerMinipoolStake"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProtocolSettingsNode contract binding
func NewDaoProtocolSettingsNode(rp *rocketpool.RocketPool, daoProtocolContract *protocol.DaoProtocol, opts *bind.CallOpts) (*DaoProtocolSettingsNode, error) {
	// Create the contract
	contract, err := rp.GetContract(nodeSettingsContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO protocol settings node contract: %w", err)
	}

	return &DaoProtocolSettingsNode{
		Details:             DaoProtocolSettingsNodeDetails{},
		rp:                  rp,
		contract:            contract,
		daoProtocolContract: daoProtocolContract,
	}, nil
}

// =============
// === Calls ===
// =============

// Check if node registration is currently enabled
func (c *DaoProtocolSettingsNode) GetRegistrationEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.IsRegistrationEnabled, "getRegistrationEnabled")
}

// Check if node deposits are currently enabled
func (c *DaoProtocolSettingsNode) GetDepositEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.IsDepositingEnabled, "getDepositEnabled")
}

// Check if creating vacant minipools is currently enabled
func (c *DaoProtocolSettingsNode) GetVacantMinipoolsEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.AreVacantMinipoolsEnabled, "getVacantMinipoolsEnabled")
}

// Get the minimum RPL stake per minipool as a fraction of assigned user ETH
func (c *DaoProtocolSettingsNode) GetMinimumPerMinipoolStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.MinimumPerMinipoolStake.RawValue, "getMinimumPerMinipoolStake")
}

// Get the maximum RPL stake per minipool as a fraction of assigned user ETH
func (c *DaoProtocolSettingsNode) GetMaximumPerMinipoolStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.MaximumPerMinipoolStake.RawValue, "getMaximumPerMinipoolStake")
}

// Get all basic details
func (c *DaoProtocolSettingsNode) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetRegistrationEnabled(mc)
	c.GetDepositEnabled(mc)
	c.GetVacantMinipoolsEnabled(mc)
	c.GetMinimumPerMinipoolStake(mc)
	c.GetMaximumPerMinipoolStake(mc)
}

// ====================
// === Transactions ===
// ====================

func (c *DaoProtocolSettingsNode) BootstrapNodeRegistrationEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapBool(nodeSettingsContractName, "node.registration.enabled", value, opts)
}

func (c *DaoProtocolSettingsNode) BootstrapNodeDepositEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapBool(nodeSettingsContractName, "node.deposit.enabled", value, opts)
}

func (c *DaoProtocolSettingsNode) BootstrapVacantMinipoolsEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapBool(nodeSettingsContractName, "node.vacant.minipools.enabled", value, opts)
}

func (c *DaoProtocolSettingsNode) BootstrapMinimumPerMinipoolStake(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(nodeSettingsContractName, "node.per.minipool.stake.minimum", eth.EthToWei(value), opts)
}

func (c *DaoProtocolSettingsNode) BootstrapMaximumPerMinipoolStake(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(nodeSettingsContractName, "node.per.minipool.stake.maximum", eth.EthToWei(value), opts)
}
