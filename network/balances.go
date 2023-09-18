package network

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNetworkBalances
type NetworkBalances struct {
	*NetworkBalancesDetails
	rp *rocketpool.RocketPool
	nb *core.Contract
}

// Details for network balances
type NetworkBalancesDetails struct {
	BalancesBlock                 core.Parameter[uint64]  `json:"balancesBlock"`
	TotalETHBalance               *big.Int                `json:"totalEthBalance"`
	StakingETHBalance             *big.Int                `json:"stakingEthBalance"`
	TotalRETHSupply               *big.Int                `json:"totalRethSupply"`
	EthUtilizationRate            core.Parameter[float64] `json:"ethUtilizationRate"`
	LatestReportableBalancesBlock core.Parameter[uint64]  `json:"latestReportableBalancesBlock"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NetworkBalances contract binding
func NewNetworkBalances(rp *rocketpool.RocketPool) (*NetworkBalances, error) {
	// Create the contract
	nb, err := rp.GetContract(rocketpool.ContractName_RocketNetworkBalances)
	if err != nil {
		return nil, fmt.Errorf("error getting network balances contract: %w", err)
	}

	return &NetworkBalances{
		NetworkBalancesDetails: &NetworkBalancesDetails{},
		rp:                     rp,
		nb:                     nb,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the block number which network balances are current for
func (c *NetworkBalances) GetBalancesBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nb, &c.BalancesBlock.RawValue, "getBalancesBlock")
}

// Get the current network total ETH balance
func (c *NetworkBalances) GetTotalETHBalance(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nb, &c.TotalETHBalance, "getTotalETHBalance")
}

// Get the current network staking ETH balance
func (c *NetworkBalances) GetStakingETHBalance(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nb, &c.StakingETHBalance, "getStakingETHBalance")
}

// Get the current network total rETH supply
func (c *NetworkBalances) GetTotalRETHSupply(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nb, &c.TotalRETHSupply, "getTotalRETHSupply")
}

// Get the current network ETH utilization rate
func (c *NetworkBalances) GetEthUtilizationRate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nb, &c.EthUtilizationRate.RawValue, "getETHUtilizationRate")
}

// Returns the latest block number that oracles should be reporting balances for
func (c *NetworkBalances) GetLatestReportableBalancesBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nb, &c.LatestReportableBalancesBlock.RawValue, "getLatestReportableBlock")
}

// Get all basic details
func (c *NetworkBalances) GetAllDetails(mc *batch.MultiCaller) {
	c.GetBalancesBlock(mc)
	c.GetTotalETHBalance(mc)
	c.GetStakingETHBalance(mc)
	c.GetTotalRETHSupply(mc)
	c.GetEthUtilizationRate(mc)
	c.GetLatestReportableBalancesBlock(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for network balance submission
func (c *NetworkBalances) SubmitBalances(block uint64, totalEth, stakingEth, rethSupply *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.nb, "submitBalances", opts, block, totalEth, stakingEth, rethSupply)
}

// =============
// === Utils ===
// =============

// Returns an array of block numbers for balances submissions the given trusted node has submitted since fromBlock
func (c *NetworkBalances) GetBalancesSubmissions(nodeAddress common.Address, fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) (*[]uint64, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*c.nb.Address}
	topicFilter := [][]common.Hash{{c.nb.ABI.Events["BalancesSubmitted"].ID}, {nodeAddress.Hash()}}

	// Get the event logs
	logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, big.NewInt(int64(fromBlock)), nil, nil)
	if err != nil {
		return nil, err
	}

	timestamps := make([]uint64, len(logs))
	for i, log := range logs {
		values := make(map[string]interface{})
		// Decode the event
		if c.nb.ABI.Events["BalancesSubmitted"].Inputs.UnpackIntoMap(values, log.Data) != nil {
			return nil, err
		}
		timestamps[i] = values["block"].(*big.Int).Uint64()
	}
	return &timestamps, nil
}

// Returns an array of members who submitted a balance since fromBlock
func (c *NetworkBalances) GetLatestBalancesSubmissions(fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) ([]common.Address, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*c.nb.Address}
	topicFilter := [][]common.Hash{{c.nb.ABI.Events["BalancesSubmitted"].ID}}

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
