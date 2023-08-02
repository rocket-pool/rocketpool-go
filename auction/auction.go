package auction

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
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
	contract *core.Contract
}

// Details for RocketAuctionManager
type AuctionManagerDetails struct {
	TotalRplBalance     *big.Int               `json:"totalRplBalance"`
	AllottedRplBalance  *big.Int               `json:"allottedRplBalance"`
	RemainingRplBalance *big.Int               `json:"remainingRplBalance"`
	LotCount            core.Parameter[uint64] `json:"lotCount"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new AuctionManager contract binding
func NewAuctionManager(rp *rocketpool.RocketPool) (*AuctionManager, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketAuctionManager)
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
func (c *AuctionManager) CreateLot(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "createLot", opts)
}
