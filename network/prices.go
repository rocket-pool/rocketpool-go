package network

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNetworkPrices
type NetworkPrices struct {
	Details  NetworkPricesDetails
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// Details for network prices
type NetworkPricesDetails struct {
	PricesBlock                 rocketpool.Parameter[uint64]  `json:"pricesBlock"`
	RplPrice                    rocketpool.Parameter[float64] `json:"rplPrice"`
	LatestReportablePricesBlock rocketpool.Parameter[uint64]  `json:"latestReportablePricesBlock"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new NetworkPrices contract binding
func NewNetworkPrices(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*NetworkPrices, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketNetworkPrices", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting network prices contract: %w", err)
	}

	return &NetworkPrices{
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the block number which network prices are current for
func (c *NetworkPrices) GetPricesBlock(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.PricesBlock.RawValue, "getPricesBlock")
}

// Get the current network RPL price in ETH
func (c *NetworkPrices) GetRplPrice(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.RplPrice.RawValue, "getRPLPrice")
}

// Returns the latest block number that oracles should be reporting prices for
func (c *NetworkPrices) GetLatestReportablePricesBlock(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.LatestReportablePricesBlock.RawValue, "getLatestReportableBlock")
}

// Get all basic details
func (c *NetworkPrices) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetPricesBlock(mc)
	c.GetRplPrice(mc)
	c.GetLatestReportablePricesBlock(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for network price submission
func (c *NetworkPrices) SubmitPrices(block uint64, rplPrice *big.Int, opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, "submitPrices", opts, big.NewInt(int64(block)), rplPrice)
}
