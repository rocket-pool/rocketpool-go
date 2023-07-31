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

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProtocolSettingsRewards
type DaoProtocolSettingsRewards struct {
	Details             DaoProtocolSettingsRewardsDetails
	rp                  *rocketpool.RocketPool
	contract            *core.Contract
	daoProtocolContract *protocol.DaoProtocol
}

// Details for RocketDAOProtocolSettingsRewards
type DaoProtocolSettingsRewardsDetails struct {
	ClaimerPercentageTotal core.Parameter[float64]       `json:"claimerPercentageTotal"`
	ClaimIntervalTime      core.Parameter[time.Duration] `json:"claimIntervalTime"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProtocolSettingsRewards contract binding
func NewDaoProtocolSettingsRewards(rp *rocketpool.RocketPool, daoProtocolContract *protocol.DaoProtocol, opts *bind.CallOpts) (*DaoProtocolSettingsRewards, error) {
	// Create the contract
	contract, err := rp.GetContract(rewardsSettingsContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO protocol settings rewards contract: %w", err)
	}

	return &DaoProtocolSettingsRewards{
		Details:             DaoProtocolSettingsRewardsDetails{},
		rp:                  rp,
		contract:            contract,
		daoProtocolContract: daoProtocolContract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the total claim amount for all claimers as a fraction
func (c *DaoProtocolSettingsRewards) GetClaimerPercentageTotal(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.ClaimerPercentageTotal.RawValue, "getRewardsClaimersPercTotal")
}

// Get the rewards claim interval time
func (c *DaoProtocolSettingsRewards) GetClaimIntervalTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.ClaimIntervalTime.RawValue, "getRewardsClaimIntervalTime")
}

// Get all basic details
func (c *DaoProtocolSettingsRewards) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetClaimerPercentageTotal(mc)
	c.GetClaimIntervalTime(mc)
}

// Get the claim amount for a claimer as a fraction
func (c *DaoProtocolSettingsRewards) GetClaimerPercentage(mc *multicall.MultiCaller, percentage_Out *core.Parameter[float64], contractName string) {
	multicall.AddCall(mc, c.contract, &percentage_Out.RawValue, "getRewardsClaimerPerc", contractName)
}

// Get the time that a claimer's share was last updated
func (c *DaoProtocolSettingsRewards) GetClaimerPercentageTimeUpdated(mc *multicall.MultiCaller, time_Out *core.Parameter[time.Time], contractName string) {
	multicall.AddCall(mc, c.contract, &time_Out.RawValue, "getRewardsClaimerPercTimeUpdated", contractName)
}

// ====================
// === Transactions ===
// ====================

func (c *DaoProtocolSettingsRewards) BootstrapClaimIntervalTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocolContract.BootstrapUint(rewardsSettingsContractName, "rpl.rewards.claim.period.time", big.NewInt(int64(value)), opts)
}
