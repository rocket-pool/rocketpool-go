package protocol

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProtocolSettingsInflation
type DaoProtocolSettingsInflation struct {
	Details             DaoProtocolSettingsInflationDetails
	rp                  *rocketpool.RocketPool
	contract            *core.Contract
	daoProtocolContract *protocol.DaoProtocol
}

// Details for RocketDAOProtocolSettingsInflation
type DaoProtocolSettingsInflationDetails struct {
	IntervalRate core.Parameter[float64]   `json:"intervalRate"`
	StartTime    core.Parameter[time.Time] `json:"startTime"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProtocolSettingsInflation contract binding
func NewDaoProtocolSettingsInflation(rp *rocketpool.RocketPool, daoProtocolContract *protocol.DaoProtocol, opts *bind.CallOpts) (*DaoProtocolSettingsInflation, error) {
	// Create the contract
	contract, err := rp.GetContract(inflationSettingsContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO protocol settings inflation contract: %w", err)
	}

	return &DaoProtocolSettingsInflation{
		Details:             DaoProtocolSettingsInflationDetails{},
		rp:                  rp,
		contract:            contract,
		daoProtocolContract: daoProtocolContract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the RPL inflation rate per interval
func (c *DaoProtocolSettingsInflation) GetIntervalRate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.IntervalRate.RawValue, "getInflationIntervalRate")
}

// Get the RPL inflation start time
func (c *DaoProtocolSettingsInflation) GetStartTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.StartTime.RawValue, "getInflationIntervalStartTime")
}

// Get all basic details
func (c *DaoProtocolSettingsInflation) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetIntervalRate(mc)
	c.GetStartTime(mc)
}

// ====================
// === Transactions ===
// ====================

// Set the RPL inflation rate per interval
func (c *DaoProtocolSettingsInflation) BootstrapIntervalRate(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(inflationSettingsContractName, "rpl.inflation.interval.rate", eth.EthToWei(value), opts)
}

// Set the RPL inflation start time
func (c *DaoProtocolSettingsInflation) BootstrapStartTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(inflationSettingsContractName, "rpl.inflation.interval.start", big.NewInt(int64(value)), opts)
}
