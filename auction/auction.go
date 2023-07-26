package auction

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
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

// Multicall details for auction manager
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

// Get all lot details
func (c *AuctionManager) GetLots(rp *rocketpool.RocketPool, opts *bind.CallOpts) ([]AuctionLotDetails, error) {
	// Get lot count
	lotCount, err := c.GetLotCount(opts)
	if err != nil {
		return []AuctionLotDetails{}, err
	}

	// Load lot details in batches
	details := make([]AuctionLotDetails, lotCount)
	for bsi := uint64(0); bsi < lotCount; bsi += lotDetailsBatchSize {

		// Get batch start & end index
		lsi := bsi
		lei := bsi + lotDetailsBatchSize
		if lei > lotCount {
			lei = lotCount
		}

		// Load details
		var wg errgroup.Group
		for li := lsi; li < lei; li++ {
			li := li
			lot := NewAuctionLot(c, li)
			wg.Go(func() error {
				lotDetails, err := lot.GetLotDetails(opts)
				if err == nil {
					details[li] = lotDetails
				}
				return err
			})
		}
		if err := wg.Wait(); err != nil {
			return []AuctionLotDetails{}, err
		}

	}

	// Return
	return details, nil
}

// Get all lot details with bids from an address
func (c *AuctionManager) GetLotsWithBids(rp *rocketpool.RocketPool, bidder common.Address, opts *bind.CallOpts) ([]AuctionLotDetails, error) {
	// Get lot count
	lotCount, err := c.GetLotCount(opts)
	if err != nil {
		return []AuctionLotDetails{}, err
	}

	// Load lot details in batches
	details := make([]AuctionLotDetails, lotCount)
	for bsi := uint64(0); bsi < lotCount; bsi += lotDetailsBatchSize {

		// Get batch start & end index
		lsi := bsi
		lei := bsi + lotDetailsBatchSize
		if lei > lotCount {
			lei = lotCount
		}

		// Load details
		var wg errgroup.Group
		for li := lsi; li < lei; li++ {
			li := li
			lot := NewAuctionLot(c, li)
			wg.Go(func() error {
				lotDetails, err := lot.GetLotDetailsWithBids(bidder, opts)
				if err == nil {
					details[li] = lotDetails
				}
				return err
			})
		}
		if err := wg.Wait(); err != nil {
			return []AuctionLotDetails{}, err
		}

	}

	// Return
	return details, nil
}

// =========================
// === Formatted Getters ===
// =========================

// Get the number of lots for auction
func (c *AuctionManager) GetLotCount(opts *bind.CallOpts) (uint64, error) {
	raw, err := c.GetLotCountRaw(opts)
	if err != nil {
		return 0, err
	}
	return raw.Uint64(), nil
}

// ====================
// === Transactions ===
// ====================

// Get info for creating a new lot
func (c *AuctionManager) GetCreateLotInfo(opts *bind.TransactOpts) (*rocketpool.TransactionInfo, error) {
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
	details.LotCount = details.LotCountRaw.Uint64()
}

// Get all lot details
func (c *AuctionManager) GetLotsViaMulticall(multicallerAddress common.Address, lotCount uint64, opts *bind.CallOpts) ([]AuctionLotDetails, error) {
	lots := make([]*AuctionLot, lotCount)
	lotDetails := make([]AuctionLotDetails, lotCount)

	// Sync
	var wg errgroup.Group
	wg.SetLimit(int(lotDetailsBatchSize))

	// Load lot details in batches
	for i := uint64(0); i < lotCount; i += lotDetailsBatchSize {
		i := i
		max := i + lotDetailsBatchSize
		if max > lotCount {
			max = lotCount
		}

		// Load details
		wg.Go(func() error {
			var err error
			mc, err := multicall.NewMultiCaller(c.rp.Client, multicallerAddress)
			if err != nil {
				return err
			}
			for j := i; j < max; j++ {
				lot := NewAuctionLot(c, j)
				lots[i] = lot
				details := &lotDetails[j]
				lot.AddMulticallQueries(mc, details)
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("error getting lot details: %w", err)
	}

	// Do some postprocessing
	for i := range lotDetails {
		lot := lots[i]
		details := &lotDetails[i]
		lot.PostprocessAfterMulticall(details)
	}

	// Return
	return lotDetails, nil
}

// Get all lot details with bids from an address
func (c *AuctionManager) GetLotsWithBidsViaMulticall(multicallerAddress common.Address, lotCount uint64, bidder common.Address, opts *bind.CallOpts) ([]AuctionLotDetails, error) {
	lots := make([]*AuctionLot, lotCount)
	lotDetails := make([]AuctionLotDetails, lotCount)

	// Sync
	var wg errgroup.Group
	wg.SetLimit(int(lotDetailsBatchSize))

	// Load lot details in batches
	for i := uint64(0); i < lotCount; i += lotDetailsBatchSize {
		i := i
		max := i + lotDetailsBatchSize
		if max > lotCount {
			max = lotCount
		}

		// Load details
		wg.Go(func() error {
			var err error
			mc, err := multicall.NewMultiCaller(c.rp.Client, multicallerAddress)
			if err != nil {
				return err
			}
			for j := i; j < max; j++ {
				lot := NewAuctionLot(c, j)
				lots[i] = lot
				details := &lotDetails[j]
				lot.AddMulticallQueries(mc, details)
				lot.AddBidAmountToMulticallQuery(mc, details, bidder)
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("error getting lot details: %w", err)
	}

	// Do some postprocessing
	for i := range lotDetails {
		lot := lots[i]
		details := &lotDetails[i]
		lot.PostprocessAfterMulticall(details)
	}

	// Return
	return lotDetails, nil
}
