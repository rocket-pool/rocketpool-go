package trustednode

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	scrubPeriodPath               = "minipool.scrub.period"
	promotionScrubPeriodPath      = "minipool.promotion.scrub.period"
	scrubPenaltyEnabledPath       = "minipool.scrub.penalty.enabled"
	bondReductionWindowStartPath  = "minipool.bond.reduction.window.start"
	bondReductionWindowLengthPath = "minipool.bond.reduction.window.length"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAONodeTrustedSettingsMinipool
type DaoNodeTrustedSettingsMinipool struct {
	Details                         DaoNodeTrustedSettingsMinipoolDetails
	rp                              *rocketpool.RocketPool
	contract                        *core.Contract
	daoNodeTrustedContract          *trustednode.DaoNodeTrusted
	daoNodeTrustedProposalsContract *trustednode.DaoNodeTrustedProposals
}

// Details for RocketDAONodeTrustedSettingsMinipool
type DaoNodeTrustedSettingsMinipoolDetails struct {
	ScrubPeriod               core.Parameter[time.Duration] `json:"scrubPeriod"`
	PromotionScrubPeriod      core.Parameter[time.Duration] `json:"promotionScrubPeriod"`
	IsScrubPenaltyEnabled     bool                          `json:"isScrubPenaltyEnabled"`
	BondReductionWindowStart  core.Parameter[time.Duration] `json:"bondReductionWindowStart"`
	BondReductionWindowLength core.Parameter[time.Duration] `json:"bondReductionWindowLength"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoNodeTrustedSettingsMinipool contract binding
func NewDaoNodeTrustedSettingsMinipool(rp *rocketpool.RocketPool) (*DaoNodeTrustedSettingsMinipool, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted settings minipool contract: %w", err)
	}
	daoNodeTrustedContract, err := trustednode.NewDaoNodeTrusted(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted contract: %w", err)
	}
	daoNodeTrustedProposalsContract, err := trustednode.NewDaoNodeTrustedProposals(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted proposals contract: %w", err)
	}

	return &DaoNodeTrustedSettingsMinipool{
		Details:                         DaoNodeTrustedSettingsMinipoolDetails{},
		rp:                              rp,
		contract:                        contract,
		daoNodeTrustedContract:          daoNodeTrustedContract,
		daoNodeTrustedProposalsContract: daoNodeTrustedProposalsContract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the amount of time, in seconds, the scrub check lasts before a minipool can move from prelaunch to staking
func (c *DaoNodeTrustedSettingsMinipool) GetScrubPeriod(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.ScrubPeriod.RawValue, "getScrubPeriod")
}

// Get the amount of time, in seconds, the promotion scrub check lasts before a vacant minipool can be promoted
func (c *DaoNodeTrustedSettingsMinipool) GetPromotionScrubPeriod(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.PromotionScrubPeriod.RawValue, "getPromotionScrubPeriod")
}

// Check if the RPL slashing penalty is applied to scrubbed minipools
func (c *DaoNodeTrustedSettingsMinipool) GetScrubPenaltyEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.IsScrubPenaltyEnabled, "getScrubPenaltyEnabled")
}

// Get the amount of time, in seconds, a minipool must wait after beginning a bond reduction before it can apply the bond reduction (how long the Oracle DAO has to cancel the reduction if required)
func (c *DaoNodeTrustedSettingsMinipool) GetBondReductionWindowStart(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.BondReductionWindowStart.RawValue, "getBondReductionWindowStart")
}

// Get the amount of time, in seconds, a minipool has to reduce its bond once it has passed the check window
func (c *DaoNodeTrustedSettingsMinipool) GetBondReductionWindowLength(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.BondReductionWindowLength.RawValue, "getBondReductionWindowLength")
}

// Get all basic details
func (c *DaoNodeTrustedSettingsMinipool) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetScrubPeriod(mc)
	c.GetPromotionScrubPeriod(mc)
	c.GetScrubPenaltyEnabled(mc)
	c.GetBondReductionWindowStart(mc)
	c.GetBondReductionWindowLength(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for setting the amount of time, in seconds, the scrub check lasts before a minipool can move from prelaunch to staking
func (c *DaoNodeTrustedSettingsMinipool) BootstrapScrubPeriod(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, scrubPeriodPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the amount of time, in seconds, the promotion scrub check lasts before a vacant minipool can be promoted
func (c *DaoNodeTrustedSettingsMinipool) BootstrapPromotionScrubPeriod(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, promotionScrubPeriodPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the flag for the RPL slashing penalty on scrubbed minipools
func (c *DaoNodeTrustedSettingsMinipool) BootstrapScrubPenaltyEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapBool(rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, scrubPenaltyEnabledPath, value, opts)
}

// Get info for setting the amount of time, in seconds, a minipool must wait after beginning a bond reduction before it can apply the bond reduction (how long the Oracle DAO has to cancel the reduction if required)
func (c *DaoNodeTrustedSettingsMinipool) BootstrapBondReductionWindowStart(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, bondReductionWindowStartPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the amount of time, in seconds, a minipool has to reduce its bond once it has passed the check window
func (c *DaoNodeTrustedSettingsMinipool) BootstrapBondReductionWindowLength(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, bondReductionWindowLengthPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the amount of time, in seconds, the scrub check lasts before a minipool can move from prelaunch to staking
func (c *DaoNodeTrustedSettingsMinipool) ProposeScrubPeriod(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", scrubPeriodPath), rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, scrubPeriodPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the amount of time, in seconds, the promotion scrub check lasts before a vacant minipool can be promoted
func (c *DaoNodeTrustedSettingsMinipool) ProposePromotionScrubPeriod(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", promotionScrubPeriodPath), rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, promotionScrubPeriodPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the flag for the RPL slashing penalty on scrubbed minipools
func (c *DaoNodeTrustedSettingsMinipool) ProposeScrubPenaltyEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetBool(fmt.Sprintf("set %s", scrubPenaltyEnabledPath), rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, scrubPenaltyEnabledPath, value, opts)
}

// Get info for setting the amount of time, in seconds, a minipool must wait after beginning a bond reduction before it can apply the bond reduction (how long the Oracle DAO has to cancel the reduction if required)
func (c *DaoNodeTrustedSettingsMinipool) ProposeBondReductionWindowStart(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", bondReductionWindowStartPath), rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, bondReductionWindowStartPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the amount of time, in seconds, a minipool has to reduce its bond once it has passed the check window
func (c *DaoNodeTrustedSettingsMinipool) ProposeBondReductionWindowLength(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", bondReductionWindowLengthPath), rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, bondReductionWindowLengthPath, big.NewInt(int64(value)), opts)
}
