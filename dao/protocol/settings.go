package protocol

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for Protocol DAO settings
type ProtocolDaoSettings struct {
	Details           ProtocolDaoSettingsDetails
	AuctionContract   *core.Contract
	DepositContract   *core.Contract
	InflationContract *core.Contract
	MinipoolContract  *core.Contract
	NetworkContract   *core.Contract
	NodeContract      *core.Contract
	RewardsContract   *core.Contract

	rp      *rocketpool.RocketPool
	pdaoMgr *ProtocolDaoManager
}

// Details for Protocol DAO settings
type ProtocolDaoSettingsDetails struct {
	Auction struct {
		IsCreateLotEnabled    bool                    `json:"isCreateLotEnabled"`
		IsBidOnLotEnabled     bool                    `json:"isBidOnLotEnabled"`
		LotMinimumEthValue    *big.Int                `json:"lotMinimumEthValue"`
		LotMaximumEthValue    *big.Int                `json:"lotMaximumEthValue"`
		LotDuration           core.Parameter[uint64]  `json:"lotDuration"`
		LotStartingPriceRatio core.Parameter[float64] `json:"lotStartingPriceRatio"`
		LotReservePriceRatio  core.Parameter[float64] `json:"lotReservePriceRatio"`
	} `json:"auction"`

	Deposit struct {
		IsDepositingEnabled                    bool                    `json:"isDepositingEnabled"`
		AreDepositAssignmentsEnabled           bool                    `json:"areDepositAssignmentsEnabled"`
		MinimumDeposit                         *big.Int                `json:"minimumDeposit"`
		MaximumDepositPoolSize                 *big.Int                `json:"maximumDepositPoolSize"`
		MaximumAssignmentsPerDeposit           core.Parameter[uint64]  `json:"maximumAssignmentsPerDeposit"`
		MaximumSocialisedAssignmentsPerDeposit core.Parameter[uint64]  `json:"maximumSocialisedAssignmentsPerDeposit"`
		DepositFee                             core.Parameter[float64] `json:"depositFee"`
	} `json:"deposit"`

	Inflation struct {
		IntervalRate core.Parameter[float64]   `json:"intervalRate"`
		StartTime    core.Parameter[time.Time] `json:"startTime"`
	} `json:"inflation"`

	Minipool struct {
		IsSubmitWithdrawableEnabled bool                          `json:"isSubmitWithdrawableEnabled"`
		LaunchTimeout               core.Parameter[time.Duration] `json:"launchTimeout"`
		IsBondReductionEnabled      bool                          `json:"isBondReductionEnabled"`
		MaximumCount                core.Parameter[uint64]        `json:"maximumCount"`
		UserDistributeWindowStart   core.Parameter[time.Duration] `json:"userDistributeWindowStart"`
		UserDistributeWindowLength  core.Parameter[time.Duration] `json:"userDistributeWindowLength"`
	} `json:"minipool"`

	Network struct {
		OracleDaoConsensusThreshold core.Parameter[float64]       `json:"oracleDaoConsensusThreshold"`
		NodePenaltyThreshold        core.Parameter[float64]       `json:"nodePenaltyThreshold"`
		PerPenaltyRate              core.Parameter[float64]       `json:"perPenaltyRate"`
		IsSubmitBalancesEnabled     bool                          `json:"isSubmitBalancesEnabled"`
		SubmitBalancesFrequency     core.Parameter[time.Duration] `json:"submitBalancesFrequency"`
		IsSubmitPricesEnabled       bool                          `json:"isSubmitPricesEnabled"`
		SubmitPricesFrequency       core.Parameter[time.Duration] `json:"submitPricesFrequency"`
		MinimumNodeFee              core.Parameter[float64]       `json:"minimumNodeFee"`
		TargetNodeFee               core.Parameter[float64]       `json:"targetNodeFee"`
		MaximumNodeFee              core.Parameter[float64]       `json:"maximumNodeFee"`
		NodeFeeDemandRange          *big.Int                      `json:"nodeFeeDemandRange"`
		TargetRethCollateralRate    core.Parameter[float64]       `json:"targetRethCollateralRate"`
		RethDepositDelay            core.Parameter[time.Duration] `json:"rethDepositDelay"`
		IsSubmitRewardsEnabled      bool                          `json:"isSubmitRewardsEnabled"`
	} `json:"network"`

	Node struct {
		IsRegistrationEnabled              bool                    `json:"isRegistrationEnabled"`
		IsSmoothingPoolRegistrationEnabled bool                    `json:"isSmoothingPoolRegistrationEnabled"`
		IsDepositingEnabled                bool                    `json:"isDepositingEnabled"`
		AreVacantMinipoolsEnabled          bool                    `json:"areVacantMinipoolsEnabled"`
		MinimumPerMinipoolStake            core.Parameter[float64] `json:"minimumPerMinipoolStake"`
		MaximumPerMinipoolStake            core.Parameter[float64] `json:"maximumPerMinipoolStake"`
	} `json:"node"`

	Rewards struct {
		IntervalTime core.Parameter[time.Duration] `json:"intervalTime"`
	} `json:"rewards"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new ProtocolDaoSettings binding
func newProtocolDaoSettings(pdaoMgr *ProtocolDaoManager) (*ProtocolDaoSettings, error) {
	// Get the contracts
	contracts, err := pdaoMgr.rp.GetContracts([]rocketpool.ContractName{
		rocketpool.ContractName_RocketDAOProtocolSettingsAuction,
		rocketpool.ContractName_RocketDAOProtocolSettingsDeposit,
		rocketpool.ContractName_RocketDAOProtocolSettingsInflation,
		rocketpool.ContractName_RocketDAOProtocolSettingsMinipool,
		rocketpool.ContractName_RocketDAOProtocolSettingsNetwork,
		rocketpool.ContractName_RocketDAOProtocolSettingsNode,
		rocketpool.ContractName_RocketDAOProtocolSettingsRewards,
	}...)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO settings contracts: %w", err)
	}

	return &ProtocolDaoSettings{
		Details: ProtocolDaoSettingsDetails{},
		rp:      pdaoMgr.rp,
		pdaoMgr: pdaoMgr,

		AuctionContract:   contracts[0],
		DepositContract:   contracts[1],
		InflationContract: contracts[2],
		MinipoolContract:  contracts[3],
		NetworkContract:   contracts[4],
		NodeContract:      contracts[5],
		RewardsContract:   contracts[6],
	}, nil
}

// =============
// === Calls ===
// =============

// === RocketDAOProtocolSettingsAuction ===

// Check if lot creation is currently enabled
func (c *ProtocolDaoSettings) GetCreateAuctionLotEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.AuctionContract, &c.Details.Auction.IsCreateLotEnabled, "getCreateLotEnabled")
}

// Check if lot bidding is currently enabled
func (c *ProtocolDaoSettings) GetBidOnAuctionLotEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.AuctionContract, &c.Details.Auction.IsBidOnLotEnabled, "getBidOnLotEnabled")
}

// Get the minimum lot size in ETH
func (c *ProtocolDaoSettings) GetAuctionLotMinimumEthValue(mc *batch.MultiCaller) {
	core.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotMinimumEthValue, "getLotMinimumEthValue")
}

// Get the maximum lot size in ETH
func (c *ProtocolDaoSettings) GetAuctionLotMaximumEthValue(mc *batch.MultiCaller) {
	core.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotMaximumEthValue, "getLotMaximumEthValue")
}

// Get the lot duration, in blocks
func (c *ProtocolDaoSettings) GetAuctionLotDuration(mc *batch.MultiCaller) {
	core.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotDuration.RawValue, "getLotDuration")
}

// Get the lot starting price relative to current ETH price, as a fraction
func (c *ProtocolDaoSettings) GetAuctionLotStartingPriceRatio(mc *batch.MultiCaller) {
	core.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotStartingPriceRatio.RawValue, "getStartingPriceRatio")
}

// Get the reserve price relative to current ETH price, as a fraction
func (c *ProtocolDaoSettings) GetAuctionLotReservePriceRatio(mc *batch.MultiCaller) {
	core.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotReservePriceRatio.RawValue, "getReservePriceRatio")
}

// === RocketDAOProtocolSettingsDeposit ===

// Check if deposits are currently enabled
func (c *ProtocolDaoSettings) GetPoolDepositEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.DepositContract, &c.Details.Deposit.IsDepositingEnabled, "getDepositEnabled")
}

// Check if deposit assignments are currently enabled
func (c *ProtocolDaoSettings) GetAssignPoolDepositsEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.DepositContract, &c.Details.Deposit.AreDepositAssignmentsEnabled, "getAssignDepositsEnabled")
}

// Get the minimum deposit to the deposit pool
func (c *ProtocolDaoSettings) GetMinimumPoolDeposit(mc *batch.MultiCaller) {
	core.AddCall(mc, c.DepositContract, &c.Details.Deposit.MinimumDeposit, "getMinimumDeposit")
}

// Get the maximum size of the deposit pool
func (c *ProtocolDaoSettings) GetMaximumDepositPoolSize(mc *batch.MultiCaller) {
	core.AddCall(mc, c.DepositContract, &c.Details.Deposit.MaximumDepositPoolSize, "getMaximumDepositPoolSize")
}

// Get the total maximum assignments per deposit transaction, including socialized deposits
func (c *ProtocolDaoSettings) GetMaximumPoolDepositAssignments(mc *batch.MultiCaller) {
	core.AddCall(mc, c.DepositContract, &c.Details.Deposit.MaximumAssignmentsPerDeposit.RawValue, "getMaximumDepositAssignments")
}

// Get the number of "socialized" assignments for a pool deposit - these are assignments that always occur if the pool has enough ETH to support them, regardless of deposit size
func (c *ProtocolDaoSettings) GetMaximumPoolDepositSocialisedAssignments(mc *batch.MultiCaller) {
	core.AddCall(mc, c.DepositContract, &c.Details.Deposit.MaximumSocialisedAssignmentsPerDeposit.RawValue, "getMaximumDepositSocialisedAssignments")
}

// Get the fee that is applied to new pool deposits, as a fraction
func (c *ProtocolDaoSettings) GetDepositFee(mc *batch.MultiCaller) {
	core.AddCall(mc, c.DepositContract, &c.Details.Deposit.DepositFee.RawValue, "getDepositFee")
}

// === RocketDAOProtocolSettingsInflation ===

// Get the RPL inflation rate per interval
func (c *ProtocolDaoSettings) GetInflationIntervalRate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.InflationContract, &c.Details.Inflation.IntervalRate.RawValue, "getInflationIntervalRate")
}

// Get the RPL inflation start time
func (c *ProtocolDaoSettings) GetInflationIntervalStartTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.InflationContract, &c.Details.Inflation.StartTime.RawValue, "getInflationIntervalStartTime")
}

// === RocketDAOProtocolSettingsMinipool ===

// Check if minipool withdrawable event submissions are currently enabled
func (c *ProtocolDaoSettings) GetSubmitWithdrawableEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.IsSubmitWithdrawableEnabled, "getSubmitWithdrawableEnabled")
}

// Get the timeout period, in seconds, for prelaunch minipools to launch
func (c *ProtocolDaoSettings) GetMinipoolLaunchTimeout(mc *batch.MultiCaller) {
	core.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.LaunchTimeout.RawValue, "getLaunchTimeout")
}

// Check if minipool bond reductions are currently enabled
func (c *ProtocolDaoSettings) GetBondReductionEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.IsBondReductionEnabled, "getBondReductionEnabled")
}

// Get the limit on the total number of active minipools (non-finalized)
func (c *ProtocolDaoSettings) GetMaximumCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.MaximumCount.RawValue, "getMaximumCount")
}

// Gets the amount of time that must pass once someone calls beginUserDistribute() before the users can distribute a minipool
func (c *ProtocolDaoSettings) GetUserDistributeWindowStart(mc *batch.MultiCaller) {
	core.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.UserDistributeWindowStart.RawValue, "getUserDistributeWindowStart")
}

// Gets the amount of time the users have once UserDistributeWindowStart has passed to distribute a minipool before it expires
func (c *ProtocolDaoSettings) GetUserDistributeWindowLength(mc *batch.MultiCaller) {
	core.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.UserDistributeWindowLength.RawValue, "getUserDistributeWindowLength")
}

// === RocketDAOProtocolSettingsNetwork ===

// Get the threshold of Oracle DAO nodes that must reach consensus on oracle data to commit it
func (c *ProtocolDaoSettings) GetOracleDaoConsensusThreshold(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.OracleDaoConsensusThreshold.RawValue, "getNodeConsensusThreshold")
}

// Get the threshold of Oracle DAO nodes that must reach consensus on a penalty to apply it
func (c *ProtocolDaoSettings) GetNodePenaltyThreshold(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.NodePenaltyThreshold.RawValue, "getNodePenaltyThreshold")
}

// Get the fraction of a minipool's total node bond to penalize for each penalty
func (c *ProtocolDaoSettings) GetPerPenaltyRate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.PerPenaltyRate.RawValue, "getPerPenaltyRate")
}

// Check if network balance submission is enabled
func (c *ProtocolDaoSettings) GetSubmitBalancesEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.IsSubmitBalancesEnabled, "getSubmitBalancesEnabled")
}

// Get the frequency, in blocks, at which network balances should be submitted by the Oracle DAO
func (c *ProtocolDaoSettings) GetSubmitBalancesFrequency(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.SubmitBalancesFrequency.RawValue, "getSubmitBalancesFrequency")
}

// Check if network price submission is enabled
func (c *ProtocolDaoSettings) GetSubmitPricesEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.IsSubmitPricesEnabled, "getSubmitPricesEnabled")
}

// Get the frequency, in blocks, at which network prices should be submitted by the Oracle DAO
func (c *ProtocolDaoSettings) GetSubmitPricesFrequency(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.SubmitPricesFrequency.RawValue, "getSubmitPricesFrequency")
}

// Get the minimum node commission rate
func (c *ProtocolDaoSettings) GetMinimumNodeFee(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.MinimumNodeFee.RawValue, "getMinimumNodeFee")
}

// Get the target node commission rate
func (c *ProtocolDaoSettings) GetTargetNodeFee(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.TargetNodeFee.RawValue, "getTargetNodeFee")
}

// Get the maximum node commission rate
func (c *ProtocolDaoSettings) GetMaximumNodeFee(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.MaximumNodeFee.RawValue, "getMaximumNodeFee")
}

// Get the range of node demand values to base fee calculations on
func (c *ProtocolDaoSettings) GetNodeFeeDemandRange(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.NodeFeeDemandRange, "getNodeFeeDemandRange")
}

// Get the target collateralization rate for the rETH NetworkContract as a fraction
func (c *ProtocolDaoSettings) GetTargetRethCollateralRate(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.TargetRethCollateralRate.RawValue, "getTargetRethCollateralRate")
}

// Get the delay on pool deposits
func (c *ProtocolDaoSettings) GetRethDepositDelay(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.RethDepositDelay.RawValue, "getRethDepositDelay")
}

// Check if rewards submissions are enabled
func (c *ProtocolDaoSettings) GetSubmitRewardsEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NetworkContract, &c.Details.Network.IsSubmitRewardsEnabled, "getSubmitRewardsEnabled")
}

// === RocketDAOProtocolSettingsNode ===

// Check if node registration is currently enabled
func (c *ProtocolDaoSettings) GetNodeRegistrationEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NodeContract, &c.Details.Node.IsRegistrationEnabled, "getRegistrationEnabled")
}

// Check if smoothing pool registration is currently enabled
func (c *ProtocolDaoSettings) GetSmoothingPoolRegistrationEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NodeContract, &c.Details.Node.IsSmoothingPoolRegistrationEnabled, "getSmoothingPoolRegistrationEnabled")
}

// Check if node deposits are currently enabled
func (c *ProtocolDaoSettings) GetNodeDepositEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NodeContract, &c.Details.Node.IsDepositingEnabled, "getDepositEnabled")
}

// Check if creating vacant minipools is currently enabled
func (c *ProtocolDaoSettings) GetVacantMinipoolsEnabled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NodeContract, &c.Details.Node.AreVacantMinipoolsEnabled, "getVacantMinipoolsEnabled")
}

// Get the minimum RPL stake per minipool as a fraction of assigned user ETH
func (c *ProtocolDaoSettings) GetMinimumPerMinipoolStake(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NodeContract, &c.Details.Node.MinimumPerMinipoolStake.RawValue, "getMinimumPerMinipoolStake")
}

// Get the maximum RPL stake per minipool as a fraction of assigned user ETH
func (c *ProtocolDaoSettings) GetMaximumPerMinipoolStake(mc *batch.MultiCaller) {
	core.AddCall(mc, c.NodeContract, &c.Details.Node.MaximumPerMinipoolStake.RawValue, "getMaximumPerMinipoolStake")
}

// === RocketDAOProtocolSettingsRewards ===

// Get the rewards interval time
func (c *ProtocolDaoSettings) GetRewardsIntervalTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.RewardsContract, &c.Details.Rewards.IntervalTime.RawValue, "getRewardsClaimIntervalTime")
}

// === Universal ===

// Get all basic details
func (c *ProtocolDaoSettings) GetAllDetails(mc *batch.MultiCaller) {
	// Auction
	c.GetCreateAuctionLotEnabled(mc)
	c.GetBidOnAuctionLotEnabled(mc)
	c.GetAuctionLotMinimumEthValue(mc)
	c.GetAuctionLotMaximumEthValue(mc)
	c.GetAuctionLotDuration(mc)
	c.GetAuctionLotStartingPriceRatio(mc)
	c.GetAuctionLotReservePriceRatio(mc)

	// Deposit
	c.GetPoolDepositEnabled(mc)
	c.GetAssignPoolDepositsEnabled(mc)
	c.GetMinimumPoolDeposit(mc)
	c.GetMaximumDepositPoolSize(mc)
	c.GetMaximumPoolDepositAssignments(mc)
	c.GetMaximumPoolDepositSocialisedAssignments(mc)
	c.GetDepositFee(mc)

	// Inflation
	c.GetInflationIntervalRate(mc)
	c.GetInflationIntervalStartTime(mc)

	// Minipool
	c.GetSubmitWithdrawableEnabled(mc)
	c.GetMinipoolLaunchTimeout(mc)
	c.GetBondReductionEnabled(mc)
	c.GetMaximumCount(mc)
	c.GetUserDistributeWindowStart(mc)
	c.GetUserDistributeWindowLength(mc)

	// Network
	c.GetOracleDaoConsensusThreshold(mc)
	c.GetNodePenaltyThreshold(mc)
	c.GetPerPenaltyRate(mc)
	c.GetSubmitBalancesEnabled(mc)
	c.GetSubmitBalancesFrequency(mc)
	c.GetSubmitPricesEnabled(mc)
	c.GetSubmitPricesFrequency(mc)
	c.GetMinimumNodeFee(mc)
	c.GetTargetNodeFee(mc)
	c.GetMaximumNodeFee(mc)
	c.GetNodeFeeDemandRange(mc)
	c.GetTargetRethCollateralRate(mc)
	c.GetRethDepositDelay(mc)
	c.GetSubmitRewardsEnabled(mc)

	// Node
	c.GetNodeRegistrationEnabled(mc)
	c.GetSmoothingPoolRegistrationEnabled(mc)
	c.GetNodeDepositEnabled(mc)
	c.GetVacantMinipoolsEnabled(mc)
	c.GetMinimumPerMinipoolStake(mc)
	c.GetMaximumPerMinipoolStake(mc)

	// Rewards
	c.GetRewardsIntervalTime(mc)
}

// ====================
// === Transactions ===
// ====================

// === RocketDAOProtocolSettingsAuction ===

// Set the create lot enabled flag
func (c *ProtocolDaoSettings) BootstrapCreateAuctionLotEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.create.enabled", value, opts)
}

// Set the create lot enabled flag
func (c *ProtocolDaoSettings) BootstrapBidOnAuctionLotEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.bidding.enabled", value, opts)
}

// Set the minimum ETH value for lots
func (c *ProtocolDaoSettings) BootstrapAuctionLotMinimumEthValue(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.value.minimum", value, opts)
}

// Set the maximum ETH value for lots
func (c *ProtocolDaoSettings) BootstrapAuctionLotMaximumEthValue(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.value.maximum", value, opts)
}

// Set the duration value for lots, in blocks
func (c *ProtocolDaoSettings) BootstrapAuctionLotDuration(value core.Parameter[uint64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.duration", value.RawValue, opts)
}

// Set the starting price ratio for lots
func (c *ProtocolDaoSettings) BootstrapAuctionLotStartingPriceRatio(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.price.start", value.RawValue, opts)
}

// Set the reserve price ratio for lots
func (c *ProtocolDaoSettings) BootstrapAuctionLotReservePriceRatio(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.price.reserve", value.RawValue, opts)
}

// === RocketDAOProtocolSettingsDeposit ===

// Set the deposit enabled flag
func (c *ProtocolDaoSettings) BootstrapPoolDepositEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.enabled", value, opts)
}

// Set the deposit assignments enabled flag
func (c *ProtocolDaoSettings) BootstrapAssignPoolDepositsEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.assign.enabled", value, opts)
}

// Set the minimum deposit amount
func (c *ProtocolDaoSettings) BootstrapMinimumPoolDeposit(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.minimum", value, opts)
}

// Set the maximum deposit pool size
func (c *ProtocolDaoSettings) BootstrapMaximumDepositPoolSize(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.pool.maximum", value, opts)
}

// Set the max assignments per deposit
func (c *ProtocolDaoSettings) BootstrapMaximumPoolDepositAssignments(value core.Parameter[uint64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.assign.maximum", value.RawValue, opts)
}

// Set the max socialised assignments per deposit
func (c *ProtocolDaoSettings) BootstrapMaximumSocialisedPoolDepositAssignments(value core.Parameter[uint64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.assign.socialised.maximum", value.RawValue, opts)
}

// Set the deposit fee
func (c *ProtocolDaoSettings) BootstrapDepositFee(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.fee", value.RawValue, opts)
}

// === RocketDAOProtocolSettingsInflation ===

// Set the RPL inflation rate per interval
func (c *ProtocolDaoSettings) BootstrapInflationIntervalRate(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsInflation, "rpl.inflation.interval.rate", value.RawValue, opts)
}

// Set the RPL inflation start time
func (c *ProtocolDaoSettings) BootstrapInflationIntervalStartTime(value core.Parameter[time.Time], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsInflation, "rpl.inflation.interval.start", value.RawValue, opts)
}

// === RocketDAOProtocolSettingsMinipool ===

// Set the flag for enabling minipool withdrawable event submissions
func (c *ProtocolDaoSettings) BootstrapSubmitWithdrawableEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.submit.withdrawable.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapMinipoolLaunchTimeout(value core.Parameter[time.Duration], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.launch.timeout", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapBondReductionEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.bond.reduction.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapMaximumMinipoolCount(value core.Parameter[uint64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.maximum.count", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapUserDistributeWindowStart(value core.Parameter[time.Duration], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.user.distribute.window.start", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapUserDistributeWindowLength(value core.Parameter[time.Duration], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.user.distribute.window.length", value.RawValue, opts)
}

// === RocketDAOProtocolSettingsNetwork ===

func (c *ProtocolDaoSettings) BootstrapOracleDaoConsensusThreshold(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.consensus.threshold", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapNodePenaltyThreshold(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.penalty.threshold", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapPerPenaltyRate(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.penalty.per.rate", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapSubmitBalancesEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.balances.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapSubmitBalancesFrequency(value core.Parameter[time.Duration], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.balances.frequency", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapSubmitPricesEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.prices.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapSubmitPricesFrequency(value core.Parameter[time.Duration], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.prices.frequency", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapMinimumNodeFee(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.minimum", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapTargetNodeFee(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.target", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapMaximumNodeFee(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.maximum", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapNodeFeeDemandRange(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.demand.range", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapTargetRethCollateralRate(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.reth.collateral.target", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapRethDepositDelay(value core.Parameter[time.Duration], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.reth.deposit.delay", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapSubmitRewardsEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.rewards.enabled", value, opts)
}

// === RocketDAOProtocolSettingsNode ===

func (c *ProtocolDaoSettings) BootstrapNodeRegistrationEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.registration.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapSmoothingPoolRegistrationEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.smoothing.pool.registration.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapNodeDepositEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.deposit.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapVacantMinipoolsEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.vacant.minipools.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapMinimumPerMinipoolStake(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.per.minipool.stake.minimum", value.RawValue, opts)
}

func (c *ProtocolDaoSettings) BootstrapMaximumPerMinipoolStake(value core.Parameter[float64], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.per.minipool.stake.maximum", value.RawValue, opts)
}

// === RocketDAOProtocolSettingsRewards ===

func (c *ProtocolDaoSettings) BootstrapRewardsIntervalTime(value core.Parameter[time.Duration], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.pdaoMgr.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsRewards, "rpl.rewards.claim.period.time", value.RawValue, opts)
}
