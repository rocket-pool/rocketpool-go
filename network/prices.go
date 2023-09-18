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

// Binding for RocketNetworkPrices
type NetworkPrices struct {
	*NetworkPricesDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for network prices
type NetworkPricesDetails struct {
	PricesBlock                 core.Parameter[uint64]  `json:"pricesBlock"`
	RplPrice                    core.Parameter[float64] `json:"rplPrice"`
	LatestReportablePricesBlock core.Parameter[uint64]  `json:"latestReportablePricesBlock"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NetworkPrices contract binding
func NewNetworkPrices(rp *rocketpool.RocketPool) (*NetworkPrices, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketNetworkPrices)
	if err != nil {
		return nil, fmt.Errorf("error getting network prices contract: %w", err)
	}

	return &NetworkPrices{
		NetworkPricesDetails: &NetworkPricesDetails{},
		rp:                   rp,
		contract:             contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the block number which network prices are current for
func (c *NetworkPrices) GetPricesBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.PricesBlock.RawValue, "getPricesBlock")
}

// Get the current network RPL price in ETH
func (c *NetworkPrices) GetRplPrice(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.RplPrice.RawValue, "getRPLPrice")
}

// Returns the latest block number that oracles should be reporting prices for
func (c *NetworkPrices) GetLatestReportablePricesBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.LatestReportablePricesBlock.RawValue, "getLatestReportableBlock")
}

// Get all basic details
func (c *NetworkPrices) GetAllDetails(mc *batch.MultiCaller) {
	c.GetPricesBlock(mc)
	c.GetRplPrice(mc)
	c.GetLatestReportablePricesBlock(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for network price submission
func (c *NetworkPrices) SubmitPrices(block uint64, rplPrice *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "submitPrices", opts, big.NewInt(int64(block)), rplPrice)
}

// =============
// === Utils ===
// =============

// Returns an array of block numbers for prices submissions the given trusted node has submitted since fromBlock
func (c *NetworkPrices) GetPricesSubmissions(nodeAddress common.Address, fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) (*[]uint64, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*c.contract.Address}
	topicFilter := [][]common.Hash{{c.contract.ABI.Events["PricesSubmitted"].ID}, {nodeAddress.Hash()}}

	// Get the event logs
	logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, big.NewInt(int64(fromBlock)), nil, nil)
	if err != nil {
		return nil, err
	}
	timestamps := make([]uint64, len(logs))
	for i, log := range logs {
		values := make(map[string]interface{})
		// Decode the event
		if c.contract.ABI.Events["PricesSubmitted"].Inputs.UnpackIntoMap(values, log.Data) != nil {
			return nil, err
		}
		timestamps[i] = values["block"].(*big.Int).Uint64()
	}
	return &timestamps, nil
}

// Returns an array of members who submitted prices since fromBlock
func (c *NetworkPrices) GetLatestPricesSubmissions(fromBlock uint64, intervalSize *big.Int, opts *bind.CallOpts) ([]common.Address, error) {
	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*c.contract.Address}
	topicFilter := [][]common.Hash{{c.contract.ABI.Events["PricesSubmitted"].ID}}

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
