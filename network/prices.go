package network

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	// Contract names
	NetworkPrices_ContractName string = "rocketNetworkPrices"

	// Calls
	networkPrices_getPricesBlock           string = "getPricesBlock"
	networkPrices_getRPLPrice              string = "getRPLPrice"
	networkPrices_getLatestReportableBlock string = "getLatestReportableBlock"

	// Transactions
	networkPrices_submitPrices string = "submitPrices"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNetworkPrices
type NetworkPrices struct {
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// Multicall details for network prices
type NetworkPricesDetails struct {
	// Raw parameters
	PricesBlockRaw                 *big.Int `json:"pricesBlockRaw"`
	RplPriceRaw                    *big.Int `json:"rplPriceRaw"`
	LatestReportablePricesBlockRaw *big.Int `json:"latestReportablePricesBlockRaw"`

	// Formatted parameters
	PricesBlock                 uint64  `json:"pricesBlock"`
	RplPrice                    float64 `json:"rplPrice"`
	LatestReportablePricesBlock uint64  `json:"latestReportablePricesBlock"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NetworkPrices contract binding
func NewNetworkPrices(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*NetworkPrices, error) {
	// Create the contract
	contract, err := rp.GetContract(NetworkPrices_ContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting network prices contract: %w", err)
	}

	return &NetworkPrices{
		rp:       rp,
		contract: contract,
	}, nil
}

// ===================
// === Raw Getters ===
// ===================

// Get the block number which network prices are current for
func (c *NetworkPrices) GetPricesBlockRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkPrices_getPricesBlock)
}

// Get the current network RPL price in ETH
func (c *NetworkPrices) GetRplPriceRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkPrices_getRPLPrice)
}

// Returns the latest block number that oracles should be reporting prices for
func (c *NetworkPrices) GetLatestReportablePricesBlockRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, networkPrices_getLatestReportableBlock)
}

// =========================
// === Formatted Getters ===
// =========================

// Get the block number which network prices are current for
func (c *NetworkPrices) GetPricesBlock(opts *bind.CallOpts) (uint64, error) {
	raw, err := c.GetPricesBlockRaw(opts)
	if err != nil {
		return 0, err
	}
	return raw.Uint64(), nil
}

// Get the current network RPL price in ETH
func (c *NetworkPrices) GetRplPrice(opts *bind.CallOpts) (float64, error) {
	raw, err := c.GetRplPriceRaw(opts)
	if err != nil {
		return 0, err
	}
	return eth.WeiToEth(raw), nil
}

// Returns the latest block number that oracles should be reporting prices for
func (c *NetworkPrices) GetLatestReportablePricesBlock(opts *bind.CallOpts) (uint64, error) {
	raw, err := c.GetLatestReportablePricesBlockRaw(opts)
	if err != nil {
		return 0, err
	}
	return raw.Uint64(), nil
}

// ====================
// === Transactions ===
// ====================

// Get info for network price submission
func (c *NetworkPrices) GetSubmitPricesInfo(block uint64, rplPrice *big.Int, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, networkPrices_submitPrices, opts, big.NewInt(int64(block)), rplPrice)
}

// =================
// === Multicall ===
// =================

// Add queries to a multicall batcher
func (c *NetworkPrices) AddMulticallQueries(mc *multicall.MultiCaller, details *NetworkPricesDetails) {
	mc.AddCall(c.contract, &details.PricesBlockRaw, networkPrices_getPricesBlock)
	mc.AddCall(c.contract, &details.RplPriceRaw, networkPrices_getRPLPrice)
	mc.AddCall(c.contract, &details.LatestReportablePricesBlockRaw, networkPrices_getLatestReportableBlock)
}

// Postprocess the multicalled data to get the formatted parameters
func (c *NetworkPrices) PostprocessAfterMulticall(details *NetworkPricesDetails) {
	details.PricesBlock = details.PricesBlockRaw.Uint64()
	details.RplPrice = eth.WeiToEth(details.RplPriceRaw)
	details.LatestReportablePricesBlock = details.LatestReportablePricesBlockRaw.Uint64()
}
