package auction

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	// Settings
	lotDetailsBatchSize uint64 = 10

	// Contract names
	AuctionManager_ContractName string = "rocketAuctionManager"

	// Calls
	auctionManager_getTotalRPLBalance     string = "getTotalRPLBalance"
	auctionManager_getAllottedRPLBalance  string = "getAllottedRPLBalance"
	auctionManager_getRemainingRPLBalance string = "getRemainingRPLBalance"
	auctionManager_getLotCount            string = "getLotCount"

	// Transactions
	auctionManager_createLot string = "createLot"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketAuctionManager
type AuctionManager struct {
	rp       *rocketpool.RocketPool
	contract *rocketpool.Contract
}

// Details for auction manager
type AuctionManagerDetails struct {
	// Raw parameters
	TotalRplBalance     *big.Int `json:"totalRplBalance"`
	AllottedRplBalance  *big.Int `json:"allottedRplBalance"`
	RemainingRplBalance *big.Int `json:"remainingRplBalance"`
	LotCountRaw         *big.Int `json:"lotCountRaw"`

	// Formatted parameters
	LotCount uint64 `json:"lotCount"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new AuctionManager contract binding
func NewAuctionManager(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*AuctionManager, error) {
	// Create the contract
	contract, err := rp.GetContract(AuctionManager_ContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting auction manager contract: %w", err)
	}

	return &AuctionManager{
		rp:       rp,
		contract: contract,
	}, nil
}

// ===================
// === Raw Getters ===
// ===================

// Get the total RPL balance of the auction contract
func (c *AuctionManager) GetTotalRPLBalance(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, auctionManager_getTotalRPLBalance)
}

// Get the allotted RPL balance of the auction contract
func (c *AuctionManager) GetAllottedRPLBalance(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, auctionManager_getAllottedRPLBalance)
}

// Get the remaining RPL balance of the auction contract
func (c *AuctionManager) GetRemainingRPLBalance(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, auctionManager_getRemainingRPLBalance)
}

// Get the number of lots for auction
func (c *AuctionManager) GetLotCountRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.contract, opts, auctionManager_getLotCount)
}

// =========================
// === Formatted Getters ===
// =========================

// Get the number of lots for auction
func (c *AuctionManager) GetLotCount(opts *bind.CallOpts) (uint64, error) {
	return utils.ConvertRawToUint(c.GetLotCountRaw(opts))
}

// ====================
// === Transactions ===
// ====================

// Get info for creating a new lot
func (c *AuctionManager) CreateLot(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
	return rocketpool.NewTransactionInfo(c.contract, auctionManager_createLot, opts)
}

// =================
// === Multicall ===
// =================

// Add queries to a multicall batcher
func (c *AuctionManager) AddMulticallQueries(mc *multicall.MultiCaller, details *AuctionManagerDetails) {
	mc.AddCall(c.contract, &details.TotalRplBalance, auctionManager_getTotalRPLBalance)
	mc.AddCall(c.contract, &details.AllottedRplBalance, auctionManager_getAllottedRPLBalance)
	mc.AddCall(c.contract, &details.RemainingRplBalance, auctionManager_getRemainingRPLBalance)
	mc.AddCall(c.contract, &details.LotCountRaw, auctionManager_getLotCount)
}

// Postprocess the multicalled data to get the formatted parameters
func (c *AuctionManager) PostprocessAfterMulticall(details *AuctionManagerDetails) {
	details.LotCount = utils.ConvertToUint(details.LotCountRaw)
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
		func(mc *multicall.MultiCaller) *AuctionLot {
			lot := NewAuctionLot(c, index)
			lot.AddMulticallQueries(mc, &lot.Details)
			if bidder != nil {
				lot.AddBidAmountToMulticallQuery(mc, &lot.Details, *bidder)
			}
			return lot
		},
		func(lot *AuctionLot) {
			lot.PostprocessAfterMulticall(&lot.Details)
		},
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
		func(lots []*AuctionLot, index uint64, mc *multicall.MultiCaller) {
			lot := NewAuctionLot(c, index)
			lots[index] = lot
			details := &lot.Details
			lot.AddMulticallQueries(mc, details)
			if bidder != nil {
				lot.AddBidAmountToMulticallQuery(mc, details, *bidder)
			}
		},
		func(lot *AuctionLot) {
			lot.PostprocessAfterMulticall(&lot.Details)
		},
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting lot details: %w", err)
	}

	// Return
	return lots, nil
}
