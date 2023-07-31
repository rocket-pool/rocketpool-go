package protocol

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProtocolSettingsNetwork
type DaoProtocolSettingsNetwork struct {
	Details             DaoProtocolSettingsNetworkDetails
	rp                  *rocketpool.RocketPool
	contract            *core.Contract
	daoProtocolContract *protocol.DaoProtocol
}

// Details for RocketDAOProtocolSettingsNetwork
type DaoProtocolSettingsNetworkDetails struct {
	OracleDaoConsensusThreshold core.Parameter[float64] `json:"oracleDaoConsensusThreshold"`
	IsSubmitBalancesEnabled     bool                    `json:"isSubmitBalancesEnabled"`
	SubmitBalancesFrequency     core.Parameter[uint64]  `json:"submitBalancesFrequency"`
	IsSubmitPricesEnabled       bool                    `json:"isSubmitPricesEnabled"`
	SubmitPricesFrequency       core.Parameter[uint64]  `json:"submitPricesFrequency"`
	MinimumNodeFee              core.Parameter[float64] `json:"minimumNodeFee"`
	TargetNodeFee               core.Parameter[float64] `json:"targetNodeFee"`
	MaximumNodeFee              core.Parameter[float64] `json:"maximumNodeFee"`
	NodeFeeDemandRange          *big.Int                `json:"nodeFeeDemandRange"`
	TargetRethCollateralRate    core.Parameter[float64] `json:"targetRethCollateralRate"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProtocolSettingsNetwork contract binding
func NewDaoProtocolSettingsNetwork(rp *rocketpool.RocketPool, daoProtocolContract *protocol.DaoProtocol, opts *bind.CallOpts) (*DaoProtocolSettingsNetwork, error) {
	// Create the contract
	contract, err := rp.GetContract(networkSettingsContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO protocol settings network contract: %w", err)
	}

	return &DaoProtocolSettingsNetwork{
		Details:             DaoProtocolSettingsNetworkDetails{},
		rp:                  rp,
		contract:            contract,
		daoProtocolContract: daoProtocolContract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the threshold of Oracle DAO nodes that must reach consensus on oracle data to commit it
func (c *DaoProtocolSettingsNetwork) GetOracleDaoConsensusThreshold(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.OracleDaoConsensusThreshold.RawValue, "getNodeConsensusThreshold")
}

// Check if network balance submission is enabled
func (c *DaoProtocolSettingsNetwork) GetSubmitBalancesEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.IsSubmitBalancesEnabled, "getSubmitBalancesEnabled")
}

// Get the frequency, in blocks, at which network balances should be submitted by the Oracle DAO
func (c *DaoProtocolSettingsNetwork) GetSubmitBalancesFrequency(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.SubmitBalancesFrequency.RawValue, "getSubmitBalancesFrequency")
}

// Check if network price submission is enabled
func (c *DaoProtocolSettingsNetwork) GetSubmitPricesEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.IsSubmitPricesEnabled, "getSubmitPricesEnabled")
}

// Get the frequency, in blocks, at which network prices should be submitted by the Oracle DAO
func (c *DaoProtocolSettingsNetwork) GetSubmitPricesFrequency(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.SubmitPricesFrequency.RawValue, "getSubmitPricesFrequency")
}

// Get the minimum node commission rate
func (c *DaoProtocolSettingsNetwork) GetMinimumNodeFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.MinimumNodeFee.RawValue, "getMinimumNodeFee")
}

// Get the target node commission rate
func (c *DaoProtocolSettingsNetwork) GetTargetNodeFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.TargetNodeFee.RawValue, "getTargetNodeFee")
}

// Get the maximum node commission rate
func (c *DaoProtocolSettingsNetwork) GetMaximumNodeFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.MaximumNodeFee.RawValue, "getMaximumNodeFee")
}

// Get the range of node demand values to base fee calculations on
func (c *DaoProtocolSettingsNetwork) GetNodeFeeDemandRange(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.NodeFeeDemandRange, "getNodeFeeDemandRange")
}

// Get the target collateralization rate for the rETH contract as a fraction
func (c *DaoProtocolSettingsNetwork) GetTargetRethCollateralRate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.TargetRethCollateralRate.RawValue, "getTargetRethCollateralRate")
}

// Get all basic details
func (c *DaoProtocolSettingsNetwork) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetOracleDaoConsensusThreshold(mc)
	c.GetSubmitBalancesEnabled(mc)
	c.GetSubmitBalancesFrequency(mc)
	c.GetSubmitPricesEnabled(mc)
	c.GetSubmitPricesFrequency(mc)
	c.GetMinimumNodeFee(mc)
	c.GetTargetNodeFee(mc)
	c.GetMaximumNodeFee(mc)
	c.GetNodeFeeDemandRange(mc)
	c.GetTargetRethCollateralRate(mc)
}

// ====================
// === Transactions ===
// ====================

func (c *DaoProtocolSettingsAuction) BootstrapNodeConsensusThreshold(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(networkSettingsContractName, "network.consensus.threshold", eth.EthToWei(value), opts)
}

func (c *DaoProtocolSettingsAuction) BootstrapSubmitBalancesEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapBool(networkSettingsContractName, "network.submit.balances.enabled", value, opts)
}

func (c *DaoProtocolSettingsAuction) BootstrapSubmitBalancesFrequency(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(networkSettingsContractName, "network.submit.balances.frequency", big.NewInt(int64(value)), opts)
}

func (c *DaoProtocolSettingsAuction) BootstrapSubmitPricesEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapBool(networkSettingsContractName, "network.submit.prices.enabled", value, opts)
}

func (c *DaoProtocolSettingsAuction) BootstrapSubmitPricesFrequency(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(networkSettingsContractName, "network.submit.prices.frequency", big.NewInt(int64(value)), opts)
}

func (c *DaoProtocolSettingsAuction) BootstrapMinimumNodeFee(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(networkSettingsContractName, "network.node.fee.minimum", eth.EthToWei(value), opts)
}

func (c *DaoProtocolSettingsAuction) BootstrapTargetNodeFee(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(networkSettingsContractName, "network.node.fee.target", eth.EthToWei(value), opts)
}

func (c *DaoProtocolSettingsAuction) BootstrapMaximumNodeFee(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(networkSettingsContractName, "network.node.fee.maximum", eth.EthToWei(value), opts)
}

func (c *DaoProtocolSettingsAuction) BootstrapNodeFeeDemandRange(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(networkSettingsContractName, "network.node.fee.demand.range", value, opts)
}

func (c *DaoProtocolSettingsAuction) BootstrapTargetRethCollateralRate(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(networkSettingsContractName, "network.reth.collateral.target", eth.EthToWei(value), opts)
}
