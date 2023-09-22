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
	// Members
	Members struct {
		Quorum                 *OracleDaoCompoundSetting[float64]       `json:"quorum"`
		RplBond                *OracleDaoUintSetting                    `json:"rplBond"`
		UnbondedMinipoolMax    *OracleDaoCompoundSetting[uint64]        `json:"unbondedMinipoolMax"`
		UnbondedMinipoolMinFee *OracleDaoCompoundSetting[float64]       `json:"unbondedMinipoolMinFee"`
		ChallengeCooldown      *OracleDaoCompoundSetting[time.Duration] `json:"challengeCooldown"`
		ChallengeWindow        *OracleDaoCompoundSetting[time.Duration] `json:"challengeWindow"`
		ChallengeCost          *OracleDaoUintSetting                    `json:"challengeCost"`
	} `json:"members"`

	// Minipools
	Minipools struct {
		ScrubPeriod                     *OracleDaoCompoundSetting[time.Duration] `json:"scrubPeriod"`
		ScrubQuorum                     *OracleDaoCompoundSetting[float64]       `json:"scrubQuorum"`
		PromotionScrubPeriod            *OracleDaoCompoundSetting[time.Duration] `json:"promotionScrubPeriod"`
		IsScrubPenaltyEnabled           *OracleDaoBoolSetting                    `json:"isScrubPenaltyEnabled"`
		BondReductionWindowStart        *OracleDaoCompoundSetting[time.Duration] `json:"bondReductionWindowStart"`
		BondReductionWindowLength       *OracleDaoCompoundSetting[time.Duration] `json:"bondReductionWindowLength"`
		BondReductionCancellationQuorum *OracleDaoCompoundSetting[float64]       `json:"bondReductionCancellationQuorum"`
	} `json:"minipools"`

	// Proposals
	Proposals struct {
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

	// Member settings
	s.Members.Quorum = newCompoundSetting[float64](s.dnts_members, odaoMgr, "members.quorum")
	s.Members.RplBond = newUintSetting(s.dnts_members, odaoMgr, "members.rplbond")
	s.Members.UnbondedMinipoolMax = newCompoundSetting[uint64](s.dnts_members, odaoMgr, "members.minipool.unbonded.max")
	s.Members.UnbondedMinipoolMinFee = newCompoundSetting[float64](s.dnts_members, odaoMgr, "members.minipool.unbonded.min.fee")
	s.Members.ChallengeCooldown = newCompoundSetting[time.Duration](s.dnts_members, odaoMgr, "members.challenge.cooldown")
	s.Members.ChallengeWindow = newCompoundSetting[time.Duration](s.dnts_members, odaoMgr, "members.challenge.window")
	s.Members.ChallengeCost = newUintSetting(s.dnts_members, odaoMgr, "members.challenge.cost")

	// Minipool settings
	s.Minipools.ScrubPeriod = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, "minipool.scrub.period")
	s.Minipools.ScrubQuorum = newCompoundSetting[float64](s.dnts_minipool, odaoMgr, "minipool.scrub.quorum")
	s.Minipools.PromotionScrubPeriod = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, "minipool.promotion.scrub.period")
	s.Minipools.IsScrubPenaltyEnabled = newBoolSetting(s.dnts_minipool, odaoMgr, "minipool.scrub.penalty.enabled")
	s.Minipools.BondReductionWindowStart = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, "minipool.bond.reduction.window.start")
	s.Minipools.BondReductionWindowLength = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, "minipool.bond.reduction.window.length")
	s.Minipools.BondReductionCancellationQuorum = newCompoundSetting[float64](s.dnts_minipool, odaoMgr, "minipool.cancel.bond.reduction.quorum")

	// Proposal settings
	s.Proposals.CooldownTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, "proposal.cooldown.time")
	s.Proposals.VoteTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, "proposal.vote.time")
	s.Proposals.VoteDelayTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, "proposal.vote.delay.time")
	s.Proposals.ExecuteTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, "proposal.execute.time")
	s.Proposals.ActionTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, "proposal.action.time")

	return s, nil
}

// =============
// === Calls ===
// =============

func (c *OracleDaoSettings) GetAllDetails(mc *batch.MultiCaller) {
	// Members
	c.Members.Quorum.Get(mc)
	c.Members.RplBond.Get(mc)
	c.Members.UnbondedMinipoolMax.Get(mc)
	c.Members.UnbondedMinipoolMinFee.Get(mc)
	c.Members.ChallengeCooldown.Get(mc)
	c.Members.ChallengeWindow.Get(mc)
	c.Members.ChallengeCost.Get(mc)

	/// Minipools
	c.Minipools.ScrubPeriod.Get(mc)
	c.Minipools.ScrubQuorum.Get(mc)
	c.Minipools.PromotionScrubPeriod.Get(mc)
	c.Minipools.IsScrubPenaltyEnabled.Get(mc)
	c.Minipools.BondReductionWindowStart.Get(mc)
	c.Minipools.BondReductionWindowLength.Get(mc)
	c.Minipools.BondReductionCancellationQuorum.Get(mc)

	// Proposals
	c.Proposals.CooldownTime.Get(mc)
	c.Proposals.VoteTime.Get(mc)
	c.Proposals.VoteDelayTime.Get(mc)
	c.Proposals.ExecuteTime.Get(mc)
	c.Proposals.ActionTime.Get(mc)
}

// === RocketDAONodeTrustedSettingsRewards ===

// Get whether or not the provided rewards network is enabled
func (c *OracleDaoSettings) GetNetworkEnabled(mc *batch.MultiCaller, enabled_Out *bool, network uint64) {
	core.AddCall(mc, c.dnts_rewards, enabled_Out, "getNetworkEnabled", big.NewInt(0).SetUint64(network))
}
