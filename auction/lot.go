package auction

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
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
	Index   *big.Int
	Details AuctionLotDetails
	mgr     *AuctionManager
}

// Details for auction lots
type AuctionLotDetails struct {
	// Raw parameters
	IndexRaw               *big.Int `json:"indexRaw"`
	Exists                 bool     `json:"exists"`
	StartBlockRaw          *big.Int `json:"startBlockRaw"`
	EndBlockRaw            *big.Int `json:"endBlockRaw"`
	StartPriceRaw          *big.Int `json:"startPriceRaw"`
	ReservePriceRaw        *big.Int `json:"reservePriceRaw"`
	PriceAtCurrentBlockRaw *big.Int `json:"priceAtCurrentBlockRaw"`
	PriceByTotalBidsRaw    *big.Int `json:"priceByTotalBidsRaw"`
	CurrentPriceRaw        *big.Int `json:"currentPriceRaw"`
	TotalRplAmount         *big.Int `json:"totalRplAmount"`
	ClaimedRplAmount       *big.Int `json:"claimedRplAmount"`
	RemainingRplAmount     *big.Int `json:"remainingRplAmount"`
	TotalBidAmount         *big.Int `json:"totalBidAmount"`
	AddressBidAmount       *big.Int `json:"addressBidAmount"`
	Cleared                bool     `json:"cleared"`
	RplRecovered           bool     `json:"rplRecovered"`

	// Formatted parameters
	Index               uint64  `json:"index"`
	StartBlock          uint64  `json:"startBlock"`
	EndBlock            uint64  `json:"endBlock"`
	StartPrice          float64 `json:"startPrice"`
	ReservePrice        float64 `json:"reservePrice"`
	PriceAtCurrentBlock float64 `json:"priceAtCurrentBlock"`
	PriceByTotalBids    float64 `json:"priceByTotalBids"`
	CurrentPrice        float64 `json:"currentPrice"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new AuctionLot instance
func NewAuctionLot(mgr *AuctionManager, index uint64) *AuctionLot {
	return &AuctionLot{
		Index:   big.NewInt(int64(index)),
		Details: AuctionLotDetails{},
		mgr:     mgr,
	}
}

// ===================
// === Raw Getters ===
// ===================

// Check whether or not the lot exists
func (c *AuctionLot) GetLotExists(opts *bind.CallOpts) (bool, error) {
	return rocketpool.Call[bool](c.mgr.contract, opts, auctionManager_getLotExists, c.Index)
}

// Get the lot's start block
func (c *AuctionLot) GetLotStartBlockRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotStartBlock, c.Index)
}

// Get the lot's end block
func (c *AuctionLot) GetLotEndBlockRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotEndBlock, c.Index)
}

// Get the lot's starting price
func (c *AuctionLot) GetLotStartPriceRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotStartPrice, c.Index)
}

// Get the lot's reserve price
func (c *AuctionLot) GetLotReservePriceRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotReservePrice, c.Index)
}

// Get the lot's total RPL
func (c *AuctionLot) GetLotTotalRPLAmount(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotTotalRPLAmount, c.Index)
}

// Get the lot's total bid amount
func (c *AuctionLot) GetLotTotalBidAmount(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotTotalBidAmount, c.Index)
}

// Check whether RPL has been recovered by the lot
func (c *AuctionLot) GetLotRPLRecovered(opts *bind.CallOpts) (bool, error) {
	return rocketpool.Call[bool](c.mgr.contract, opts, auctionManager_getLotRPLRecovered, c.Index)
}

// Get the price of the lot in RPL/ETH at the given block
func (c *AuctionLot) GetLotPriceAtCurrentBlockRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotPriceAtCurrentBlock, c.Index)
}

// Get the price of the lot by the total bids
func (c *AuctionLot) GetLotPriceByTotalBidsRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotPriceByTotalBids, c.Index)
}

// Get the price of the lot at the current block
func (c *AuctionLot) GetLotCurrentPriceRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotCurrentPrice, c.Index)
}

// Get the amount of RPL claimed for the lot
func (c *AuctionLot) GetLotClaimedRPLAmount(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotClaimedRPLAmount, c.Index)
}

// Get the amount of RPL remaining for the lot
func (c *AuctionLot) GetLotRemainingRPLAmount(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotRemainingRPLAmount, c.Index)
}

// Check if the lot has been cleared already
func (c *AuctionLot) GetLotIsCleared(opts *bind.CallOpts) (bool, error) {
	return rocketpool.Call[bool](c.mgr.contract, opts, auctionManager_getLotIsCleared, c.Index)
}

// Get the price of the lot at the given block
func (c *AuctionLot) GetLotPriceAtBlockRaw(blockNumber uint64, opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotPriceAtBlock, c.Index, big.NewInt(int64(blockNumber)))
}

// Get the ETH amount bid on the lot by an address
func (c *AuctionLot) GetLotAddressBidAmount(bidder common.Address, opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, auctionManager_getLotAddressBidAmount, c.Index, bidder)
}

// =========================
// === Formatted Getters ===
// =========================

// Get the lot's start block
func (c *AuctionLot) GetLotStartBlock(opts *bind.CallOpts) (uint64, error) {
	raw, err := c.GetLotStartBlockRaw(opts)
	if err != nil {
		return 0, err
	}
	return raw.Uint64(), nil
}

// Get the lot's end block
func (c *AuctionLot) GetLotEndBlock(opts *bind.CallOpts) (uint64, error) {
	raw, err := c.GetLotEndBlockRaw(opts)
	if err != nil {
		return 0, err
	}
	return raw.Uint64(), nil
}

// Get the lot's starting price
func (c *AuctionLot) GetLotStartPrice(opts *bind.CallOpts) (float64, error) {
	raw, err := c.GetLotStartPriceRaw(opts)
	if err != nil {
		return 0, err
	}
	return eth.WeiToEth(raw), nil
}

// Get the lot's reserve price
func (c *AuctionLot) GetLotReservePrice(opts *bind.CallOpts) (float64, error) {
	raw, err := c.GetLotReservePriceRaw(opts)
	if err != nil {
		return 0, err
	}
	return eth.WeiToEth(raw), nil
}

// Get the price of the lot in RPL/ETH at the given block
func (c *AuctionLot) GetLotPriceAtCurrentBlock(opts *bind.CallOpts) (float64, error) {
	raw, err := c.GetLotPriceAtCurrentBlockRaw(opts)
	if err != nil {
		return 0, err
	}
	return eth.WeiToEth(raw), nil
}

// Get the price of the lot by the total bids
func (c *AuctionLot) GetLotPriceByTotalBids(opts *bind.CallOpts) (float64, error) {
	raw, err := c.GetLotPriceByTotalBidsRaw(opts)
	if err != nil {
		return 0, err
	}
	return eth.WeiToEth(raw), nil
}

// Get the price of the lot at the current block
func (c *AuctionLot) GetLotCurrentPrice(opts *bind.CallOpts) (float64, error) {
	raw, err := c.GetLotCurrentPriceRaw(opts)
	if err != nil {
		return 0, err
	}
	return eth.WeiToEth(raw), nil
}

// Get the price of the lot at the given block
func (c *AuctionLot) GetLotPriceAtBlock(blockNumber uint64, opts *bind.CallOpts) (float64, error) {
	raw, err := c.GetLotPriceAtBlockRaw(blockNumber, opts)
	if err != nil {
		return 0, err
	}
	return eth.WeiToEth(raw), nil
}

// ====================
// === Transactions ===
// ====================

// Get info for placing a bid on a lot
func (c *AuctionLot) PlaceBid(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.mgr.contract, auctionManager_placeBid, opts, c.Index)
}

// Get info for claiming RPL from a lot that was bid on
func (c *AuctionLot) ClaimBid(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.mgr.contract, auctionManager_claimBid, opts, c.Index)
}

// Get info for recovering unclaimed RPL from a lot
func (c *AuctionLot) RecoverUnclaimedRpl(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.mgr.contract, auctionManager_recoverUnclaimedRPL, opts, c.Index)
}

// =================
// === Multicall ===
// =================

// Add queries to a multicall batcher
func (c *AuctionLot) AddMulticallQueries(mc *multicall.MultiCaller, details *AuctionLotDetails) {
	details.IndexRaw = c.Index
	mc.AddCall(c.mgr.contract, &details.Exists, auctionManager_getLotExists)
	mc.AddCall(c.mgr.contract, &details.StartBlockRaw, auctionManager_getLotStartBlock)
	mc.AddCall(c.mgr.contract, &details.EndBlockRaw, auctionManager_getLotEndBlock)
	mc.AddCall(c.mgr.contract, &details.StartPriceRaw, auctionManager_getLotStartPrice)
	mc.AddCall(c.mgr.contract, &details.ReservePriceRaw, auctionManager_getLotReservePrice)
	mc.AddCall(c.mgr.contract, &details.PriceAtCurrentBlockRaw, auctionManager_getLotPriceAtCurrentBlock)
	mc.AddCall(c.mgr.contract, &details.PriceByTotalBidsRaw, auctionManager_getLotPriceByTotalBids)
	mc.AddCall(c.mgr.contract, &details.CurrentPriceRaw, auctionManager_getLotCurrentPrice)
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

// Postprocess the multicalled data to get the formatted parameters
func (c *AuctionLot) PostprocessAfterMulticall(details *AuctionLotDetails) {
	details.Index = details.IndexRaw.Uint64()
	details.StartBlock = details.StartBlockRaw.Uint64()
	details.EndBlock = details.EndBlockRaw.Uint64()
	details.StartPrice = eth.WeiToEth(details.StartPriceRaw)
	details.ReservePrice = eth.WeiToEth(details.ReservePriceRaw)
	details.PriceAtCurrentBlock = eth.WeiToEth(details.PriceAtCurrentBlockRaw)
	details.PriceByTotalBids = eth.WeiToEth(details.PriceByTotalBidsRaw)
	details.CurrentPrice = eth.WeiToEth(details.CurrentPriceRaw)
}
