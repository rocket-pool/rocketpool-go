package oracle

import (
	"fmt"
	"math/big"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for Oracle DAO settings
type OracleDaoSettings struct {
	*OracleDaoSettingsDetails
	dnts_members   *core.Contract
	dnts_minipool  *core.Contract
	dnts_proposals *core.Contract
	dnts_rewards   *core.Contract

	rp      *rocketpool.RocketPool
	odaoMgr *OracleDaoManager
}

// Details for Oracle DAO settings
type OracleDaoSettingsDetails struct {
	// Member
	Member struct {
		Quorum                 *OracleDaoCompoundSetting[float64]       `json:"quorum"`
		RplBond                *OracleDaoUintSetting                    `json:"rplBond"`
		UnbondedMinipoolMax    *OracleDaoCompoundSetting[uint64]        `json:"unbondedMinipoolMax"`
		UnbondedMinipoolMinFee *OracleDaoCompoundSetting[float64]       `json:"unbondedMinipoolMinFee"`
		ChallengeCooldown      *OracleDaoCompoundSetting[time.Duration] `json:"challengeCooldown"`
		ChallengeWindow        *OracleDaoCompoundSetting[time.Duration] `json:"challengeWindow"`
		ChallengeCost          *OracleDaoUintSetting                    `json:"challengeCost"`
	} `json:"members"`

	// Minipool
	Minipool struct {
		ScrubPeriod                     *OracleDaoCompoundSetting[time.Duration] `json:"scrubPeriod"`
		ScrubQuorum                     *OracleDaoCompoundSetting[float64]       `json:"scrubQuorum"`
		PromotionScrubPeriod            *OracleDaoCompoundSetting[time.Duration] `json:"promotionScrubPeriod"`
		IsScrubPenaltyEnabled           *OracleDaoBoolSetting                    `json:"isScrubPenaltyEnabled"`
		BondReductionWindowStart        *OracleDaoCompoundSetting[time.Duration] `json:"bondReductionWindowStart"`
		BondReductionWindowLength       *OracleDaoCompoundSetting[time.Duration] `json:"bondReductionWindowLength"`
		BondReductionCancellationQuorum *OracleDaoCompoundSetting[float64]       `json:"bondReductionCancellationQuorum"`
	} `json:"minipools"`

	// Proposal
	Proposal struct {
		CooldownTime  *OracleDaoCompoundSetting[time.Duration] `json:"cooldownTime"`
		VoteTime      *OracleDaoCompoundSetting[time.Duration] `json:"voteTime"`
		VoteDelayTime *OracleDaoCompoundSetting[time.Duration] `json:"voteDelayTime"`
		ExecuteTime   *OracleDaoCompoundSetting[time.Duration] `json:"executeTime"`
		ActionTime    *OracleDaoCompoundSetting[time.Duration] `json:"actionTime"`
	} `json:"proposals"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new Oracle DAO settings binding
func newOracleDaoSettings(odaoMgr *OracleDaoManager) (*OracleDaoSettings, error) {
	// Get the contracts
	contracts, err := odaoMgr.rp.GetContracts([]rocketpool.ContractName{
		rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers,
		rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool,
		rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals,
		rocketpool.ContractName_RocketDAONodeTrustedSettingsRewards,
	}...)
	if err != nil {
		return nil, fmt.Errorf("error getting Oracle DAO settings contracts: %w", err)
	}

	s := &OracleDaoSettings{
		OracleDaoSettingsDetails: &OracleDaoSettingsDetails{},
		rp:                       odaoMgr.rp,
		odaoMgr:                  odaoMgr,

		dnts_members:   contracts[0],
		dnts_minipool:  contracts[1],
		dnts_proposals: contracts[2],
		dnts_rewards:   contracts[3],
	}

	// Member
	s.Member.Quorum = newCompoundSetting[float64](s.dnts_members, odaoMgr, "members.quorum")
	s.Member.RplBond = newUintSetting(s.dnts_members, odaoMgr, "members.rplbond")
	s.Member.UnbondedMinipoolMax = newCompoundSetting[uint64](s.dnts_members, odaoMgr, "members.minipool.unbonded.max")
	s.Member.UnbondedMinipoolMinFee = newCompoundSetting[float64](s.dnts_members, odaoMgr, "members.minipool.unbonded.min.fee")
	s.Member.ChallengeCooldown = newCompoundSetting[time.Duration](s.dnts_members, odaoMgr, "members.challenge.cooldown")
	s.Member.ChallengeWindow = newCompoundSetting[time.Duration](s.dnts_members, odaoMgr, "members.challenge.window")
	s.Member.ChallengeCost = newUintSetting(s.dnts_members, odaoMgr, "members.challenge.cost")

	// Minipool
	s.Minipool.ScrubPeriod = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, "minipool.scrub.period")
	s.Minipool.ScrubQuorum = newCompoundSetting[float64](s.dnts_minipool, odaoMgr, "minipool.scrub.quorum")
	s.Minipool.PromotionScrubPeriod = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, "minipool.promotion.scrub.period")
	s.Minipool.IsScrubPenaltyEnabled = newBoolSetting(s.dnts_minipool, odaoMgr, "minipool.scrub.penalty.enabled")
	s.Minipool.BondReductionWindowStart = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, "minipool.bond.reduction.window.start")
	s.Minipool.BondReductionWindowLength = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, "minipool.bond.reduction.window.length")
	s.Minipool.BondReductionCancellationQuorum = newCompoundSetting[float64](s.dnts_minipool, odaoMgr, "minipool.cancel.bond.reduction.quorum")

	// Proposal
	s.Proposal.CooldownTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, "proposal.cooldown.time")
	s.Proposal.VoteTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, "proposal.vote.time")
	s.Proposal.VoteDelayTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, "proposal.vote.delay.time")
	s.Proposal.ExecuteTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, "proposal.execute.time")
	s.Proposal.ActionTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, "proposal.action.time")

	return s, nil
}

// =============
// === Calls ===
// =============

func (s *OracleDaoSettings) GetAllDetails(mc *batch.MultiCaller) {
	// Member
	s.Member.Quorum.Get(mc)
	s.Member.RplBond.Get(mc)
	s.Member.UnbondedMinipoolMax.Get(mc)
	s.Member.UnbondedMinipoolMinFee.Get(mc)
	s.Member.ChallengeCooldown.Get(mc)
	s.Member.ChallengeWindow.Get(mc)
	s.Member.ChallengeCost.Get(mc)

	/// Minipool
	s.Minipool.ScrubPeriod.Get(mc)
	s.Minipool.ScrubQuorum.Get(mc)
	s.Minipool.PromotionScrubPeriod.Get(mc)
	s.Minipool.IsScrubPenaltyEnabled.AddToQuery(mc)
	s.Minipool.BondReductionWindowStart.Get(mc)
	s.Minipool.BondReductionWindowLength.Get(mc)
	s.Minipool.BondReductionCancellationQuorum.Get(mc)

	// Proposal
	s.Proposal.CooldownTime.Get(mc)
	s.Proposal.VoteTime.Get(mc)
	s.Proposal.VoteDelayTime.Get(mc)
	s.Proposal.ExecuteTime.Get(mc)
	s.Proposal.ActionTime.Get(mc)
}

// === RocketDAONodeTrustedSettingsRewards ===

// Get whether or not the provided rewards network is enabled
func (c *OracleDaoSettings) GetNetworkEnabled(mc *batch.MultiCaller, enabled_Out *bool, network uint64) {
	core.AddCall(mc, c.dnts_rewards, enabled_Out, "getNetworkEnabled", big.NewInt(0).SetUint64(network))
}
