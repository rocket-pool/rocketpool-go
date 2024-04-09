package protocol

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/rocket-pool/rocketpool-go/v2/core"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
)

// =====================
// === Setting Names ===
// =====================

type SettingName string

const (
	// Auction
	SettingName_Auction_IsCreateLotEnabled    SettingName = "auction.lot.create.enabled"
	SettingName_Auction_IsBidOnLotEnabled     SettingName = "auction.lot.bidding.enabled"
	SettingName_Auction_LotMinimumEthValue    SettingName = "auction.lot.value.minimum"
	SettingName_Auction_LotMaximumEthValue    SettingName = "auction.lot.value.maximum"
	SettingName_Auction_LotDuration           SettingName = "auction.lot.duration"
	SettingName_Auction_LotStartingPriceRatio SettingName = "auction.price.start"
	SettingName_Auction_LotReservePriceRatio  SettingName = "auction.price.reserve"

	// Deposit
	SettingName_Deposit_IsDepositingEnabled                    SettingName = "deposit.enabled"
	SettingName_Deposit_AreDepositAssignmentsEnabled           SettingName = "deposit.assign.enabled"
	SettingName_Deposit_MinimumDeposit                         SettingName = "deposit.minimum"
	SettingName_Deposit_MaximumDepositPoolSize                 SettingName = "deposit.pool.maximum"
	SettingName_Deposit_MaximumAssignmentsPerDeposit           SettingName = "deposit.assign.maximum"
	SettingName_Deposit_MaximumSocialisedAssignmentsPerDeposit SettingName = "deposit.assign.socialised.maximum"
	SettingName_Deposit_DepositFee                             SettingName = "deposit.fee"

	// Inflation
	SettingName_Inflation_IntervalRate SettingName = "rpl.inflation.interval.rate"
	SettingName_Inflation_StartTime    SettingName = "rpl.inflation.interval.start"

	// Minipool
	SettingName_Minipool_IsSubmitWithdrawableEnabled SettingName = "minipool.submit.withdrawable.enabled"
	SettingName_Minipool_LaunchTimeout               SettingName = "minipool.launch.timeout"
	SettingName_Minipool_IsBondReductionEnabled      SettingName = "minipool.bond.reduction.enabled"
	SettingName_Minipool_MaximumCount                SettingName = "minipool.maximum.count"
	SettingName_Minipool_UserDistributeWindowStart   SettingName = "minipool.user.distribute.window.start"
	SettingName_Minipool_UserDistributeWindowLength  SettingName = "minipool.user.distribute.window.length"

	// Network
	SettingName_Network_OracleDaoConsensusThreshold SettingName = "network.consensus.threshold"
	SettingName_Network_NodePenaltyThreshold        SettingName = "network.penalty.threshold"
	SettingName_Network_PerPenaltyRate              SettingName = "network.penalty.per.rate"
	SettingName_Network_IsSubmitBalancesEnabled     SettingName = "network.submit.balances.enabled"
	SettingName_Network_SubmitBalancesFrequency     SettingName = "network.submit.balances.frequency"
	SettingName_Network_IsSubmitPricesEnabled       SettingName = "network.submit.prices.enabled"
	SettingName_Network_SubmitPricesFrequency       SettingName = "network.submit.prices.frequency"
	SettingName_Network_MinimumNodeFee              SettingName = "network.node.fee.minimum"
	SettingName_Network_TargetNodeFee               SettingName = "network.node.fee.target"
	SettingName_Network_MaximumNodeFee              SettingName = "network.node.fee.maximum"
	SettingName_Network_NodeFeeDemandRange          SettingName = "network.node.fee.demand.range"
	SettingName_Network_TargetRethCollateralRate    SettingName = "network.reth.collateral.target"
	SettingName_Network_IsSubmitRewardsEnabled      SettingName = "network.submit.rewards.enabled"

	// Node
	SettingName_Node_IsRegistrationEnabled              SettingName = "node.registration.enabled"
	SettingName_Node_IsSmoothingPoolRegistrationEnabled SettingName = "node.smoothing.pool.registration.enabled"
	SettingName_Node_IsDepositingEnabled                SettingName = "node.deposit.enabled"
	SettingName_Node_AreVacantMinipoolsEnabled          SettingName = "node.vacant.minipools.enabled"
	SettingName_Node_MinimumPerMinipoolStake            SettingName = "node.per.minipool.stake.minimum"
	SettingName_Node_MaximumPerMinipoolStake            SettingName = "node.per.minipool.stake.maximum"

	// Proposals
	SettingName_Proposals_VotePhase1Time      SettingName = "proposal.vote.phase1.time"
	SettingName_Proposals_VotePhase2Time      SettingName = "proposal.vote.phase2.time"
	SettingName_Proposals_VoteDelayTime       SettingName = "proposal.vote.delay.time"
	SettingName_Proposals_ExecuteTime         SettingName = "proposal.execute.time"
	SettingName_Proposals_ProposalBond        SettingName = "proposal.bond"
	SettingName_Proposals_ChallengeBond       SettingName = "proposal.challenge.bond"
	SettingName_Proposals_ChallengePeriod     SettingName = "proposal.challenge.period"
	SettingName_Proposals_ProposalQuorum      SettingName = "proposal.quorum"
	SettingName_Proposals_ProposalVetoQuorum  SettingName = "proposal.veto.quorum"
	SettingName_Proposals_ProposalMaxBlockAge SettingName = "proposal.max.block.age"

	// Rewards
	SettingName_Rewards_IntervalPeriods SettingName = "rewards.claimsperiods"

	// Security
	SettingName_Security_MembersQuorum       SettingName = "members.quorum"
	SettingName_Security_MembersLeaveTime    SettingName = "members.leave.time"
	SettingName_Security_ProposalVoteTime    SettingName = "proposal.vote.time"
	SettingName_Security_ProposalExecuteTime SettingName = "proposal.execute.time"
	SettingName_Security_ProposalActionTime  SettingName = "proposal.action.time"
)

// ===============
// === Structs ===
// ===============

// Wrapper for a settings category, with all of its settings
type SettingsCategory struct {
	ContractName rocketpool.ContractName
	BoolSettings []IProtocolDaoSetting[bool]
	UintSettings []IProtocolDaoSetting[*big.Int]
}

// Binding for Protocol DAO settings
type ProtocolDaoSettings struct {
	Auction struct {
		IsCreateLotEnabled    *ProtocolDaoBoolSetting
		IsBidOnLotEnabled     *ProtocolDaoBoolSetting
		LotMinimumEthValue    *ProtocolDaoUintSetting
		LotMaximumEthValue    *ProtocolDaoUintSetting
		LotDuration           *ProtocolDaoCompoundSetting[time.Duration]
		LotStartingPriceRatio *ProtocolDaoCompoundSetting[float64]
		LotReservePriceRatio  *ProtocolDaoCompoundSetting[float64]
	}

	Deposit struct {
		IsDepositingEnabled                    *ProtocolDaoBoolSetting
		AreDepositAssignmentsEnabled           *ProtocolDaoBoolSetting
		MinimumDeposit                         *ProtocolDaoUintSetting
		MaximumDepositPoolSize                 *ProtocolDaoUintSetting
		MaximumAssignmentsPerDeposit           *ProtocolDaoCompoundSetting[uint64]
		MaximumSocialisedAssignmentsPerDeposit *ProtocolDaoCompoundSetting[uint64]
		DepositFee                             *ProtocolDaoCompoundSetting[float64]
	}

	Inflation struct {
		IntervalRate *ProtocolDaoCompoundSetting[float64]
		StartTime    *ProtocolDaoCompoundSetting[time.Time]
	}

	Minipool struct {
		IsSubmitWithdrawableEnabled *ProtocolDaoBoolSetting
		LaunchTimeout               *ProtocolDaoCompoundSetting[time.Duration]
		IsBondReductionEnabled      *ProtocolDaoBoolSetting
		MaximumCount                *ProtocolDaoCompoundSetting[uint64]
		UserDistributeWindowStart   *ProtocolDaoCompoundSetting[time.Duration]
		UserDistributeWindowLength  *ProtocolDaoCompoundSetting[time.Duration]
	}

	Network struct {
		OracleDaoConsensusThreshold *ProtocolDaoCompoundSetting[float64]
		NodePenaltyThreshold        *ProtocolDaoCompoundSetting[float64]
		PerPenaltyRate              *ProtocolDaoCompoundSetting[float64]
		IsSubmitBalancesEnabled     *ProtocolDaoBoolSetting
		SubmitBalancesFrequency     *ProtocolDaoCompoundSetting[time.Duration]
		IsSubmitPricesEnabled       *ProtocolDaoBoolSetting
		SubmitPricesFrequency       *ProtocolDaoCompoundSetting[time.Duration]
		MinimumNodeFee              *ProtocolDaoCompoundSetting[float64]
		TargetNodeFee               *ProtocolDaoCompoundSetting[float64]
		MaximumNodeFee              *ProtocolDaoCompoundSetting[float64]
		NodeFeeDemandRange          *ProtocolDaoUintSetting
		TargetRethCollateralRate    *ProtocolDaoCompoundSetting[float64]
		IsSubmitRewardsEnabled      *ProtocolDaoBoolSetting
	}

	Node struct {
		IsRegistrationEnabled              *ProtocolDaoBoolSetting
		IsSmoothingPoolRegistrationEnabled *ProtocolDaoBoolSetting
		IsDepositingEnabled                *ProtocolDaoBoolSetting
		AreVacantMinipoolsEnabled          *ProtocolDaoBoolSetting
		MinimumPerMinipoolStake            *ProtocolDaoCompoundSetting[float64]
		MaximumPerMinipoolStake            *ProtocolDaoCompoundSetting[float64]
	}

	Proposals struct {
		VotePhase1Time      *ProtocolDaoCompoundSetting[time.Duration]
		VotePhase2Time      *ProtocolDaoCompoundSetting[time.Duration]
		VoteDelayTime       *ProtocolDaoCompoundSetting[time.Duration]
		ExecuteTime         *ProtocolDaoCompoundSetting[time.Duration]
		ProposalBond        *ProtocolDaoUintSetting
		ChallengeBond       *ProtocolDaoUintSetting
		ChallengePeriod     *ProtocolDaoCompoundSetting[time.Duration]
		ProposalQuorum      *ProtocolDaoCompoundSetting[float64]
		ProposalVetoQuorum  *ProtocolDaoCompoundSetting[float64]
		ProposalMaxBlockAge *ProtocolDaoCompoundSetting[uint64]
	}

	Rewards struct {
		IntervalPeriods *ProtocolDaoCompoundSetting[uint64]
	}

	Security struct {
		MembersQuorum       *ProtocolDaoCompoundSetting[float64]
		MembersLeaveTime    *ProtocolDaoCompoundSetting[time.Duration]
		ProposalVoteTime    *ProtocolDaoCompoundSetting[time.Duration]
		ProposalExecuteTime *ProtocolDaoCompoundSetting[time.Duration]
		ProposalActionTime  *ProtocolDaoCompoundSetting[time.Duration]
	}

	// === Internal fields ===
	rp              *rocketpool.RocketPool
	pdaoMgr         *ProtocolDaoManager
	dps_auction     *core.Contract
	dps_deposit     *core.Contract
	dps_inflation   *core.Contract
	dps_minipool    *core.Contract
	dps_network     *core.Contract
	dps_node        *core.Contract
	dps_proposals   *core.Contract
	dps_rewards     *core.Contract
	dps_security    *core.Contract
	contractNameMap map[string]rocketpool.ContractName
}

// ====================
// === Constructors ===
// ====================

// Creates a new ProtocolDaoSettings binding
func newProtocolDaoSettings(pdaoMgr *ProtocolDaoManager) (*ProtocolDaoSettings, error) {
	// Get the contracts
	contractNames := []rocketpool.ContractName{
		rocketpool.ContractName_RocketDAOProtocolSettingsAuction,
		rocketpool.ContractName_RocketDAOProtocolSettingsDeposit,
		rocketpool.ContractName_RocketDAOProtocolSettingsInflation,
		rocketpool.ContractName_RocketDAOProtocolSettingsMinipool,
		rocketpool.ContractName_RocketDAOProtocolSettingsNetwork,
		rocketpool.ContractName_RocketDAOProtocolSettingsNode,
		rocketpool.ContractName_RocketDAOProtocolSettingsProposals,
		rocketpool.ContractName_RocketDAOProtocolSettingsRewards,
		rocketpool.ContractName_RocketDAOProtocolSettingsSecurity,
	}
	contracts, err := pdaoMgr.rp.GetContracts(contractNames...)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO settings contracts: %w", err)
	}

	s := &ProtocolDaoSettings{
		rp:      pdaoMgr.rp,
		pdaoMgr: pdaoMgr,

		dps_auction:   contracts[0],
		dps_deposit:   contracts[1],
		dps_inflation: contracts[2],
		dps_minipool:  contracts[3],
		dps_network:   contracts[4],
		dps_node:      contracts[5],
		dps_proposals: contracts[6],
		dps_rewards:   contracts[7],
		dps_security:  contracts[8],
	}
	s.contractNameMap = map[string]rocketpool.ContractName{
		"Auction":   contractNames[0],
		"Deposit":   contractNames[1],
		"Inflation": contractNames[2],
		"Minipool":  contractNames[3],
		"Network":   contractNames[4],
		"Node":      contractNames[5],
		"Proposals": contractNames[6],
		"Rewards":   contractNames[7],
		"Security":  contractNames[8],
	}

	// Auction
	s.Auction.IsCreateLotEnabled = newBoolSetting(s.dps_auction, pdaoMgr, SettingName_Auction_IsCreateLotEnabled)
	s.Auction.IsBidOnLotEnabled = newBoolSetting(s.dps_auction, pdaoMgr, SettingName_Auction_IsBidOnLotEnabled)
	s.Auction.LotMinimumEthValue = newUintSetting(s.dps_auction, pdaoMgr, SettingName_Auction_LotMinimumEthValue)
	s.Auction.LotMaximumEthValue = newUintSetting(s.dps_auction, pdaoMgr, SettingName_Auction_LotMaximumEthValue)
	s.Auction.LotDuration = newCompoundSetting[time.Duration](s.dps_auction, pdaoMgr, SettingName_Auction_LotDuration)
	s.Auction.LotStartingPriceRatio = newCompoundSetting[float64](s.dps_auction, pdaoMgr, SettingName_Auction_LotStartingPriceRatio)
	s.Auction.LotReservePriceRatio = newCompoundSetting[float64](s.dps_auction, pdaoMgr, SettingName_Auction_LotReservePriceRatio)

	// Deposit
	s.Deposit.IsDepositingEnabled = newBoolSetting(s.dps_deposit, pdaoMgr, SettingName_Deposit_IsDepositingEnabled)
	s.Deposit.AreDepositAssignmentsEnabled = newBoolSetting(s.dps_deposit, pdaoMgr, SettingName_Deposit_AreDepositAssignmentsEnabled)
	s.Deposit.MinimumDeposit = newUintSetting(s.dps_deposit, pdaoMgr, SettingName_Deposit_MinimumDeposit)
	s.Deposit.MaximumDepositPoolSize = newUintSetting(s.dps_deposit, pdaoMgr, SettingName_Deposit_MaximumDepositPoolSize)
	s.Deposit.MaximumAssignmentsPerDeposit = newCompoundSetting[uint64](s.dps_deposit, pdaoMgr, SettingName_Deposit_MaximumAssignmentsPerDeposit)
	s.Deposit.MaximumSocialisedAssignmentsPerDeposit = newCompoundSetting[uint64](s.dps_deposit, pdaoMgr, SettingName_Deposit_MaximumSocialisedAssignmentsPerDeposit)
	s.Deposit.DepositFee = newCompoundSetting[float64](s.dps_deposit, pdaoMgr, SettingName_Deposit_DepositFee)

	// Inflation
	s.Inflation.IntervalRate = newCompoundSetting[float64](s.dps_inflation, pdaoMgr, SettingName_Inflation_IntervalRate)
	s.Inflation.StartTime = newCompoundSetting[time.Time](s.dps_inflation, pdaoMgr, SettingName_Inflation_StartTime)

	// Minipool
	s.Minipool.IsSubmitWithdrawableEnabled = newBoolSetting(s.dps_minipool, pdaoMgr, SettingName_Minipool_IsSubmitWithdrawableEnabled)
	s.Minipool.LaunchTimeout = newCompoundSetting[time.Duration](s.dps_minipool, pdaoMgr, SettingName_Minipool_LaunchTimeout)
	s.Minipool.IsBondReductionEnabled = newBoolSetting(s.dps_minipool, pdaoMgr, SettingName_Minipool_IsBondReductionEnabled)
	s.Minipool.MaximumCount = newCompoundSetting[uint64](s.dps_minipool, pdaoMgr, SettingName_Minipool_MaximumCount)
	s.Minipool.UserDistributeWindowStart = newCompoundSetting[time.Duration](s.dps_minipool, pdaoMgr, SettingName_Minipool_UserDistributeWindowStart)
	s.Minipool.UserDistributeWindowLength = newCompoundSetting[time.Duration](s.dps_minipool, pdaoMgr, SettingName_Minipool_UserDistributeWindowLength)

	// Network
	s.Network.OracleDaoConsensusThreshold = newCompoundSetting[float64](s.dps_network, pdaoMgr, SettingName_Network_OracleDaoConsensusThreshold)
	s.Network.NodePenaltyThreshold = newCompoundSetting[float64](s.dps_network, pdaoMgr, SettingName_Network_NodePenaltyThreshold)
	s.Network.PerPenaltyRate = newCompoundSetting[float64](s.dps_network, pdaoMgr, SettingName_Network_PerPenaltyRate)
	s.Network.IsSubmitBalancesEnabled = newBoolSetting(s.dps_network, pdaoMgr, SettingName_Network_IsSubmitBalancesEnabled)
	s.Network.SubmitBalancesFrequency = newCompoundSetting[time.Duration](s.dps_network, pdaoMgr, SettingName_Network_SubmitBalancesFrequency)
	s.Network.IsSubmitPricesEnabled = newBoolSetting(s.dps_network, pdaoMgr, SettingName_Network_IsSubmitPricesEnabled)
	s.Network.SubmitPricesFrequency = newCompoundSetting[time.Duration](s.dps_network, pdaoMgr, SettingName_Network_SubmitPricesFrequency)
	s.Network.MinimumNodeFee = newCompoundSetting[float64](s.dps_network, pdaoMgr, SettingName_Network_MinimumNodeFee)
	s.Network.TargetNodeFee = newCompoundSetting[float64](s.dps_network, pdaoMgr, SettingName_Network_TargetNodeFee)
	s.Network.MaximumNodeFee = newCompoundSetting[float64](s.dps_network, pdaoMgr, SettingName_Network_MaximumNodeFee)
	s.Network.NodeFeeDemandRange = newUintSetting(s.dps_network, pdaoMgr, SettingName_Network_NodeFeeDemandRange)
	s.Network.TargetRethCollateralRate = newCompoundSetting[float64](s.dps_network, pdaoMgr, SettingName_Network_TargetRethCollateralRate)
	s.Network.IsSubmitRewardsEnabled = newBoolSetting(s.dps_network, pdaoMgr, SettingName_Network_IsSubmitRewardsEnabled)

	// Node
	s.Node.IsRegistrationEnabled = newBoolSetting(s.dps_node, pdaoMgr, SettingName_Node_IsRegistrationEnabled)
	s.Node.IsSmoothingPoolRegistrationEnabled = newBoolSetting(s.dps_node, pdaoMgr, SettingName_Node_IsSmoothingPoolRegistrationEnabled)
	s.Node.IsDepositingEnabled = newBoolSetting(s.dps_node, pdaoMgr, SettingName_Node_IsDepositingEnabled)
	s.Node.AreVacantMinipoolsEnabled = newBoolSetting(s.dps_node, pdaoMgr, SettingName_Node_AreVacantMinipoolsEnabled)
	s.Node.MinimumPerMinipoolStake = newCompoundSetting[float64](s.dps_node, pdaoMgr, SettingName_Node_MinimumPerMinipoolStake)
	s.Node.MaximumPerMinipoolStake = newCompoundSetting[float64](s.dps_node, pdaoMgr, SettingName_Node_MaximumPerMinipoolStake)

	// Proposals
	s.Proposals.VotePhase1Time = newCompoundSetting[time.Duration](s.dps_proposals, pdaoMgr, SettingName_Proposals_VotePhase1Time)
	s.Proposals.VotePhase2Time = newCompoundSetting[time.Duration](s.dps_proposals, pdaoMgr, SettingName_Proposals_VotePhase2Time)
	s.Proposals.VoteDelayTime = newCompoundSetting[time.Duration](s.dps_proposals, pdaoMgr, SettingName_Proposals_VoteDelayTime)
	s.Proposals.ExecuteTime = newCompoundSetting[time.Duration](s.dps_proposals, pdaoMgr, SettingName_Proposals_ExecuteTime)
	s.Proposals.ProposalBond = newUintSetting(s.dps_proposals, pdaoMgr, SettingName_Proposals_ProposalBond)
	s.Proposals.ChallengeBond = newUintSetting(s.dps_proposals, pdaoMgr, SettingName_Proposals_ChallengeBond)
	s.Proposals.ChallengePeriod = newCompoundSetting[time.Duration](s.dps_proposals, pdaoMgr, SettingName_Proposals_ChallengePeriod)
	s.Proposals.ProposalQuorum = newCompoundSetting[float64](s.dps_proposals, pdaoMgr, SettingName_Proposals_ProposalQuorum)
	s.Proposals.ProposalVetoQuorum = newCompoundSetting[float64](s.dps_proposals, pdaoMgr, SettingName_Proposals_ProposalVetoQuorum)
	s.Proposals.ProposalMaxBlockAge = newCompoundSetting[uint64](s.dps_proposals, pdaoMgr, SettingName_Proposals_ProposalMaxBlockAge)

	// Rewards
	s.Rewards.IntervalPeriods = newCompoundSetting[uint64](s.dps_rewards, pdaoMgr, SettingName_Rewards_IntervalPeriods)

	// Security
	s.Security.MembersQuorum = newCompoundSetting[float64](s.dps_security, pdaoMgr, SettingName_Security_MembersQuorum)
	s.Security.MembersLeaveTime = newCompoundSetting[time.Duration](s.dps_security, pdaoMgr, SettingName_Security_MembersLeaveTime)
	s.Security.ProposalVoteTime = newCompoundSetting[time.Duration](s.dps_security, pdaoMgr, SettingName_Security_ProposalVoteTime)
	s.Security.ProposalExecuteTime = newCompoundSetting[time.Duration](s.dps_security, pdaoMgr, SettingName_Security_ProposalExecuteTime)
	s.Security.ProposalActionTime = newCompoundSetting[time.Duration](s.dps_security, pdaoMgr, SettingName_Security_ProposalActionTime)

	return s, nil
}

// =============
// === Calls ===
// =============

// Get all of the settings, organized by the type used in proposals and boostraps
func (c *ProtocolDaoSettings) GetSettings() map[rocketpool.ContractName]SettingsCategory {
	catMap := map[rocketpool.ContractName]SettingsCategory{}

	settingsType := reflect.TypeOf(c)
	settingsVal := reflect.ValueOf(c)
	fieldCount := settingsType.NumField()
	for i := 0; i < fieldCount; i++ {
		categoryField := settingsType.Field(i)
		categoryFieldType := categoryField.Type

		// A container struct for settings by category
		if categoryFieldType.Kind() == reflect.Struct {
			// Get the contract name of this category
			name, exists := c.contractNameMap[categoryField.Name]
			if !exists {
				panic(fmt.Sprintf("Protocol DAO settings field named %s does not exist in the contract map.", name))
			}
			boolSettings := []IProtocolDaoSetting[bool]{}
			uintSettings := []IProtocolDaoSetting[*big.Int]{}

			// Get all of the settings in this cateogry
			categoryFieldVal := settingsVal.Field(i)
			settingCount := categoryFieldType.NumField()
			for j := 0; j < settingCount; j++ {
				setting := categoryFieldVal.Field(i).Interface()

				// Try bool settings
				boolSetting, isBoolSetting := setting.(IProtocolDaoSetting[bool])
				if isBoolSetting {
					boolSettings = append(boolSettings, boolSetting)
					continue
				}

				// Try uint settings
				uintSetting, isUintSetting := setting.(IProtocolDaoSetting[*big.Int])
				if isUintSetting {
					uintSettings = append(uintSettings, uintSetting)
				}
			}

			settingsCat := SettingsCategory{
				ContractName: name,
				BoolSettings: boolSettings,
				UintSettings: uintSettings,
			}
			catMap[name] = settingsCat
		}
	}

	return catMap
}
