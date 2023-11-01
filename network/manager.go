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
	// The block number which network balances are current for
	BalancesBlock *core.FormattedUint256Field[uint64]

	// The current network total ETH balance
	TotalEthBalance *core.SimpleField[*big.Int]

	// The current network staking ETH balance
	StakingEthBalance *core.SimpleField[*big.Int]

	// The current network total rETH supply
	TotalRethSupply *core.SimpleField[*big.Int]

	// The current network ETH utilization rate
	EthUtilizationRate *core.FormattedUint256Field[float64]

	// The latest block number that oracles should be reporting balances for
	LatestReportableBalancesBlock *core.FormattedUint256Field[uint64]

	// The current network node demand in ETH
	NodeDemand *core.SimpleField[*big.Int]

	// The current network node commission rate
	NodeFee *core.FormattedUint256Field[float64]

	// The block number which network prices are current for
	PricesBlock *core.FormattedUint256Field[uint64]

	// The current network RPL price in ETH
	RplPrice *core.FormattedUint256Field[float64]

	// The latest block number that oracles should be reporting prices for
	LatestReportablePricesBlock *core.FormattedUint256Field[uint64]

	// === Internal fields ===
	rp               *rocketpool.RocketPool
	networkBalances  *core.Contract
	networkFees      *core.Contract
	networkPenalties *core.Contract
	networkPrices    *core.Contract
	networkVoting    *core.Contract
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
	networkVoting, err := rp.GetContract(rocketpool.ContractName_RocketNetworkVoting)
	if err != nil {
		return nil, fmt.Errorf("error getting network voting binding: %w", err)
	}

	return &NetworkManager{
		// NetworkBalances
		BalancesBlock:                 core.NewFormattedUint256Field[uint64](networkBalances, "getBalancesBlock"),
		TotalEthBalance:               core.NewSimpleField[*big.Int](networkBalances, "getTotalETHBalance"),
		StakingEthBalance:             core.NewSimpleField[*big.Int](networkBalances, "getStakingETHBalance"),
		TotalRethSupply:               core.NewSimpleField[*big.Int](networkBalances, "getTotalRETHSupply"),
		EthUtilizationRate:            core.NewFormattedUint256Field[float64](networkBalances, "getETHUtilizationRate"),
		LatestReportableBalancesBlock: core.NewFormattedUint256Field[uint64](networkBalances, "getLatestReportableBlock"),

		// NetworkFees
		NodeDemand: core.NewSimpleField[*big.Int](networkFees, "getNodeDemand"),
		NodeFee:    core.NewFormattedUint256Field[float64](networkFees, "getNodeFee"),

		// NetworkPrices
		PricesBlock:                 core.NewFormattedUint256Field[uint64](networkPrices, "getPricesBlock"),
		RplPrice:                    core.NewFormattedUint256Field[float64](networkPrices, "getRPLPrice"),
		LatestReportablePricesBlock: core.NewFormattedUint256Field[uint64](networkPrices, "getLatestReportableBlock"),

		rp:               rp,
		networkBalances:  networkBalances,
		networkFees:      networkFees,
		networkPenalties: networkPenalties,
		networkPrices:    networkPrices,
		networkVoting:    networkVoting,
	}, nil
}

// =============
// === Calls ===
// =============

// === NetworkFees ===

// Get the network node fee for a node demand value
func (c *NetworkManager) GetNodeFeeByDemand(mc *batch.MultiCaller, out **big.Int, demand *big.Int) {
	core.AddCall(mc, c.networkFees, out, "getNodeFeeByDemand", demand)
}

// === NetworkVoting ===

// Get the number of nodes that were present in the network at the provided block
func (c *NetworkManager) GetVotingNodeCountAtBlock(mc *batch.MultiCaller, out **big.Int, blockNumber uint32) {
	core.AddCall(mc, c.networkVoting, out, "getNodeCount", blockNumber)
}

// ====================
// === Transactions ===
// ====================

// === NetworkBalances ===

// Get info for network balance submission
func (c *NetworkManager) SubmitBalances(block uint64, totalEth, stakingEth, rethSupply *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.networkBalances, "submitBalances", opts, block, totalEth, stakingEth, rethSupply)
}

// === NetworkPenalties ===

// Get info for minipool penalty submission
func (c *NetworkManager) SubmitPenalty(minipoolAddress common.Address, block *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.networkPenalties, "submitPenalty", opts, minipoolAddress, block)
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
	topicFilter := [][]common.Hash{{c.networkBalances.ABI.Events["BalancesSubmitted"].ID}, {common.BytesToHash(nodeAddress[:])}}

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
	topicFilter := [][]common.Hash{{c.networkPrices.ABI.Events["PricesSubmitted"].ID}, {common.BytesToHash(nodeAddress[:])}}

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
