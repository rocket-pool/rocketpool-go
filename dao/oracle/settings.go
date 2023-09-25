package oracle

import (
	"fmt"
	"math/big"
	"reflect"
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
	// Member
	Member struct {
		Quorum            *OracleDaoCompoundSetting[float64]
		RplBond           *OracleDaoUintSetting
		ChallengeCooldown *OracleDaoCompoundSetting[time.Duration]
		ChallengeWindow   *OracleDaoCompoundSetting[time.Duration]
		ChallengeCost     *OracleDaoUintSetting
	}

	// Minipool
	Minipool struct {
		ScrubPeriod                     *OracleDaoCompoundSetting[time.Duration]
		ScrubQuorum                     *OracleDaoCompoundSetting[float64]
		PromotionScrubPeriod            *OracleDaoCompoundSetting[time.Duration]
		IsScrubPenaltyEnabled           *OracleDaoBoolSetting
		BondReductionWindowStart        *OracleDaoCompoundSetting[time.Duration]
		BondReductionWindowLength       *OracleDaoCompoundSetting[time.Duration]
		BondReductionCancellationQuorum *OracleDaoCompoundSetting[float64]
	}

	// Proposal
	Proposal struct {
		CooldownTime  *OracleDaoCompoundSetting[time.Duration]
		VoteTime      *OracleDaoCompoundSetting[time.Duration]
		VoteDelayTime *OracleDaoCompoundSetting[time.Duration]
		ExecuteTime   *OracleDaoCompoundSetting[time.Duration]
		ActionTime    *OracleDaoCompoundSetting[time.Duration]
	}

	// === Internal fields ===
	rp             *rocketpool.RocketPool
	odaoMgr        *OracleDaoManager
	dnts_members   *core.Contract
	dnts_minipool  *core.Contract
	dnts_proposals *core.Contract
	dnts_rewards   *core.Contract
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
		rp:      odaoMgr.rp,
		odaoMgr: odaoMgr,

		dnts_members:   contracts[0],
		dnts_minipool:  contracts[1],
		dnts_proposals: contracts[2],
		dnts_rewards:   contracts[3],
	}

	// Member
	s.Member.Quorum = newCompoundSetting[float64](s.dnts_members, odaoMgr, "members.quorum")
	s.Member.RplBond = newUintSetting(s.dnts_members, odaoMgr, "members.rplbond")
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

// Get all of the settings, organized by the type used in proposals and boostraps
func (c *OracleDaoSettings) GetSettings() ([]IOracleDaoSetting[bool], []IOracleDaoSetting[*big.Int]) {
	boolSettings := []IOracleDaoSetting[bool]{}
	uintSettings := []IOracleDaoSetting[*big.Int]{}

	settingsType := reflect.TypeOf(c)
	settingsVal := reflect.ValueOf(c)
	fieldCount := settingsType.NumField()
	for i := 0; i < fieldCount; i++ {
		categoryFieldType := settingsType.Field(i).Type

		// A container struct for settings by category
		if categoryFieldType.Kind() == reflect.Struct {
			// Get all of the settings in this cateogry
			categoryFieldVal := settingsVal.Field(i)
			settingCount := categoryFieldType.NumField()
			for j := 0; j < settingCount; j++ {
				setting := categoryFieldVal.Field(i).Interface()

				// Try bool settings
				boolSetting, isBoolSetting := setting.(IOracleDaoSetting[bool])
				if isBoolSetting {
					boolSettings = append(boolSettings, boolSetting)
					continue
				}

				// Try uint settings
				uintSetting, isUintSetting := setting.(IOracleDaoSetting[*big.Int])
				if isUintSetting {
					uintSettings = append(uintSettings, uintSetting)
				}
			}

		}
	}

	return boolSettings, uintSettings
}

// Get whether or not the provided rewards network is enabled
func (c *OracleDaoSettings) GetNetworkEnabled(mc *batch.MultiCaller, enabled_Out *bool, network uint64) {
	core.AddCall(mc, c.dnts_rewards, enabled_Out, "getNetworkEnabled", big.NewInt(0).SetUint64(network))
}
