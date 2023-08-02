package auction

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for auction lots
type AuctionLot struct {
	Index core.Parameter[uint64] `json:"index"`
	mgr   *core.Contract
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
	Cleared             bool                    `json:"cleared"`
	RplRecovered        bool                    `json:"rplRecovered"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new AuctionLot instance
func NewAuctionLot(rp *rocketpool.RocketPool, index uint64) (*AuctionLot, error) {
	mgr, err := rp.GetContract(rocketpool.ContractName_RocketAuctionManager)
	if err != nil {
		return nil, fmt.Errorf("error getting auction manager contract: %w", err)
	}

	return &AuctionLot{
		Index: core.Parameter[uint64]{
			RawValue: big.NewInt(int64(index)),
		},
		mgr: mgr,
	}, nil
}

// =============
// === Calls ===
// =============

// Check whether or not the lot exists
func (c *AuctionLot) GetLotExists(mc *multicall.MultiCaller, out *bool) {
	multicall.AddCall(mc, c.mgr, out, "getLotExists", c.Index.RawValue)
}

// Get the lot's start block
func (c *AuctionLot) GetLotStartBlock(mc *multicall.MultiCaller, out *core.Parameter[uint64]) {
	multicall.AddCall(mc, c.mgr, &out.RawValue, "getLotStartBlock", c.Index.RawValue)
}

// Get the lot's end block
func (c *AuctionLot) GetLotEndBlock(mc *multicall.MultiCaller, out *core.Parameter[uint64]) {
	multicall.AddCall(mc, c.mgr, &out.RawValue, "getLotEndBlock", c.Index.RawValue)
}

// Get the lot's starting price
func (c *AuctionLot) GetLotStartPrice(mc *multicall.MultiCaller, out *core.Parameter[float64]) {
	multicall.AddCall(mc, c.mgr, &out.RawValue, "getLotStartPrice", c.Index.RawValue)
}

// Get the lot's reserve price
func (c *AuctionLot) GetLotReservePrice(mc *multicall.MultiCaller, out *core.Parameter[float64]) {
	multicall.AddCall(mc, c.mgr, &out.RawValue, "getLotReservePrice", c.Index.RawValue)
}

// Get the price of the lot in RPL/ETH at the given block
func (c *AuctionLot) GetLotPriceAtCurrentBlock(mc *multicall.MultiCaller, out *core.Parameter[float64]) {
	multicall.AddCall(mc, c.mgr, &out.RawValue, "getLotPriceAtCurrentBlock", c.Index.RawValue)
}

// Get the price of the lot by the total bids
func (c *AuctionLot) GetLotPriceByTotalBids(mc *multicall.MultiCaller, out *core.Parameter[float64]) {
	multicall.AddCall(mc, c.mgr, &out.RawValue, "getLotPriceByTotalBids", c.Index.RawValue)
}

// Get the price of the lot at the current block
func (c *AuctionLot) GetLotCurrentPrice(mc *multicall.MultiCaller, out *core.Parameter[float64]) {
	multicall.AddCall(mc, c.mgr, &out.RawValue, "getLotCurrentPrice", c.Index.RawValue)
}

// Get the lot's total RPL
func (c *AuctionLot) GetLotTotalRplAmount(mc *multicall.MultiCaller, out **big.Int) {
	multicall.AddCall(mc, c.mgr, out, "getLotTotalRPLAmount", c.Index.RawValue)
}

// Get the amount of RPL claimed for the lot
func (c *AuctionLot) GetLotClaimedRplAmount(mc *multicall.MultiCaller, out **big.Int) {
	multicall.AddCall(mc, c.mgr, out, "getLotClaimedRPLAmount", c.Index.RawValue)
}

// Get the amount of RPL remaining for the lot
func (c *AuctionLot) GetLotRemainingRplAmount(mc *multicall.MultiCaller, out **big.Int) {
	multicall.AddCall(mc, c.mgr, out, "getLotRemainingRPLAmount", c.Index.RawValue)
}

// Get the lot's total bid amount
func (c *AuctionLot) GetLotTotalBidAmount(mc *multicall.MultiCaller, out **big.Int) {
	multicall.AddCall(mc, c.mgr, out, "getLotTotalBidAmount", c.Index.RawValue)
}

// Check if the lot has been cleared already
func (c *AuctionLot) GetLotIsCleared(mc *multicall.MultiCaller, out *bool) {
	multicall.AddCall(mc, c.mgr, out, "getLotIsCleared", c.Index.RawValue)
}

// Check whether RPL has been recovered by the lot
func (c *AuctionLot) GetLotRplRecovered(mc *multicall.MultiCaller, out *bool) {
	multicall.AddCall(mc, c.mgr, out, "getLotRPLRecovered", c.Index.RawValue)
}

// Get the price of the lot at the given block
func (c *AuctionLot) GetLotPriceAtBlock(mc *multicall.MultiCaller, out *core.Parameter[float64], blockNumber uint64) {
	multicall.AddCall(mc, c.mgr, &out.RawValue, "getLotPriceAtBlock", c.Index.RawValue, big.NewInt(int64(blockNumber)))
}

// Get the ETH amount bid on the lot by an address
func (c *AuctionLot) GetLotAddressBidAmount(mc *multicall.MultiCaller, out **big.Int, bidder common.Address) {
	multicall.AddCall(mc, c.mgr, out, "getLotAddressBidAmount", c.Index.RawValue, bidder)
}

// Get all basic details
func (c *AuctionLot) GetAllDetails(mc *multicall.MultiCaller) (*AuctionLotDetails, error) {
	details := &AuctionLotDetails{}
	err := utils.GetAllDetails(c, details, mc)
	if err != nil {
		return nil, fmt.Errorf("error getting details: %w", err)
	}
	return details, nil
}

// ====================
// === Transactions ===
// ====================

// Get info for placing a bid on a lot
func (c *AuctionLot) PlaceBid(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.mgr, "placeBid", opts, c.Index.RawValue)
}

// Get info for claiming RPL from a lot that was bid on
func (c *AuctionLot) ClaimBid(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.mgr, "claimBid", opts, c.Index.RawValue)
}

// Get info for recovering unclaimed RPL from a lot
func (c *AuctionLot) RecoverUnclaimedRpl(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.mgr, "recoverUnclaimedRPL", opts, c.Index.RawValue)
}
