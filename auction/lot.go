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
	// The index of the auction lot
	Index uint64

	// True if an auction lot with this index exists
	Exists *core.SimpleField[bool]

	// The lot's start block
	StartBlock *core.FormattedUint256Field[uint64]

	// The lot's end block
	EndBlock *core.FormattedUint256Field[uint64]

	// The lot's starting price
	StartPrice *core.FormattedUint256Field[float64]

	// The lot's reserve price
	ReservePrice *core.FormattedUint256Field[float64]

	// The price of the lot in RPL/ETH in the current block
	PriceAtCurrentBlock *core.FormattedUint256Field[float64]

	// The price of the lot by the total bids
	PriceByTotalBids *core.FormattedUint256Field[float64]

	// The price of the lot at the current block
	CurrentPrice *core.FormattedUint256Field[float64]

	// The lot's total RPL
	TotalRplAmount *core.SimpleField[*big.Int]

	// The amount of RPL claimed for the lot
	ClaimedRplAmount *core.SimpleField[*big.Int]

	// The amount of RPL remaining for the lot
	RemainingRplAmount *core.SimpleField[*big.Int]

	// The lot's total bid amount
	TotalBidAmount *core.SimpleField[*big.Int]

	// True if the lot has been cleared already
	IsCleared *core.SimpleField[bool]

	// True if RPL has been recovered by the lot
	RplRecovered *core.SimpleField[bool]

	// === Internal fields ===
	am       *core.Contract
	indexBig *big.Int
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

	indexBig := big.NewInt(0).SetUint64(index)
	return &AuctionLot{
		Index:               index,
		Exists:              core.NewSimpleField[bool](am, "getLotExists", indexBig),
		StartBlock:          core.NewFormattedUint256Field[uint64](am, "getLotStartBlock", indexBig),
		EndBlock:            core.NewFormattedUint256Field[uint64](am, "getLotEndBlock", indexBig),
		StartPrice:          core.NewFormattedUint256Field[float64](am, "getLotStartPrice", indexBig),
		ReservePrice:        core.NewFormattedUint256Field[float64](am, "getLotReservePrice", indexBig),
		PriceAtCurrentBlock: core.NewFormattedUint256Field[float64](am, "getLotPriceAtCurrentBlock", indexBig),
		PriceByTotalBids:    core.NewFormattedUint256Field[float64](am, "getLotPriceByTotalBids", indexBig),
		CurrentPrice:        core.NewFormattedUint256Field[float64](am, "getLotCurrentPrice", indexBig),
		TotalRplAmount:      core.NewSimpleField[*big.Int](am, "getLotTotalRPLAmount", indexBig),
		ClaimedRplAmount:    core.NewSimpleField[*big.Int](am, "getLotClaimedRPLAmount", indexBig),
		RemainingRplAmount:  core.NewSimpleField[*big.Int](am, "getLotRemainingRPLAmount", indexBig),
		TotalBidAmount:      core.NewSimpleField[*big.Int](am, "getLotTotalBidAmount", indexBig),
		IsCleared:           core.NewSimpleField[bool](am, "getLotIsCleared", indexBig),
		RplRecovered:        core.NewSimpleField[bool](am, "getLotRPLRecovered", indexBig),

		am:       am,
		indexBig: indexBig,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the price of the lot at the given block
func (c *AuctionLot) GetLotPriceAtBlock(mc *batch.MultiCaller, out **big.Int, blockNumber uint64) {
	core.AddCall(mc, c.am, out, "getLotPriceAtBlock", c.indexBig, big.NewInt(int64(blockNumber)))
}

// Get the ETH amount bid on the lot by an address
func (c *AuctionLot) GetLotAddressBidAmount(mc *batch.MultiCaller, out **big.Int, bidder common.Address) {
	core.AddCall(mc, c.am, out, "getLotAddressBidAmount", c.indexBig, bidder)
}

// ====================
// === Transactions ===
// ====================

// Get info for placing a bid on a lot
func (c *AuctionLot) PlaceBid(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.am, "placeBid", opts, c.indexBig)
}

// Get info for claiming RPL from a lot that was bid on
func (c *AuctionLot) ClaimBid(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.am, "claimBid", opts, c.indexBig)
}

// Get info for recovering unclaimed RPL from a lot
func (c *AuctionLot) RecoverUnclaimedRpl(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.am, "recoverUnclaimedRPL", opts, c.indexBig)
}
