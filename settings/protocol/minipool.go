package protocol

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	minipoolSettingsContractName string = "rocketDAOProtocolSettingsMinipool"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProtocolSettingsMinipool
type DaoProtocolSettingsMinipool struct {
	Details             DaoProtocolSettingsMinipoolDetails
	rp                  *rocketpool.RocketPool
	contract            *core.Contract
	daoProtocolContract *protocol.DaoProtocol
}

// Details for RocketDAOProtocolSettingsMinipool
type DaoProtocolSettingsMinipoolDetails struct {
	LaunchBalance               *big.Int                      `json:"launchBalance"`
	FullDepositNodeAmount       *big.Int                      `json:"fullDepositNodeAmount"`
	HalfDepositNodeAmount       *big.Int                      `json:"halfDepositNodeAmount"`
	EmptyDepositNodeAmount      *big.Int                      `json:"emptyDepositNodeAmount"`
	FullDepositUserAmount       *big.Int                      `json:"fullDepositUserAmount"`
	HalfDepositUserAmount       *big.Int                      `json:"halfDepositUserAmount"`
	EmptyDepositUserAmount      *big.Int                      `json:"emptyDepositUserAmount"`
	IsSubmitWithdrawableEnabled bool                          `json:"isSubmitWithdrawableEnabled"`
	LaunchTimeout               core.Parameter[time.Duration] `json:"launchTimeout"`
	IsBondReductionEnabled      bool                          `json:"isBondReductionEnabled"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProtocolSettingsMinipool contract binding
func NewDaoProtocolSettingsMinipool(rp *rocketpool.RocketPool, daoProtocolContract *protocol.DaoProtocol, opts *bind.CallOpts) (*DaoProtocolSettingsMinipool, error) {
	// Create the contract
	contract, err := rp.GetContract(minipoolSettingsContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO protocol settings minipool contract: %w", err)
	}

	return &DaoProtocolSettingsMinipool{
		Details:             DaoProtocolSettingsMinipoolDetails{},
		rp:                  rp,
		contract:            contract,
		daoProtocolContract: daoProtocolContract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the minipool launch balance
func (c *DaoProtocolSettingsMinipool) GetLaunchBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.LaunchBalance, "getLaunchBalance")
}

// Get the amount required from the node for a full deposit
func (c *DaoProtocolSettingsMinipool) GetFullDepositNodeAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.FullDepositNodeAmount, "getFullDepositNodeAmount")
}

// Get the amount required from the node for a half deposit
func (c *DaoProtocolSettingsMinipool) GetHalfDepositNodeAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.HalfDepositNodeAmount, "getHalfDepositNodeAmount")
}

// Get the amount required from the node for an empty deposit
func (c *DaoProtocolSettingsMinipool) GetEmptyDepositNodeAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.EmptyDepositNodeAmount, "getEmptyDepositNodeAmount")
}

// Get the amount required from the pool stakers for a full deposit
func (c *DaoProtocolSettingsMinipool) GetFullDepositUserAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.FullDepositUserAmount, "getFullDepositUserAmount")
}

// Get the amount required from the pool stakers for a half deposit
func (c *DaoProtocolSettingsMinipool) GetHalfDepositUserAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.HalfDepositUserAmount, "getHalfDepositUserAmount")
}

// Get the amount required from the pool stakers for an empty deposit
func (c *DaoProtocolSettingsMinipool) GetEmptyDepositUserAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.EmptyDepositUserAmount, "getEmptyDepositUserAmount")
}

// Check if minipool withdrawable event submissions are currently enabled
func (c *DaoProtocolSettingsMinipool) GetSubmitWithdrawableEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.IsSubmitWithdrawableEnabled, "getSubmitWithdrawableEnabled")
}

// Get the timeout period, in seconds, for prelaunch minipools to launch
func (c *DaoProtocolSettingsMinipool) GetLaunchTimeout(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.LaunchTimeout.RawValue, "getLaunchTimeout")
}

// Check if minipool bond reductions are currently enabled
func (c *DaoProtocolSettingsMinipool) GetBondReductionEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.IsBondReductionEnabled, "getBondReductionEnabled")
}

// Get all basic details
func (c *DaoProtocolSettingsMinipool) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetLaunchBalance(mc)
	c.GetFullDepositNodeAmount(mc)
	c.GetHalfDepositNodeAmount(mc)
	c.GetEmptyDepositNodeAmount(mc)
	c.GetFullDepositUserAmount(mc)
	c.GetHalfDepositUserAmount(mc)
	c.GetEmptyDepositUserAmount(mc)
	c.GetSubmitWithdrawableEnabled(mc)
	c.GetLaunchTimeout(mc)
	c.GetBondReductionEnabled(mc)
}

// ====================
// === Transactions ===
// ====================

// Set the flag for enabling minipool withdrawable event submissions
func (c *DaoProtocolSettingsMinipool) BootstrapSubmitWithdrawableEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapBool(minipoolSettingsContractName, "minipool.submit.withdrawable.enabled", value, opts)
}

func (c *DaoProtocolSettingsMinipool) BootstrapLaunchTimeout(value time.Duration, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(minipoolSettingsContractName, "minipool.launch.timeout", big.NewInt(int64(value.Seconds())), opts)
}

func (c *DaoProtocolSettingsMinipool) BootstrapBondReductionEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapBool(minipoolSettingsContractName, "minipool.bond.reduction.enabled", value, opts)
}
