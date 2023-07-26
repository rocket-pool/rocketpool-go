package network

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

const (
	NetworkBalances_ContractName string = "rocketNetworkBalances"

	networkBalances_getBalancesBlock         string = "getBalancesBlock"
	networkBalances_getTotalETHBalance       string = "getTotalETHBalance"
	networkBalances_getStakingETHBalance     string = "getStakingETHBalance"
	networkBalances_getTotalRETHSupply       string = "getTotalRETHSupply"
	networkBalances_getETHUtilizationRate    string = "getETHUtilizationRate"
	networkBalances_getLatestReportableBlock string = "getLatestReportableBlock"
	networkBalances_submitBalances           string = "submitBalances"
)

// Binding for RocketNetworkBalances
type NetworkBalances struct {
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// Creates a new NetworkBalances contract binding
func NewNetworkBalances(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*NetworkBalances, error) {
	// Create the contract
	contract, err := rp.GetContract(NetworkBalances_ContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting network balances contract: %w", err)
	}

	return &NetworkBalances{
		rp:       rp,
		contract: contract,
	}, nil
}

// ===================
// === Raw Getters ===
// ===================

// Get the block number which network balances are current for
func (c *NetworkBalances) GetBalancesBlockRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkBalances_getBalancesBlock)
}

// Get the current network total ETH balance
func (c *NetworkBalances) GetTotalETHBalance(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkBalances_getTotalETHBalance)
}

// Get the current network staking ETH balance
func (c *NetworkBalances) GetStakingETHBalance(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkBalances_getStakingETHBalance)
}

// Get the current network total rETH supply
func (c *NetworkBalances) GetTotalRETHSupply(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkBalances_getTotalRETHSupply)
}

// Get the current network ETH utilization rate
func (c *NetworkBalances) GetETHUtilizationRateRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkBalances_getETHUtilizationRate)
}

// Returns the latest block number that oracles should be reporting balances for
func (c *NetworkBalances) GetLatestReportableBalancesBlockRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkBalances_getLatestReportableBlock)
}

// =========================
// === Formatted Getters ===
// =========================

// Get the block number which network balances are current for
func (c *NetworkBalances) GetBalancesBlock(opts *bind.CallOpts) (uint64, error) {
	raw, err := c.GetBalancesBlockRaw(opts)
	if err != nil {
		return 0, err
	}
	return raw.Uint64(), nil
}

// Get the current network ETH utilization rate
func (c *NetworkBalances) GetETHUtilizationRate(rp *rocketpool.RocketPool, opts *bind.CallOpts) (float64, error) {
	raw, err := c.GetETHUtilizationRateRaw(opts)
	if err != nil {
		return 0, err
	}
	return eth.WeiToEth(raw), nil
}

// Get the current network ETH utilization rate
func (c *NetworkBalances) GetLatestReportableBalancesBlock(rp *rocketpool.RocketPool, opts *bind.CallOpts) (uint64, error) {
	raw, err := c.GetLatestReportableBalancesBlockRaw(opts)
	if err != nil {
		return 0, err
	}
	return raw.Uint64(), nil
}

// ====================
// === Transactions ===
// ====================

// Get info for network balance submission
func (c *NetworkBalances) SubmitBalances(rp *rocketpool.RocketPool, block uint64, totalEth, stakingEth, rethSupply *big.Int, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, networkBalances_submitBalances, opts, block, totalEth, stakingEth, rethSupply)
}
