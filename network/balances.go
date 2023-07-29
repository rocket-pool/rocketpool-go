package network

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNetworkBalances
type NetworkBalances struct {
	Details  NetworkBalancesDetails
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// Details for network balances
type NetworkBalancesDetails struct {
	BalancesBlock                 rocketpool.Parameter[uint64]  `json:"balancesBlock"`
	TotalETHBalance               *big.Int                      `json:"totalEthBalance"`
	StakingETHBalance             *big.Int                      `json:"stakingEthBalance"`
	TotalRETHSupply               *big.Int                      `json:"totalRethSupply"`
	ETHUtilizationRate            rocketpool.Parameter[float64] `json:"ethUtilizationRate"`
	LatestReportableBalancesBlock rocketpool.Parameter[uint64]  `json:"latestReportableBalancesBlock"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NetworkBalances contract binding
func NewNetworkBalances(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*NetworkBalances, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketNetworkBalances", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting network balances contract: %w", err)
	}

	return &NetworkBalances{
		Details:  NetworkBalancesDetails{},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the block number which network balances are current for
func (c *NetworkBalances) GetBalancesBlock(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.BalancesBlock.RawValue, "getBalancesBlock")
}

// Get the current network total ETH balance
func (c *NetworkBalances) GetTotalETHBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.TotalETHBalance, "getTotalETHBalance")
}

// Get the current network staking ETH balance
func (c *NetworkBalances) GetStakingETHBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.StakingETHBalance, "getStakingETHBalance")
}

// Get the current network total rETH supply
func (c *NetworkBalances) GetTotalRETHSupply(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.TotalRETHSupply, "getTotalRETHSupply")
}

// Get the current network ETH utilization rate
func (c *NetworkBalances) GetETHUtilizationRate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.ETHUtilizationRate.RawValue, "getETHUtilizationRate")
}

// Returns the latest block number that oracles should be reporting balances for
func (c *NetworkBalances) GetLatestReportableBalancesBlock(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.LatestReportableBalancesBlock.RawValue, "getLatestReportableBlock")
}

// Get all basic details
func (c *NetworkBalances) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetBalancesBlock(mc)
	c.GetTotalETHBalance(mc)
	c.GetStakingETHBalance(mc)
	c.GetTotalRETHSupply(mc)
	c.GetETHUtilizationRate(mc)
	c.GetLatestReportableBalancesBlock(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for network balance submission
func (c *NetworkBalances) SubmitBalances(block uint64, totalEth, stakingEth, rethSupply *big.Int, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, "submitBalances", opts, block, totalEth, stakingEth, rethSupply)
}

// =============
// === Utils ===
// =============

// Returns an array of block numbers for balances submissions the given trusted node has submitted since fromBlock
func (c *NetworkBalances) GetBalancesSubmissions(nodeAddress common.Address, fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) (*[]uint64, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*c.contract.Address}
	topicFilter := [][]common.Hash{{c.contract.ABI.Events["BalancesSubmitted"].ID}, {nodeAddress.Hash()}}

	// Get the event logs
	logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, big.NewInt(int64(fromBlock)), nil, nil)
	if err != nil {
		return nil, err
	}

	timestamps := make([]uint64, len(logs))
	for i, log := range logs {
		values := make(map[string]interface{})
		// Decode the event
		if c.contract.ABI.Events["BalancesSubmitted"].Inputs.UnpackIntoMap(values, log.Data) != nil {
			return nil, err
		}
		timestamps[i] = values["block"].(*big.Int).Uint64()
	}
	return &timestamps, nil
}

// Returns an array of members who submitted a balance since fromBlock
func (c *NetworkBalances) GetLatestBalancesSubmissions(fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) ([]common.Address, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*c.contract.Address}
	topicFilter := [][]common.Hash{{c.contract.ABI.Events["BalancesSubmitted"].ID}}

	// Get the event logs
	logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, big.NewInt(int64(fromBlock)), nil, nil)
	if err != nil {
		return nil, err
	}

	results := make([]common.Address, len(logs))
	for i, log := range logs {
		// Topic 0 is the event, topic 1 is the "from" address
		address := common.BytesToAddress(log.Topics[1].Bytes())
		results[i] = address
	}
	return results, nil
}
