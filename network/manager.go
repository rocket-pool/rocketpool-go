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

// Binding for the Network Manager
type NetworkManager struct {
	*NetworkManagerDetails
	rp               *rocketpool.RocketPool
	networkBalances  *core.Contract
	networkFees      *core.Contract
	networkPenalties *core.Contract
	networkPrices    *core.Contract
}

// Details for network balances
type NetworkManagerDetails struct {
	// Balances
	BalancesBlock                 core.Uint256Parameter[uint64]  `json:"balancesBlock"`
	TotalETHBalance               *big.Int                       `json:"totalEthBalance"`
	StakingETHBalance             *big.Int                       `json:"stakingEthBalance"`
	TotalRETHSupply               *big.Int                       `json:"totalRethSupply"`
	EthUtilizationRate            core.Uint256Parameter[float64] `json:"ethUtilizationRate"`
	LatestReportableBalancesBlock core.Uint256Parameter[uint64]  `json:"latestReportableBalancesBlock"`

	// Fees
	NodeDemand      *big.Int                       `json:"nodeDemand"`
	NodeFee         core.Uint256Parameter[float64] `json:"nodeFee"`
	NodeFeeByDemand core.Uint256Parameter[float64] `json:"nodeFeeByDemand"`

	// Prices
	PricesBlock                 core.Uint256Parameter[uint64]  `json:"pricesBlock"`
	RplPrice                    core.Uint256Parameter[float64] `json:"rplPrice"`
	LatestReportablePricesBlock core.Uint256Parameter[uint64]  `json:"latestReportablePricesBlock"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NetworkManager
func NewNetworkManager(rp *rocketpool.RocketPool) (*NetworkManager, error) {
	// Create the contracts
	networkBalances, err := rp.GetContract(rocketpool.ContractName_RocketNetworkBalances)
	if err != nil {
		return nil, fmt.Errorf("error getting network balances contract: %w", err)
	}
	// Create the contract
	networkFees, err := rp.GetContract(rocketpool.ContractName_RocketNetworkFees)
	if err != nil {
		return nil, fmt.Errorf("error getting network fees contract: %w", err)
	}
	networkPenalties, err := rp.GetContract(rocketpool.ContractName_RocketNetworkPenalties)
	if err != nil {
		return nil, fmt.Errorf("error getting network penalties contract: %w", err)
	}
	networkPrices, err := rp.GetContract(rocketpool.ContractName_RocketNetworkPrices)
	if err != nil {
		return nil, fmt.Errorf("error getting network prices contract: %w", err)
	}

	return &NetworkManager{
		NetworkManagerDetails: &NetworkManagerDetails{},
		rp:                    rp,
		networkBalances:       networkBalances,
		networkFees:           networkFees,
		networkPenalties:      networkPenalties,
		networkPrices:         networkPrices,
	}, nil
}

// =============
// === Calls ===
// =============

// === NetworkBalance ===

// Get the block number which network balances are current for
func (c *NetworkManager) GetBalancesBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkBalances, &c.BalancesBlock.RawValue, "getBalancesBlock")
}

// Get the current network total ETH balance
func (c *NetworkManager) GetTotalETHBalance(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkBalances, &c.TotalETHBalance, "getTotalETHBalance")
}

// Get the current network staking ETH balance
func (c *NetworkManager) GetStakingETHBalance(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkBalances, &c.StakingETHBalance, "getStakingETHBalance")
}

// Get the current network total rETH supply
func (c *NetworkManager) GetTotalRETHSupply(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkBalances, &c.TotalRETHSupply, "getTotalRETHSupply")
}

// Get the current network ETH utilization rate
func (c *NetworkManager) GetEthUtilizationRate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkBalances, &c.EthUtilizationRate.RawValue, "getETHUtilizationRate")
}

// Returns the latest block number that oracles should be reporting balances for
func (c *NetworkManager) GetLatestReportableBalancesBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkBalances, &c.LatestReportableBalancesBlock.RawValue, "getLatestReportableBlock")
}

// === NetworkFees ===

// Get the current network node demand in ETH
func (c *NetworkManager) GetNodeDemand(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkFees, &c.NodeDemand, "getNodeDemand")
}

// Get the current network node commission rate
func (c *NetworkManager) GetNodeFee(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkFees, &c.NodeFee.RawValue, "getNodeFee")
}

// Get the network node fee for a node demand value
func (c *NetworkManager) GetNodeFeeByDemand(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkFees, &c.NodeFeeByDemand.RawValue, "getNodeFeeByDemand")
}

// === NetworkPenalties ===

// Get info for minipool penalty submission
func (c *NetworkManager) SubmitPenalty(minipoolAddress common.Address, block *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.networkPenalties, "submitPenalty", opts, minipoolAddress, block)
}

// === NetworkPrices ===

// Get the block number which network prices are current for
func (c *NetworkManager) GetPricesBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkPrices, &c.PricesBlock.RawValue, "getPricesBlock")
}

// Get the current network RPL price in ETH
func (c *NetworkManager) GetRplPrice(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkPrices, &c.RplPrice.RawValue, "getRPLPrice")
}

// Returns the latest block number that oracles should be reporting prices for
func (c *NetworkManager) GetLatestReportablePricesBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.networkPrices, &c.LatestReportablePricesBlock.RawValue, "getLatestReportableBlock")
}

// ====================
// === Transactions ===
// ====================

// === NetworkBalances ===

// Get info for network balance submission
func (c *NetworkManager) SubmitBalances(block uint64, totalEth, stakingEth, rethSupply *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.networkBalances, "submitBalances", opts, block, totalEth, stakingEth, rethSupply)
}

// === NetworkPrices ===

// Get info for network price submission
func (c *NetworkManager) SubmitPrices(block uint64, rplPrice *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.networkPrices, "submitPrices", opts, big.NewInt(int64(block)), rplPrice)
}

// =============
// === Utils ===
// =============

// Returns an array of block numbers for balances submissions the given trusted node has submitted since fromBlock
func (c *NetworkManager) GetBalancesSubmissions(nodeAddress common.Address, fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) (*[]uint64, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*c.networkBalances.Address}
	topicFilter := [][]common.Hash{{c.networkBalances.ABI.Events["BalancesSubmitted"].ID}, {nodeAddress.Hash()}}

	// Get the event logs
	logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, big.NewInt(int64(fromBlock)), nil, nil)
	if err != nil {
		return nil, err
	}

	timestamps := make([]uint64, len(logs))
	for i, log := range logs {
		values := make(map[string]interface{})
		// Decode the event
		if c.networkBalances.ABI.Events["BalancesSubmitted"].Inputs.UnpackIntoMap(values, log.Data) != nil {
			return nil, err
		}
		timestamps[i] = values["block"].(*big.Int).Uint64()
	}
	return &timestamps, nil
}

// Returns an array of members who submitted a balance since fromBlock
func (c *NetworkManager) GetLatestBalancesSubmissions(fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) ([]common.Address, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*c.networkBalances.Address}
	topicFilter := [][]common.Hash{{c.networkBalances.ABI.Events["BalancesSubmitted"].ID}}

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

// Returns an array of block numbers for prices submissions the given trusted node has submitted since fromBlock
func (c *NetworkManager) GetPricesSubmissions(nodeAddress common.Address, fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) (*[]uint64, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*c.networkPrices.Address}
	topicFilter := [][]common.Hash{{c.networkPrices.ABI.Events["PricesSubmitted"].ID}, {nodeAddress.Hash()}}

	// Get the event logs
	logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, big.NewInt(int64(fromBlock)), nil, nil)
	if err != nil {
		return nil, err
	}
	timestamps := make([]uint64, len(logs))
	for i, log := range logs {
		values := make(map[string]interface{})
		// Decode the event
		if c.networkPrices.ABI.Events["PricesSubmitted"].Inputs.UnpackIntoMap(values, log.Data) != nil {
			return nil, err
		}
		timestamps[i] = values["block"].(*big.Int).Uint64()
	}
	return &timestamps, nil
}

// Returns an array of members who submitted prices since fromBlock
func (c *NetworkManager) GetLatestPricesSubmissions(fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) ([]common.Address, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*c.networkPrices.Address}
	topicFilter := [][]common.Hash{{c.networkPrices.ABI.Events["PricesSubmitted"].ID}}

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
