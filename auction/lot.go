package auction

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	// Calls
	auctionManager_getLotExists              string = "getLotExists"
	auctionManager_getLotStartBlock          string = "getLotStartBlock"
	auctionManager_getLotEndBlock            string = "getLotEndBlock"
	auctionManager_getLotStartPrice          string = "getLotStartPrice"
	auctionManager_getLotReservePrice        string = "getLotReservePrice"
	auctionManager_getLotTotalRPLAmount      string = "getLotTotalRPLAmount"
	auctionManager_getLotTotalBidAmount      string = "getLotTotalBidAmount"
	auctionManager_getLotRPLRecovered        string = "getLotRPLRecovered"
	auctionManager_getLotPriceAtCurrentBlock string = "getLotPriceAtCurrentBlock"
	auctionManager_getLotPriceByTotalBids    string = "getLotPriceByTotalBids"
	auctionManager_getLotCurrentPrice        string = "getLotCurrentPrice"
	auctionManager_getLotClaimedRPLAmount    string = "getLotClaimedRPLAmount"
	auctionManager_getLotRemainingRPLAmount  string = "getLotRemainingRPLAmount"
	auctionManager_getLotIsCleared           string = "getLotIsCleared"
	auctionManager_getLotPriceAtBlock        string = "getLotPriceAtBlock"
	auctionManager_getLotAddressBidAmount    string = "getLotAddressBidAmount"

	// Transactions
	auctionManager_placeBid            string = "placeBid"
	auctionManager_claimBid            string = "claimBid"
	auctionManager_recoverUnclaimedRPL string = "recoverUnclaimedRPL"
)

// ===============
// === Structs ===
// ===============

// Binding for auction lots
type AuctionLot struct {
	Details AuctionLotDetails
	index   *big.Int
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
		Details: AuctionLotDetails{},
		index:   big.NewInt(int64(index)),
		mgr:     mgr,
	}
}

// ===============
// === Getters ===
// ===============

// Get the lot's index
func (c *AuctionLot) GetLotIndex() uint64 {
	return c.index.Uint64()
}

// Check whether or not the lot exists
func (c *AuctionLot) GetLotExists(opts *bind.CallOpts) (bool, error) {
	return rocketpool.Call[bool](c.mgr.contract, opts, auctionManager_getLotExists, c.index)
}

// Get the lot's start block
func (c *AuctionLot) GetLotStartBlock(opts *bind.CallOpts) (rocketpool.Parameter[uint64], error) {
	return rocketpool.CallForParameter[uint64](c.mgr.contract, opts, auctionManager_getLotStartBlock, c.index)
}

// Get the lot's end block
func (c *AuctionLot) GetLotEndBlock(opts *bind.CallOpts) (rocketpool.Parameter[uint64], error) {
	return rocketpool.CallForParameter[uint64](c.mgr.contract, opts, auctionManager_getLotEndBlock, c.index)
}

// Get the lot's starting price
func (c *AuctionLot) GetLotStartPrice(opts *bind.CallOpts) (rocketpool.Parameter[float64], error) {
	return rocketpool.CallForParameter[float64](c.mgr.contract, opts, auctionManager_getLotStartPrice, c.index)
}

// Get the lot's reserve price
func (c *AuctionLot) GetLotReservePrice(opts *bind.CallOpts) (rocketpool.Parameter[float64], error) {
	return rocketpool.CallForParameter[float64](c.mgr.contract, opts, auctionManager_getLotReservePrice, c.index)
}

// Get the lot's total RPL
func (c *AuctionLot) GetLotTotalRPLAmount(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotTotalRPLAmount, c.index)
}

// Get the lot's total bid amount
func (c *AuctionLot) GetLotTotalBidAmount(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotTotalBidAmount, c.index)
}

// Check whether RPL has been recovered by the lot
func (c *AuctionLot) GetLotRPLRecovered(opts *bind.CallOpts) (bool, error) {
	return rocketpool.Call[bool](c.mgr.contract, opts, auctionManager_getLotRPLRecovered, c.index)
}

// Get the price of the lot in RPL/ETH at the given block
func (c *AuctionLot) GetLotPriceAtCurrentBlock(opts *bind.CallOpts) (rocketpool.Parameter[float64], error) {
	return rocketpool.CallForParameter[float64](c.mgr.contract, opts, auctionManager_getLotPriceAtCurrentBlock, c.index)
}

// Get the price of the lot by the total bids
func (c *AuctionLot) GetLotPriceByTotalBids(opts *bind.CallOpts) (rocketpool.Parameter[float64], error) {
	return rocketpool.CallForParameter[float64](c.mgr.contract, opts, auctionManager_getLotPriceByTotalBids, c.index)
}

// Get the price of the lot at the current block
func (c *AuctionLot) GetLotCurrentPrice(opts *bind.CallOpts) (rocketpool.Parameter[float64], error) {
	return rocketpool.CallForParameter[float64](c.mgr.contract, opts, auctionManager_getLotCurrentPrice, c.index)
}

// Get the amount of RPL claimed for the lot
func (c *AuctionLot) GetLotClaimedRPLAmount(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotClaimedRPLAmount, c.index)
}

// Get the amount of RPL remaining for the lot
func (c *AuctionLot) GetLotRemainingRPLAmount(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotRemainingRPLAmount, c.index)
}

// Check if the lot has been cleared already
func (c *AuctionLot) GetLotIsCleared(opts *bind.CallOpts) (bool, error) {
	return rocketpool.Call[bool](c.mgr.contract, opts, auctionManager_getLotIsCleared, c.index)
}

// Get the price of the lot at the given block
func (c *AuctionLot) GetLotPriceAtBlockRaw(blockNumber uint64, opts *bind.CallOpts) (rocketpool.Parameter[float64], error) {
	return rocketpool.CallForParameter[float64](c.mgr.contract, opts, auctionManager_getLotPriceAtBlock, c.index, big.NewInt(int64(blockNumber)))
}

// Get the ETH amount bid on the lot by an address
func (c *AuctionLot) GetLotAddressBidAmount(bidder common.Address, opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotAddressBidAmount, c.index, bidder)
}

// ====================
// === Transactions ===
// ====================

// Get info for placing a bid on a lot
func (c *AuctionLot) PlaceBid(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.mgr.contract, auctionManager_placeBid, opts, c.index)
}

// Get info for claiming RPL from a lot that was bid on
func (c *AuctionLot) ClaimBid(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.mgr.contract, auctionManager_claimBid, opts, c.index)
}

// Get info for recovering unclaimed RPL from a lot
func (c *AuctionLot) RecoverUnclaimedRpl(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.mgr.contract, auctionManager_recoverUnclaimedRPL, opts, c.index)
}

// =================
// === Multicall ===
// =================

// Check whether or not the lot exists
func (c *AuctionLot) GetLotExistsMC(mc *multicall.MultiCaller) {
	mc.AddCall(c.mgr.contract, &c.Details.Exists, auctionManager_getLotExists)
}

// Add queries to a multicall batcher
func (c *AuctionLot) AddMulticallQueries(mc *multicall.MultiCaller, details *AuctionLotDetails) {
	details.Index.RawValue = c.index
	mc.AddCall(c.mgr.contract, &details.Exists, auctionManager_getLotExists)
	mc.AddCall(c.mgr.contract, &details.StartBlock.RawValue, auctionManager_getLotStartBlock)
	mc.AddCall(c.mgr.contract, &details.EndBlock.RawValue, auctionManager_getLotEndBlock)
	mc.AddCall(c.mgr.contract, &details.StartPrice.RawValue, auctionManager_getLotStartPrice)
	mc.AddCall(c.mgr.contract, &details.ReservePrice.RawValue, auctionManager_getLotReservePrice)
	mc.AddCall(c.mgr.contract, &details.PriceAtCurrentBlock.RawValue, auctionManager_getLotPriceAtCurrentBlock)
	mc.AddCall(c.mgr.contract, &details.PriceByTotalBids.RawValue, auctionManager_getLotPriceByTotalBids)
	mc.AddCall(c.mgr.contract, &details.CurrentPrice.RawValue, auctionManager_getLotCurrentPrice)
	mc.AddCall(c.mgr.contract, &details.TotalRplAmount, auctionManager_getLotTotalRPLAmount)
	mc.AddCall(c.mgr.contract, &details.ClaimedRplAmount, auctionManager_getLotClaimedRPLAmount)
	mc.AddCall(c.mgr.contract, &details.RemainingRplAmount, auctionManager_getLotRemainingRPLAmount)
	mc.AddCall(c.mgr.contract, &details.TotalBidAmount, auctionManager_getLotTotalBidAmount)
	mc.AddCall(c.mgr.contract, &details.AddressBidAmount, auctionManager_getLotStartBlock)
	mc.AddCall(c.mgr.contract, &details.Cleared, auctionManager_getLotIsCleared)
	mc.AddCall(c.mgr.contract, &details.RplRecovered, auctionManager_getLotRPLRecovered)
}

// Add a query to the amount bid by the given address to a multicall batcher
func (c *AuctionLot) AddBidAmountToMulticallQuery(mc *multicall.MultiCaller, details *AuctionLotDetails, bidder common.Address) {
	mc.AddCall(c.mgr.contract, &details.AddressBidAmount, auctionManager_getLotAddressBidAmount, bidder)
}
