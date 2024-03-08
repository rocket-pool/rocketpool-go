package auction

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rocket-pool/node-manager-core/eth"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketAuctionManager
type AuctionManager struct {
	// The total RPL balance of the auction contract
	TotalRplBalance *core.SimpleField[*big.Int]

	// The allotted RPL balance of the auction contract
	AllottedRplBalance *core.SimpleField[*big.Int]

	// The remaining RPL balance of the auction contract
	RemainingRplBalance *core.SimpleField[*big.Int]

	// The number of lots for auction
	LotCount *core.FormattedUint256Field[uint64]

	// === Internal fields ===
	rp    *rocketpool.RocketPool
	am    *core.Contract
	txMgr *eth.TransactionManager
}

// Details for RocketAuctionManager
type AuctionManagerDetails struct {
}

// ====================
// === Constructors ===
// ====================

// Creates a new AuctionManager contract binding
func NewAuctionManager(rp *rocketpool.RocketPool) (*AuctionManager, error) {
	// Create the contract
	am, err := rp.GetContract(rocketpool.ContractName_RocketAuctionManager)
	if err != nil {
		return nil, fmt.Errorf("error getting auction manager contract: %w", err)
	}

	return &AuctionManager{
		TotalRplBalance:     core.NewSimpleField[*big.Int](am, "getTotalRPLBalance"),
		AllottedRplBalance:  core.NewSimpleField[*big.Int](am, "getAllottedRPLBalance"),
		RemainingRplBalance: core.NewSimpleField[*big.Int](am, "getRemainingRPLBalance"),
		LotCount:            core.NewFormattedUint256Field[uint64](am, "getLotCount"),

		rp:    rp,
		am:    am,
		txMgr: rp.GetTransactionManager(),
	}, nil
}

// ====================
// === Transactions ===
// ====================

// Get info for creating a new lot
func (c *AuctionManager) CreateLot(opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.am.Contract, "createLot", opts)
}
