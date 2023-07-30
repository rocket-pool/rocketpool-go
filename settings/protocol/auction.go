package protocol

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	auctionSettingsContractName string = "rocketDAOProtocolSettingsAuction"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProtocolSettingsAuction
type DaoProtocolSettingsAuction struct {
	Details             DaoProtocolSettingsAuctionDetails
	rp                  *rocketpool.RocketPool
	contract            *core.Contract
	daoProtocolContract *protocol.DaoProtocol
}

// Details for RocketDAOProtocolSettingsAuction
type DaoProtocolSettingsAuctionDetails struct {
	IsCreateLotEnabled    bool                    `json:"isCreateLotEnabled"`
	IsBidOnLotEnabled     bool                    `json:"isBidOnLotEnabled"`
	LotMinimumEthValue    *big.Int                `json:"lotMinimumEthValue"`
	LotMaximumEthValue    *big.Int                `json:"lotMaximumEthValue"`
	LotDuration           core.Parameter[uint64]  `json:"lotDuration"`
	LotStartingPriceRatio core.Parameter[float64] `json:"lotStartingPriceRatio"`
	LotReservePriceRatio  core.Parameter[float64] `json:"lotReservePriceRatio"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProtocolSettingsAuction contract binding
func NewDaoProtocolSettingsAuction(rp *rocketpool.RocketPool, daoProtocolContract *protocol.DaoProtocol, opts *bind.CallOpts) (*DaoProtocolSettingsAuction, error) {
	// Create the contract
	contract, err := rp.GetContract(auctionSettingsContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO protocol settings auction contract: %w", err)
	}

	return &DaoProtocolSettingsAuction{
		Details:             DaoProtocolSettingsAuctionDetails{},
		rp:                  rp,
		contract:            contract,
		daoProtocolContract: daoProtocolContract,
	}, nil
}

// =============
// === Calls ===
// =============

// Check if lot creation is currently enabled
func (c *DaoProtocolSettingsAuction) GetCreateLotEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.IsCreateLotEnabled, "getCreateLotEnabled")
}

// Check if lot bidding is currently enabled
func (c *DaoProtocolSettingsAuction) GetBidOnLotEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.IsBidOnLotEnabled, "getBidOnLotEnabled")
}

// Get the minimum lot size in ETH
func (c *DaoProtocolSettingsAuction) GetLotMinimumEthValue(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.LotMinimumEthValue, "getLotMinimumEthValue")
}

// Get the maximum lot size in ETH
func (c *DaoProtocolSettingsAuction) GetLotMaximumEthValue(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.LotMaximumEthValue, "getLotMaximumEthValue")
}

// Get the lot duration, in blocks
func (c *DaoProtocolSettingsAuction) GetLotDuration(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.LotDuration.RawValue, "getLotDuration")
}

// Get the lot starting price relative to current ETH price, as a fraction
func (c *DaoProtocolSettingsAuction) GetLotStartingPriceRatio(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.LotStartingPriceRatio.RawValue, "getStartingPriceRatio")
}

// Get the reserve price relative to current ETH price, as a fraction
func (c *DaoProtocolSettingsAuction) GetLotReservePriceRatio(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.LotReservePriceRatio.RawValue, "getReservePriceRatio")
}

// Get all basic details
func (c *DaoProtocolSettingsAuction) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetCreateLotEnabled(mc)
	c.GetBidOnLotEnabled(mc)
	c.GetLotMinimumEthValue(mc)
	c.GetLotMaximumEthValue(mc)
	c.GetLotDuration(mc)
	c.GetLotStartingPriceRatio(mc)
	c.GetLotReservePriceRatio(mc)
}

// ====================
// === Transactions ===
// ====================

// Set the create lot enabled flag
func (c *DaoProtocolSettingsAuction) BootstrapCreateLotEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapBool(auctionSettingsContractName, "auction.lot.create.enabled", value, opts)
}

// Set the create lot enabled flag
func (c *DaoProtocolSettingsAuction) BootstrapBidOnLotEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapBool(auctionSettingsContractName, "auction.lot.bidding.enabled", value, opts)
}

// Set the minimum ETH value for lots
func (c *DaoProtocolSettingsAuction) BootstrapLotMinimumEthValue(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(auctionSettingsContractName, "auction.lot.value.minimum", value, opts)
}

// Set the maximum ETH value for lots
func (c *DaoProtocolSettingsAuction) BootstrapLotMaximumEthValue(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(auctionSettingsContractName, "auction.lot.value.maximum", value, opts)
}

// Set the duration value for lots, in blocks
func (c *DaoProtocolSettingsAuction) BootstrapLotDuration(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(auctionSettingsContractName, "auction.lot.duration", big.NewInt(int64(value)), opts)
}

// Set the starting price ratio for lots
func (c *DaoProtocolSettingsAuction) BootstrapLotStartingPriceRatio(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(auctionSettingsContractName, "auction.price.start", eth.EthToWei(value), opts)
}

// Set the reserve price ratio for lots
func (c *DaoProtocolSettingsAuction) BootstrapLotReservePriceRatio(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(auctionSettingsContractName, "auction.price.reserve", eth.EthToWei(value), opts)
}
