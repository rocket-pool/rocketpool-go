package auction

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for auction lots
type AuctionLot struct {
	*AuctionLotDetails
	am *core.Contract
}

// Details for auction lots
type AuctionLotDetails struct {
	Index               core.Parameter[uint64]  `json:"index"`
	Exists              bool                    `json:"exists"`
	StartBlock          core.Parameter[uint64]  `json:"startBlock"`
	EndBlock            core.Parameter[uint64]  `json:"endBlock"`
	StartPrice          core.Parameter[float64] `json:"startPrice"`
	ReservePrice        core.Parameter[float64] `json:"reservePrice"`
	PriceAtCurrentBlock core.Parameter[float64] `json:"priceAtCurrentBlock"`
	PriceByTotalBids    core.Parameter[float64] `json:"priceByTotalBids"`
	CurrentPrice        core.Parameter[float64] `json:"currentPrice"`
	TotalRplAmount      *big.Int                `json:"totalRplAmount"`
	ClaimedRplAmount    *big.Int                `json:"claimedRplAmount"`
	RemainingRplAmount  *big.Int                `json:"remainingRplAmount"`
	TotalBidAmount      *big.Int                `json:"totalBidAmount"`
	IsCleared           bool                    `json:"cleared"`
	RplRecovered        bool                    `json:"rplRecovered"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new AuctionLot instance
func NewAuctionLot(rp *rocketpool.RocketPool, index uint64) (*AuctionLot, error) {
	// Create the contract
	am, err := rp.GetContract(rocketpool.ContractName_RocketAuctionManager)
	if err != nil {
		return nil, fmt.Errorf("error getting auction manager contract: %w", err)
	}

	return &AuctionLot{
		AuctionLotDetails: &AuctionLotDetails{
			Index: core.Parameter[uint64]{
				RawValue: big.NewInt(int64(index)),
			},
		},
		am: am,
	}, nil
}

// =============
// === Calls ===
// =============

// Check whether or not the lot exists
func (c *AuctionLot) GetLotExists(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.Exists, "getLotExists", c.Index.RawValue)
}

// Get the lot's start block
func (c *AuctionLot) GetLotStartBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.StartBlock.RawValue, "getLotStartBlock", c.Index.RawValue)
}

// Get the lot's end block
func (c *AuctionLot) GetLotEndBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.EndBlock.RawValue, "getLotEndBlock", c.Index.RawValue)
}

// Get the lot's starting price
func (c *AuctionLot) GetLotStartPrice(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.StartPrice.RawValue, "getLotStartPrice", c.Index.RawValue)
}

// Get the lot's reserve price
func (c *AuctionLot) GetLotReservePrice(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.ReservePrice.RawValue, "getLotReservePrice", c.Index.RawValue)
}

// Get the price of the lot in RPL/ETH at the given block
func (c *AuctionLot) GetLotPriceAtCurrentBlock(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.PriceAtCurrentBlock.RawValue, "getLotPriceAtCurrentBlock", c.Index.RawValue)
}

// Get the price of the lot by the total bids
func (c *AuctionLot) GetLotPriceByTotalBids(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.PriceByTotalBids.RawValue, "getLotPriceByTotalBids", c.Index.RawValue)
}

// Get the price of the lot at the current block
func (c *AuctionLot) GetLotCurrentPrice(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.CurrentPrice.RawValue, "getLotCurrentPrice", c.Index.RawValue)
}

// Get the lot's total RPL
func (c *AuctionLot) GetLotTotalRplAmount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.TotalRplAmount, "getLotTotalRPLAmount", c.Index.RawValue)
}

// Get the amount of RPL claimed for the lot
func (c *AuctionLot) GetLotClaimedRplAmount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.ClaimedRplAmount, "getLotClaimedRPLAmount", c.Index.RawValue)
}

// Get the amount of RPL remaining for the lot
func (c *AuctionLot) GetLotRemainingRplAmount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.RemainingRplAmount, "getLotRemainingRPLAmount", c.Index.RawValue)
}

// Get the lot's total bid amount
func (c *AuctionLot) GetLotTotalBidAmount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.TotalBidAmount, "getLotTotalBidAmount", c.Index.RawValue)
}

// Check if the lot has been cleared already
func (c *AuctionLot) GetLotIsCleared(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.IsCleared, "getLotIsCleared", c.Index.RawValue)
}

// Check whether RPL has been recovered by the lot
func (c *AuctionLot) GetLotRplRecovered(mc *batch.MultiCaller) {
	core.AddCall(mc, c.am, &c.RplRecovered, "getLotRPLRecovered", c.Index.RawValue)
}

// Get the price of the lot at the given block
func (c *AuctionLot) GetLotPriceAtBlock(mc *batch.MultiCaller, out *core.Parameter[float64], blockNumber uint64) {
	core.AddCall(mc, c.am, &out.RawValue, "getLotPriceAtBlock", c.Index.RawValue, big.NewInt(int64(blockNumber)))
}

// Get the ETH amount bid on the lot by an address
func (c *AuctionLot) GetLotAddressBidAmount(mc *batch.MultiCaller, out **big.Int, bidder common.Address) {
	core.AddCall(mc, c.am, out, "getLotAddressBidAmount", c.Index.RawValue, bidder)
}

// Get all basic details
func (c *AuctionLot) GetAllDetails(mc *batch.MultiCaller) {
	c.GetLotExists(mc)
	c.GetLotStartBlock(mc)
	c.GetLotEndBlock(mc)
	c.GetLotStartPrice(mc)
	c.GetLotReservePrice(mc)
	c.GetLotPriceAtCurrentBlock(mc)
	c.GetLotPriceByTotalBids(mc)
	c.GetLotCurrentPrice(mc)
	c.GetLotTotalRplAmount(mc)
	c.GetLotClaimedRplAmount(mc)
	c.GetLotRemainingRplAmount(mc)
	c.GetLotTotalBidAmount(mc)
	c.GetLotIsCleared(mc)
	c.GetLotRplRecovered(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for placing a bid on a lot
func (c *AuctionLot) PlaceBid(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.am, "placeBid", opts, c.Index.RawValue)
}

// Get info for claiming RPL from a lot that was bid on
func (c *AuctionLot) ClaimBid(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.am, "claimBid", opts, c.Index.RawValue)
}

// Get info for recovering unclaimed RPL from a lot
func (c *AuctionLot) RecoverUnclaimedRpl(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.am, "recoverUnclaimedRPL", opts, c.Index.RawValue)
}
