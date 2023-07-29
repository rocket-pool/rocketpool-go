package auction

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	// Settings
	lotDetailsBatchSize uint64 = 10
)

// ===============
// === Structs ===
// ===============

// Binding for RocketAuctionManager
type AuctionManager struct {
	Details  AuctionManagerDetails
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// Details for RocketAuctionManager
type AuctionManagerDetails struct {
	TotalRplBalance     *big.Int                     `json:"totalRplBalance"`
	AllottedRplBalance  *big.Int                     `json:"allottedRplBalance"`
	RemainingRplBalance *big.Int                     `json:"remainingRplBalance"`
	LotCount            rocketpool.Parameter[uint64] `json:"lotCount"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new AuctionManager contract binding
func NewAuctionManager(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*AuctionManager, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketAuctionManager", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting auction manager contract: %w", err)
	}

	return &AuctionManager{
		Details:  AuctionManagerDetails{},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the total RPL balance of the auction contract
func (c *AuctionManager) GetTotalRPLBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.TotalRplBalance, "getTotalRPLBalance")
}

// Get the allotted RPL balance of the auction contract
func (c *AuctionManager) GetAllottedRPLBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.AllottedRplBalance, "getAllottedRPLBalance")
}

// Get the remaining RPL balance of the auction contract
func (c *AuctionManager) GetRemainingRPLBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.RemainingRplBalance, "getRemainingRPLBalance")
}

// Get the number of lots for auction
func (c *AuctionManager) GetLotCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.LotCount.RawValue, "getLotCount")
}

// Get all basic details
func (c *AuctionManager) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetTotalRPLBalance(mc)
	c.GetAllottedRPLBalance(mc)
	c.GetRemainingRPLBalance(mc)
	c.GetLotCount(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for creating a new lot
func (c *AuctionManager) CreateLot(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, "createLot", opts)
}

// ===================
// === Sub-Getters ===
// ===================

// Get a lot with details
func (c *AuctionManager) GetLot(index uint64, opts *bind.CallOpts) (*AuctionLot, error) {
	return c.getLotImpl(index, nil, opts)
}

// Get a lot with details and bids from the provided bidder
func (c *AuctionManager) GetLotWithBids(index uint64, bidder common.Address, opts *bind.CallOpts) (*AuctionLot, error) {
	return c.getLotImpl(index, &bidder, opts)
}

// Get lot implementation
func (c *AuctionManager) getLotImpl(index uint64, bidder *common.Address, opts *bind.CallOpts) (*AuctionLot, error) {
	// Create the lot and get details via a multicall query
	lot, err := multicall.MulticallQuery[AuctionLot](
		c.rp,
		func(mc *multicall.MultiCaller) (*AuctionLot, error) {
			lot := NewAuctionLot(c, index)
			lot.GetAllDetails(mc)
			if bidder != nil {
				lot.GetAllDetailsWithBidAmount(mc, *bidder)
			}
			return lot, nil
		},
		nil,
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting lot: %w", err)
	}

	// Return
	return lot, nil
}

// Get all lot details
func (c *AuctionManager) GetLots(lotCount uint64, opts *bind.CallOpts) ([]*AuctionLot, error) {
	return c.getLotsImpl(lotCount, nil, opts)
}

// Get all lot details with bids from an address
func (c *AuctionManager) GetLotsWithBids(lotCount uint64, bidder common.Address, opts *bind.CallOpts) ([]*AuctionLot, error) {
	return c.getLotsImpl(lotCount, &bidder, opts)
}

// Get lots implementation
func (c *AuctionManager) getLotsImpl(lotCount uint64, bidder *common.Address, opts *bind.CallOpts) ([]*AuctionLot, error) {
	// Run the multicall query for each lot
	lots, err := multicall.MulticallBatchQuery[AuctionLot](
		c.rp,
		lotCount,
		lotDetailsBatchSize,
		func(lots []*AuctionLot, index uint64, mc *multicall.MultiCaller) error {
			lot := NewAuctionLot(c, index)
			lots[index] = lot
			if bidder != nil {
				lot.GetAllDetailsWithBidAmount(mc, *bidder)
			} else {
				lot.GetAllDetails(mc)
			}
			return nil
		},
		nil,
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting lot details: %w", err)
	}

	// Return
	return lots, nil
}
