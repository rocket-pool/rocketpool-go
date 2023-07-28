package auction

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for auction lots
type AuctionLot struct {
	Details AuctionLotDetails
	mgr     *AuctionManager
}

// Details for auction lots
type AuctionLotDetails struct {
	Index               rocketpool.Parameter[uint64]  `json:"index"`
	Exists              bool                          `json:"exists"`
	StartBlock          rocketpool.Parameter[uint64]  `json:"startBlock"`
	EndBlock            rocketpool.Parameter[uint64]  `json:"endBlock"`
	StartPrice          rocketpool.Parameter[float64] `json:"startPrice"`
	ReservePrice        rocketpool.Parameter[float64] `json:"reservePrice"`
	PriceAtCurrentBlock rocketpool.Parameter[float64] `json:"priceAtCurrentBlock"`
	PriceByTotalBids    rocketpool.Parameter[float64] `json:"priceByTotalBids"`
	CurrentPrice        rocketpool.Parameter[float64] `json:"currentPrice"`
	TotalRplAmount      *big.Int                      `json:"totalRplAmount"`
	ClaimedRplAmount    *big.Int                      `json:"claimedRplAmount"`
	RemainingRplAmount  *big.Int                      `json:"remainingRplAmount"`
	TotalBidAmount      *big.Int                      `json:"totalBidAmount"`
	AddressBidAmount    *big.Int                      `json:"addressBidAmount"`
	Cleared             bool                          `json:"cleared"`
	RplRecovered        bool                          `json:"rplRecovered"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new AuctionLot instance
func NewAuctionLot(mgr *AuctionManager, index uint64) *AuctionLot {
	return &AuctionLot{
		Details: AuctionLotDetails{
			Index: rocketpool.Parameter[uint64]{
				RawValue: big.NewInt(int64(index)),
			},
		},
		mgr: mgr,
	}
}

// =============
// === Calls ===
// =============

// Check whether or not the lot exists
func (c *AuctionLot) GetLotExists(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.Exists, "getLotExists", c.Details.Index.RawValue)
}

// Get the lot's start block
func (c *AuctionLot) GetLotStartBlock(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.StartBlock.RawValue, "getLotStartBlock", c.Details.Index.RawValue)
}

// Get the lot's end block
func (c *AuctionLot) GetLotEndBlock(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.EndBlock.RawValue, "getLotEndBlock", c.Details.Index.RawValue)
}

// Get the lot's starting price
func (c *AuctionLot) GetLotStartPrice(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.StartPrice.RawValue, "getLotStartPrice", c.Details.Index.RawValue)
}

// Get the lot's reserve price
func (c *AuctionLot) GetLotReservePrice(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.ReservePrice.RawValue, "getLotReservePrice", c.Details.Index.RawValue)
}

// Get the lot's total RPL
func (c *AuctionLot) GetLotTotalRPLAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.TotalRplAmount, "getLotTotalRPLAmount", c.Details.Index.RawValue)
}

// Get the lot's total bid amount
func (c *AuctionLot) GetLotTotalBidAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.TotalBidAmount, "getLotTotalBidAmount", c.Details.Index.RawValue)
}

// Check whether RPL has been recovered by the lot
func (c *AuctionLot) GetLotRPLRecovered(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.RplRecovered, "getLotRPLRecovered", c.Details.Index.RawValue)
}

// Get the price of the lot in RPL/ETH at the given block
func (c *AuctionLot) GetLotPriceAtCurrentBlock(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.PriceAtCurrentBlock.RawValue, "getLotPriceAtCurrentBlock", c.Details.Index.RawValue)
}

// Get the price of the lot by the total bids
func (c *AuctionLot) GetLotPriceByTotalBids(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.PriceByTotalBids.RawValue, "getLotPriceByTotalBids", c.Details.Index.RawValue)
}

// Get the price of the lot at the current block
func (c *AuctionLot) GetLotCurrentPrice(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.CurrentPrice.RawValue, "getLotCurrentPrice", c.Details.Index.RawValue)
}

// Get the amount of RPL claimed for the lot
func (c *AuctionLot) GetLotClaimedRPLAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.ClaimedRplAmount, "getLotClaimedRPLAmount", c.Details.Index.RawValue)
}

// Get the amount of RPL remaining for the lot
func (c *AuctionLot) GetLotRemainingRPLAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.RemainingRplAmount, "getLotRemainingRPLAmount", c.Details.Index.RawValue)
}

// Check if the lot has been cleared already
func (c *AuctionLot) GetLotIsCleared(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.Cleared, "getLotIsCleared", c.Details.Index.RawValue)
}

// Get the price of the lot at the given block
func (c *AuctionLot) GetLotPriceAtBlock(mc *multicall.MultiCaller, blockNumber uint64, price_Out *rocketpool.Parameter[float64]) {
	*price_Out = rocketpool.Parameter[float64]{}
	multicall.AddCall(mc, c.mgr.contract, &price_Out.RawValue, "getLotPriceAtBlock", c.Details.Index.RawValue, big.NewInt(int64(blockNumber)))
}

// Get the ETH amount bid on the lot by an address
func (c *AuctionLot) GetLotAddressBidAmount(mc *multicall.MultiCaller, bidder common.Address, bidAmount_Out **big.Int) {
	multicall.AddCall(mc, c.mgr.contract, bidAmount_Out, "getLotAddressBidAmount", c.Details.Index.RawValue, bidder)
}

// Get all basic details
func (c *AuctionLot) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetLotExists(mc)
	c.GetLotStartBlock(mc)
	c.GetLotEndBlock(mc)
	c.GetLotStartPrice(mc)
	c.GetLotReservePrice(mc)
	c.GetLotTotalRPLAmount(mc)
	c.GetLotTotalBidAmount(mc)
	c.GetLotRPLRecovered(mc)
	c.GetLotPriceAtCurrentBlock(mc)
	c.GetLotPriceByTotalBids(mc)
	c.GetLotCurrentPrice(mc)
	c.GetLotClaimedRPLAmount(mc)
	c.GetLotRemainingRPLAmount(mc)
	c.GetLotIsCleared(mc)
}

// Get all basic details and the amount bid by the given address
func (c *AuctionLot) GetAllDetailsWithBidAmount(mc *multicall.MultiCaller, bidder common.Address) {
	c.GetAllDetails(mc)
	c.GetLotAddressBidAmount(mc, bidder, &c.Details.AddressBidAmount)
}

// ====================
// === Transactions ===
// ====================

// Get info for placing a bid on a lot
func (c *AuctionLot) PlaceBid(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.mgr.contract, "placeBid", opts, c.Details.Index.RawValue)
}

// Get info for claiming RPL from a lot that was bid on
func (c *AuctionLot) ClaimBid(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.mgr.contract, "claimBid", opts, c.Details.Index.RawValue)
}

// Get info for recovering unclaimed RPL from a lot
func (c *AuctionLot) RecoverUnclaimedRpl(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.mgr.contract, "recoverUnclaimedRPL", opts, c.Details.Index.RawValue)
}
