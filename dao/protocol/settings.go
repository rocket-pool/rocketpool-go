package protocol

import (
	"fmt"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for Protocol DAO settings
type ProtocolDaoSettings struct {
	*ProtocolDaoSettingsDetails
	dps_auction   *core.Contract
	dps_deposit   *core.Contract
	dps_inflation *core.Contract
	dps_minipool  *core.Contract
	dps_network   *core.Contract
	dps_node      *core.Contract
	dps_rewards   *core.Contract

	rp      *rocketpool.RocketPool
	pdaoMgr *ProtocolDaoManager
}

// Details for Protocol DAO settings
type ProtocolDaoSettingsDetails struct {
	Auction struct {
		IsCreateLotEnabled    *ProtocolDaoBoolSetting              `json:"isCreateLotEnabled"`
		IsBidOnLotEnabled     *ProtocolDaoBoolSetting              `json:"isBidOnLotEnabled"`
		LotMinimumEthValue    *ProtocolDaoUintSetting              `json:"lotMinimumEthValue"`
		LotMaximumEthValue    *ProtocolDaoUintSetting              `json:"lotMaximumEthValue"`
		LotDuration           *ProtocolDaoCompoundSetting[uint64]  `json:"lotDuration"`
		LotStartingPriceRatio *ProtocolDaoCompoundSetting[float64] `json:"lotStartingPriceRatio"`
		LotReservePriceRatio  *ProtocolDaoCompoundSetting[float64] `json:"lotReservePriceRatio"`
	} `json:"auction"`

	Deposit struct {
		IsDepositingEnabled                    *ProtocolDaoBoolSetting              `json:"isDepositingEnabled"`
		AreDepositAssignmentsEnabled           *ProtocolDaoBoolSetting              `json:"areDepositAssignmentsEnabled"`
		MinimumDeposit                         *ProtocolDaoUintSetting              `json:"minimumDeposit"`
		MaximumDepositPoolSize                 *ProtocolDaoUintSetting              `json:"maximumDepositPoolSize"`
		MaximumAssignmentsPerDeposit           *ProtocolDaoCompoundSetting[uint64]  `json:"maximumAssignmentsPerDeposit"`
		MaximumSocialisedAssignmentsPerDeposit *ProtocolDaoCompoundSetting[uint64]  `json:"maximumSocialisedAssignmentsPerDeposit"`
		DepositFee                             *ProtocolDaoCompoundSetting[float64] `json:"depositFee"`
	} `json:"deposit"`

	Inflation struct {
		IntervalRate *ProtocolDaoCompoundSetting[float64]   `json:"intervalRate"`
		StartTime    *ProtocolDaoCompoundSetting[time.Time] `json:"startTime"`
	} `json:"inflation"`

	Minipool struct {
		IsSubmitWithdrawableEnabled *ProtocolDaoBoolSetting                    `json:"isSubmitWithdrawableEnabled"`
		LaunchTimeout               *ProtocolDaoCompoundSetting[time.Duration] `json:"launchTimeout"`
		IsBondReductionEnabled      *ProtocolDaoBoolSetting                    `json:"isBondReductionEnabled"`
		MaximumCount                *ProtocolDaoCompoundSetting[uint64]        `json:"maximumCount"`
		UserDistributeWindowStart   *ProtocolDaoCompoundSetting[time.Duration] `json:"userDistributeWindowStart"`
		UserDistributeWindowLength  *ProtocolDaoCompoundSetting[time.Duration] `json:"userDistributeWindowLength"`
	} `json:"minipool"`

	Network struct {
		OracleDaoConsensusThreshold *ProtocolDaoCompoundSetting[float64]       `json:"oracleDaoConsensusThreshold"`
		NodePenaltyThreshold        *ProtocolDaoCompoundSetting[float64]       `json:"nodePenaltyThreshold"`
		PerPenaltyRate              *ProtocolDaoCompoundSetting[float64]       `json:"perPenaltyRate"`
		IsSubmitBalancesEnabled     *ProtocolDaoBoolSetting                    `json:"isSubmitBalancesEnabled"`
		SubmitBalancesFrequency     *ProtocolDaoCompoundSetting[time.Duration] `json:"submitBalancesFrequency"`
		IsSubmitPricesEnabled       *ProtocolDaoBoolSetting                    `json:"isSubmitPricesEnabled"`
		SubmitPricesFrequency       *ProtocolDaoCompoundSetting[time.Duration] `json:"submitPricesFrequency"`
		MinimumNodeFee              *ProtocolDaoCompoundSetting[float64]       `json:"minimumNodeFee"`
		TargetNodeFee               *ProtocolDaoCompoundSetting[float64]       `json:"targetNodeFee"`
		MaximumNodeFee              *ProtocolDaoCompoundSetting[float64]       `json:"maximumNodeFee"`
		NodeFeeDemandRange          *ProtocolDaoUintSetting                    `json:"nodeFeeDemandRange"`
		TargetRethCollateralRate    *ProtocolDaoCompoundSetting[float64]       `json:"targetRethCollateralRate"`
		RethDepositDelay            *ProtocolDaoCompoundSetting[time.Duration] `json:"rethDepositDelay"`
		IsSubmitRewardsEnabled      *ProtocolDaoBoolSetting                    `json:"isSubmitRewardsEnabled"`
	} `json:"network"`

	Node struct {
		IsRegistrationEnabled              *ProtocolDaoBoolSetting              `json:"isRegistrationEnabled"`
		IsSmoothingPoolRegistrationEnabled *ProtocolDaoBoolSetting              `json:"isSmoothingPoolRegistrationEnabled"`
		IsDepositingEnabled                *ProtocolDaoBoolSetting              `json:"isDepositingEnabled"`
		AreVacantMinipoolsEnabled          *ProtocolDaoBoolSetting              `json:"areVacantMinipoolsEnabled"`
		MinimumPerMinipoolStake            *ProtocolDaoCompoundSetting[float64] `json:"minimumPerMinipoolStake"`
		MaximumPerMinipoolStake            *ProtocolDaoCompoundSetting[float64] `json:"maximumPerMinipoolStake"`
	} `json:"node"`

	Rewards struct {
		IntervalTime *ProtocolDaoCompoundSetting[time.Duration] `json:"intervalTime"`
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

	s := &ProtocolDaoSettings{
		ProtocolDaoSettingsDetails: &ProtocolDaoSettingsDetails{},
		rp:                         pdaoMgr.rp,
		pdaoMgr:                    pdaoMgr,

		dps_auction:   contracts[0],
		dps_deposit:   contracts[1],
		dps_inflation: contracts[2],
		dps_minipool:  contracts[3],
		dps_network:   contracts[4],
		dps_node:      contracts[5],
		dps_rewards:   contracts[6],
	}

	// Auction
	s.Auction.IsCreateLotEnabled = newBoolSetting(s.dps_auction, pdaoMgr, "auction.lot.create.enabled")
	s.Auction.IsBidOnLotEnabled = newBoolSetting(s.dps_auction, pdaoMgr, "auction.lot.bidding.enabled")
	s.Auction.LotMinimumEthValue = newUintSetting(s.dps_auction, pdaoMgr, "auction.lot.value.minimum")
	s.Auction.LotMaximumEthValue = newUintSetting(s.dps_auction, pdaoMgr, "auction.lot.value.maximum")
	s.Auction.LotDuration = newCompoundSetting[uint64](s.dps_auction, pdaoMgr, "auction.lot.duration")
	s.Auction.LotStartingPriceRatio = newCompoundSetting[float64](s.dps_auction, pdaoMgr, "auction.price.start")
	s.Auction.LotReservePriceRatio = newCompoundSetting[float64](s.dps_auction, pdaoMgr, "auction.price.reserve")

	// Deposit
	s.Deposit.IsDepositingEnabled = newBoolSetting(s.dps_deposit, pdaoMgr, "deposit.enabled")
	s.Deposit.AreDepositAssignmentsEnabled = newBoolSetting(s.dps_deposit, pdaoMgr, "deposit.assign.enabled")
	s.Deposit.MinimumDeposit = newUintSetting(s.dps_deposit, pdaoMgr, "deposit.minimum")
	s.Deposit.MaximumDepositPoolSize = newUintSetting(s.dps_deposit, pdaoMgr, "deposit.pool.maximum")
	s.Deposit.MaximumAssignmentsPerDeposit = newCompoundSetting[uint64](s.dps_deposit, pdaoMgr, "deposit.assign.maximum")
	s.Deposit.MaximumSocialisedAssignmentsPerDeposit = newCompoundSetting[uint64](s.dps_deposit, pdaoMgr, "deposit.assign.socialised.maximum")
	s.Deposit.DepositFee = newCompoundSetting[float64](s.dps_deposit, pdaoMgr, "deposit.fee")

	// Inflation
	s.Inflation.IntervalRate = newCompoundSetting[float64](s.dps_inflation, pdaoMgr, "rpl.inflation.interval.rate")
	s.Inflation.StartTime = newCompoundSetting[time.Time](s.dps_inflation, pdaoMgr, "rpl.inflation.interval.start")

	// Minipool
	s.Minipool.IsSubmitWithdrawableEnabled = newBoolSetting(s.dps_minipool, pdaoMgr, "minipool.submit.withdrawable.enabled")
	s.Minipool.LaunchTimeout = newCompoundSetting[time.Duration](s.dps_minipool, pdaoMgr, "minipool.launch.timeout")
	s.Minipool.IsBondReductionEnabled = newBoolSetting(s.dps_minipool, pdaoMgr, "minipool.bond.reduction.enabled")
	s.Minipool.MaximumCount = newCompoundSetting[uint64](s.dps_minipool, pdaoMgr, "minipool.maximum.count")
	s.Minipool.UserDistributeWindowStart = newCompoundSetting[time.Duration](s.dps_minipool, pdaoMgr, "minipool.user.distribute.window.start")
	s.Minipool.UserDistributeWindowLength = newCompoundSetting[time.Duration](s.dps_minipool, pdaoMgr, "minipool.user.distribute.window.length")

	// Network
	s.Network.OracleDaoConsensusThreshold = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.consensus.threshold")
	s.Network.NodePenaltyThreshold = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.penalty.threshold")
	s.Network.PerPenaltyRate = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.penalty.per.rate")
	s.Network.IsSubmitBalancesEnabled = newBoolSetting(s.dps_network, pdaoMgr, "network.submit.balances.enabled")
	s.Network.SubmitBalancesFrequency = newCompoundSetting[time.Duration](s.dps_network, pdaoMgr, "network.submit.balances.frequency")
	s.Network.IsSubmitPricesEnabled = newBoolSetting(s.dps_network, pdaoMgr, "network.submit.prices.enabled")
	s.Network.SubmitPricesFrequency = newCompoundSetting[time.Duration](s.dps_network, pdaoMgr, "network.submit.prices.frequency")
	s.Network.MinimumNodeFee = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.node.fee.minimum")
	s.Network.TargetNodeFee = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.node.fee.target")
	s.Network.MaximumNodeFee = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.node.fee.maximum")
	s.Network.NodeFeeDemandRange = newUintSetting(s.dps_network, pdaoMgr, "network.node.fee.demand.range")
	s.Network.TargetRethCollateralRate = newCompoundSetting[float64](s.dps_network, pdaoMgr, "network.reth.collateral.target")
	s.Network.RethDepositDelay = newCompoundSetting[time.Duration](s.dps_network, pdaoMgr, "network.reth.deposit.delay")
	s.Network.IsSubmitRewardsEnabled = newBoolSetting(s.dps_network, pdaoMgr, "network.submit.rewards.enabled")

	// Node
	s.Node.IsRegistrationEnabled = newBoolSetting(s.dps_node, pdaoMgr, "node.registration.enabled")
	s.Node.IsSmoothingPoolRegistrationEnabled = newBoolSetting(s.dps_node, pdaoMgr, "node.smoothing.pool.registration.enabled")
	s.Node.IsDepositingEnabled = newBoolSetting(s.dps_node, pdaoMgr, "node.deposit.enabled")
	s.Node.AreVacantMinipoolsEnabled = newBoolSetting(s.dps_node, pdaoMgr, "node.vacant.minipools.enabled")
	s.Node.MinimumPerMinipoolStake = newCompoundSetting[float64](s.dps_node, pdaoMgr, "node.per.minipool.stake.minimum")
	s.Node.MaximumPerMinipoolStake = newCompoundSetting[float64](s.dps_node, pdaoMgr, "node.per.minipool.stake.maximum")

	// Rewards
	s.Rewards.IntervalTime = newCompoundSetting[time.Duration](s.dps_rewards, pdaoMgr, "rpl.rewards.claim.period.time")

	return s, nil
}

// =============
// === Calls ===
// =============

func (s *ProtocolDaoSettings) GetAllDetails(mc *batch.MultiCaller) {
	// Auction
	s.Auction.IsCreateLotEnabled.Get(mc)
	s.Auction.IsBidOnLotEnabled.Get(mc)
	s.Auction.LotMinimumEthValue.Get(mc)
	s.Auction.LotMaximumEthValue.Get(mc)
	s.Auction.LotDuration.Get(mc)
	s.Auction.LotStartingPriceRatio.Get(mc)
	s.Auction.LotReservePriceRatio.Get(mc)

	// Deposit
	s.Deposit.IsDepositingEnabled.Get(mc)
	s.Deposit.AreDepositAssignmentsEnabled.Get(mc)
	s.Deposit.MinimumDeposit.Get(mc)
	s.Deposit.MaximumDepositPoolSize.Get(mc)
	s.Deposit.MaximumAssignmentsPerDeposit.Get(mc)
	s.Deposit.MaximumSocialisedAssignmentsPerDeposit.Get(mc)
	s.Deposit.DepositFee.Get(mc)

	// Inflation
	s.Inflation.IntervalRate.Get(mc)
	s.Inflation.StartTime.Get(mc)

	// Minipool
	s.Minipool.IsSubmitWithdrawableEnabled.Get(mc)
	s.Minipool.LaunchTimeout.Get(mc)
	s.Minipool.IsBondReductionEnabled.Get(mc)
	s.Minipool.MaximumCount.Get(mc)
	s.Minipool.UserDistributeWindowStart.Get(mc)
	s.Minipool.UserDistributeWindowLength.Get(mc)

	// Network
	s.Network.OracleDaoConsensusThreshold.Get(mc)
	s.Network.NodePenaltyThreshold.Get(mc)
	s.Network.PerPenaltyRate.Get(mc)
	s.Network.IsSubmitBalancesEnabled.Get(mc)
	s.Network.SubmitBalancesFrequency.Get(mc)
	s.Network.IsSubmitPricesEnabled.Get(mc)
	s.Network.SubmitPricesFrequency.Get(mc)
	s.Network.MinimumNodeFee.Get(mc)
	s.Network.TargetNodeFee.Get(mc)
	s.Network.MaximumNodeFee.Get(mc)
	s.Network.NodeFeeDemandRange.Get(mc)
	s.Network.TargetRethCollateralRate.Get(mc)
	s.Network.RethDepositDelay.Get(mc)
	s.Network.IsSubmitRewardsEnabled.Get(mc)

	// Node
	s.Node.IsRegistrationEnabled.Get(mc)
	s.Node.IsSmoothingPoolRegistrationEnabled.Get(mc)
	s.Node.IsDepositingEnabled.Get(mc)
	s.Node.AreVacantMinipoolsEnabled.Get(mc)
	s.Node.MinimumPerMinipoolStake.Get(mc)
	s.Node.MaximumPerMinipoolStake.Get(mc)

	// Rewards
	s.Rewards.IntervalTime.Get(mc)
}
