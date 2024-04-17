package network

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/eth"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/v2/core"
	"github.com/rocket-pool/rocketpool-go/v2/node"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
	"github.com/rocket-pool/rocketpool-go/v2/types"
	"github.com/rocket-pool/rocketpool-go/v2/utils"
)

const (
	nodeVotingInfoBatchSize int = 250
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
	txMgr            *eth.TransactionManager
}

// Info for a balances updated event
type BalancesUpdated struct {
	BlockNumber    uint64    `json:"blockNumber"`
	SlotTimestamp  time.Time `json:"slotTimestamp"`
	TotalEth       *big.Int  `json:"totalEth"`
	StakingEth     *big.Int  `json:"stakingEth"`
	RethSupply     *big.Int  `json:"rethSupply"`
	BlockTimestamp time.Time `json:"blockTimestamp"`
}

// Info for a balances updated event
type balancesUpdatedRaw struct {
	Block          *big.Int `json:"block"`
	SlotTimestamp  *big.Int `json:"slotTimestamp"`
	TotalEth       *big.Int `json:"totalEth"`
	StakingEth     *big.Int `json:"stakingEth"`
	RethSupply     *big.Int `json:"rethSupply"`
	BlockTimestamp *big.Int `json:"blockTimestamp"`
}

// Info for a price updated event
type PriceUpdated struct {
	BlockNumber   uint64    `json:"blockNumber"`
	SlotTimestamp time.Time `json:"slotTimestamp"`
	RplPrice      *big.Int  `json:"rplPrice"`
	Time          time.Time `json:"time"`
}

// Info for a price updated event
type priceUpdatedRaw struct {
	Block         *big.Int `json:"block"`
	SlotTimestamp *big.Int `json:"slotTimestamp"`
	RplPrice      *big.Int `json:"rplPrice"`
	Time          *big.Int `json:"time"`
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
		BalancesBlock:      core.NewFormattedUint256Field[uint64](networkBalances, "getBalancesBlock"),
		TotalEthBalance:    core.NewSimpleField[*big.Int](networkBalances, "getTotalETHBalance"),
		StakingEthBalance:  core.NewSimpleField[*big.Int](networkBalances, "getStakingETHBalance"),
		TotalRethSupply:    core.NewSimpleField[*big.Int](networkBalances, "getTotalRETHSupply"),
		EthUtilizationRate: core.NewFormattedUint256Field[float64](networkBalances, "getETHUtilizationRate"),

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
		txMgr:            rp.GetTransactionManager(),
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
func (c *NetworkManager) SubmitBalances(block uint64, slotTimestamp uint64, totalEth *big.Int, stakingEth *big.Int, rethSupply *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.networkBalances.Contract, "submitBalances", opts, big.NewInt(int64(block)), big.NewInt(int64(slotTimestamp)), totalEth, stakingEth, rethSupply)
}

// === NetworkPenalties ===

// Get info for minipool penalty submission
func (c *NetworkManager) SubmitPenalty(minipoolAddress common.Address, block *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.networkPenalties.Contract, "submitPenalty", opts, minipoolAddress, block)
}

// === NetworkPrices ===

// Get info for network price submission
func (c *NetworkManager) SubmitPrices(block uint64, slotTimestamp uint64, rplPrice *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.networkPrices.Contract, "submitPrices", opts, big.NewInt(int64(block)), big.NewInt(int64(slotTimestamp)), rplPrice)
}

// =============
// === Utils ===
// =============

// Gets the voting power and delegation info for every node at the specified block using multicall
func (c *NetworkManager) GetNodeVotingInfo(blockNumber uint32, nodeAddresses []common.Address, opts *bind.CallOpts) ([]types.NodeVotingInfo, error) {
	infos := make([]types.NodeVotingInfo, len(nodeAddresses))
	err := c.rp.BatchQuery(len(nodeAddresses), nodeVotingInfoBatchSize, func(mc *batch.MultiCaller, i int) error {
		address := nodeAddresses[i]
		node, err := node.NewNode(c.rp, address)
		if err != nil {
			return fmt.Errorf("error creating node binding for node %s: %w", address.Hex(), err)
		}
		node.GetVotingPowerAtBlock(mc, &infos[i].VotingPower, blockNumber)
		node.GetVotingDelegateAtBlock(mc, &infos[i].Delegate, blockNumber)
		return nil
	}, opts)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

// Returns an array of block numbers for balances submissions the given trusted node has submitted since fromBlock
func (c *NetworkManager) GetBalancesSubmissions(nodeAddress common.Address, fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) (*[]uint64, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{c.networkBalances.Address}
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
	addressFilter := []common.Address{c.networkBalances.Address}
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

// Get the event emitted when the network balances were updated
func (c *NetworkManager) GetBalancesUpdatedEvent(blockNumber uint64, intervalSize *big.Int, balanceAddresses []common.Address, opts *bind.CallOpts) (bool, BalancesUpdated, error) {
	// Create the list of addresses to check
	currentAddress := c.networkBalances.Address
	if balanceAddresses == nil {
		balanceAddresses = []common.Address{currentAddress}
	} else {
		found := false
		for _, address := range balanceAddresses {
			if address == currentAddress {
				found = true
				break
			}
		}
		if !found {
			balanceAddresses = append(balanceAddresses, currentAddress)
		}
	}

	// Construct a filter query for relevant logs
	balancesUpdatedEvent := c.networkBalances.ABI.Events["BalancesUpdated"]
	addressFilter := balanceAddresses
	topicFilter := [][]common.Hash{{balancesUpdatedEvent.ID}}

	// Get the event logs, starting with the target and ending with the chain head
	logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, big.NewInt(int64(blockNumber)), nil, nil)
	if err != nil {
		return false, BalancesUpdated{}, err
	}
	if len(logs) == 0 {
		return false, BalancesUpdated{}, nil
	}

	// Get the list of events
	for _, log := range logs {
		// Get the log info values
		values, err := balancesUpdatedEvent.Inputs.Unpack(log.Data)
		if err != nil {
			return false, BalancesUpdated{}, fmt.Errorf("error unpacking balance updated event data: %w", err)
		}

		// Convert to a native struct
		var eventData balancesUpdatedRaw
		err = balancesUpdatedEvent.Inputs.Copy(&eventData, values)
		if err != nil {
			return false, BalancesUpdated{}, fmt.Errorf("error converting balance updated event data to struct: %w", err)
		}

		// Filter on events that target the block number
		if eventData.Block.Uint64() == blockNumber {
			balancesUpdated := BalancesUpdated{
				BlockNumber:    blockNumber,
				SlotTimestamp:  time.Unix(eventData.SlotTimestamp.Int64(), 0),
				TotalEth:       eventData.TotalEth,
				StakingEth:     eventData.StakingEth,
				RethSupply:     eventData.RethSupply,
				BlockTimestamp: time.Unix(eventData.BlockTimestamp.Int64(), 0),
			}

			return true, balancesUpdated, nil
		}
	}

	return false, BalancesUpdated{}, nil
}

// Returns an array of block numbers for prices submissions the given trusted node has submitted since fromBlock
func (c *NetworkManager) GetPricesSubmissions(nodeAddress common.Address, fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) (*[]uint64, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{c.networkPrices.Address}
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
	addressFilter := []common.Address{c.networkPrices.Address}
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

// Get the event info for a price update
func (c *NetworkManager) GetPriceUpdatedEvent(blockNumber uint64, intervalSize *big.Int, priceAddresses []common.Address, opts *bind.CallOpts) (bool, PriceUpdated, error) {
	// Create the list of addresses to check
	currentAddress := c.networkBalances.Address
	if priceAddresses == nil {
		priceAddresses = []common.Address{currentAddress}
	} else {
		found := false
		for _, address := range priceAddresses {
			if address == currentAddress {
				found = true
				break
			}
		}
		if !found {
			priceAddresses = append(priceAddresses, currentAddress)
		}
	}

	// Construct a filter query for relevant logs
	pricesUpdatedEvent := c.networkPrices.ABI.Events["PricesUpdated"]
	addressFilter := priceAddresses
	topicFilter := [][]common.Hash{{pricesUpdatedEvent.ID}}

	// Get the event logs, starting with the target and ending with the chain head
	logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, big.NewInt(int64(blockNumber)), nil, nil)
	if err != nil {
		return false, PriceUpdated{}, err
	}
	if len(logs) == 0 {
		return false, PriceUpdated{}, nil
	}

	// Get the list of events
	for _, log := range logs {
		// Get the log info values
		values, err := pricesUpdatedEvent.Inputs.Unpack(log.Data)
		if err != nil {
			return false, PriceUpdated{}, fmt.Errorf("error unpacking price updated event data: %w", err)
		}

		// Convert to a native struct
		var eventData priceUpdatedRaw
		err = pricesUpdatedEvent.Inputs.Copy(&eventData, values)
		if err != nil {
			return false, PriceUpdated{}, fmt.Errorf("error converting price updated event data to struct: %w", err)
		}

		fmt.Printf("Found event for block %d on block %d, eventData: %v\n", eventData.Block, log.BlockNumber, eventData)

		// Filter on events that target the block number
		if eventData.Block.Uint64() == blockNumber {
			priceUpdated := PriceUpdated{
				BlockNumber:   eventData.Block.Uint64(),
				SlotTimestamp: time.Unix(eventData.SlotTimestamp.Int64(), 0),
				RplPrice:      eventData.RplPrice,
				Time:          time.Unix(eventData.Time.Int64(), 0),
			}

			return true, priceUpdated, nil
		}
	}

	return false, PriceUpdated{}, nil
}
