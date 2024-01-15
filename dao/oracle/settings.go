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

// =====================
// === Setting Names ===
// =====================

type SettingName string

const (
	// Member
	SettingName_Member_Quorum            SettingName = "members.quorum"
	SettingName_Member_RplBond           SettingName = "members.rplbond"
	SettingName_Member_ChallengeCooldown SettingName = "members.challenge.cooldown"
	SettingName_Member_ChallengeWindow   SettingName = "members.challenge.window"
	SettingName_Member_ChallengeCost     SettingName = "members.challenge.cost"

	// Minipool
	SettingName_Minipool_ScrubPeriod                     SettingName = "minipool.scrub.period"
	SettingName_Minipool_ScrubQuorum                     SettingName = "minipool.scrub.quorum"
	SettingName_Minipool_PromotionScrubPeriod            SettingName = "minipool.promotion.scrub.period"
	SettingName_Minipool_IsScrubPenaltyEnabled           SettingName = "minipool.scrub.penalty.enabled"
	SettingName_Minipool_BondReductionWindowStart        SettingName = "minipool.bond.reduction.window.start"
	SettingName_Minipool_BondReductionWindowLength       SettingName = "minipool.bond.reduction.window.length"
	SettingName_Minipool_BondReductionCancellationQuorum SettingName = "minipool.cancel.bond.reduction.quorum"

	// Proposal
	SettingName_Proposal_CooldownTime  SettingName = "proposal.cooldown.time"
	SettingName_Proposal_VoteTime      SettingName = "proposal.vote.time"
	SettingName_Proposal_VoteDelayTime SettingName = "proposal.vote.delay.time"
	SettingName_Proposal_ExecuteTime   SettingName = "proposal.execute.time"
	SettingName_Proposal_ActionTime    SettingName = "proposal.action.time"
)

// ===============
// === Structs ===
// ===============

// Wrapper for a settings category, with all of its settings
type SettingsCategory struct {
	ContractName rocketpool.ContractName
	BoolSettings []IOracleDaoSetting[bool]
	UintSettings []IOracleDaoSetting[*big.Int]
}

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
	rp              *rocketpool.RocketPool
	odaoMgr         *OracleDaoManager
	dnts_members    *core.Contract
	dnts_minipool   *core.Contract
	dnts_proposals  *core.Contract
	dnts_rewards    *core.Contract
	contractNameMap map[string]rocketpool.ContractName
}

// ====================
// === Constructors ===
// ====================

// Creates a new Oracle DAO settings binding
func newOracleDaoSettings(odaoMgr *OracleDaoManager) (*OracleDaoSettings, error) {
	// Get the contracts
	contractNames := []rocketpool.ContractName{
		rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers,
		rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool,
		rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals,
		rocketpool.ContractName_RocketDAONodeTrustedSettingsRewards,
	}
	contracts, err := odaoMgr.rp.GetContracts(contractNames...)
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
	s.contractNameMap = map[string]rocketpool.ContractName{
		"Member":   contractNames[0],
		"Minipool": contractNames[1],
		"Proposal": contractNames[2],
	}

	// Member
	s.Member.Quorum = newCompoundSetting[float64](s.dnts_members, odaoMgr, SettingName_Member_Quorum)
	s.Member.RplBond = newUintSetting(s.dnts_members, odaoMgr, SettingName_Member_RplBond)
	s.Member.ChallengeCooldown = newCompoundSetting[time.Duration](s.dnts_members, odaoMgr, SettingName_Member_ChallengeCooldown)
	s.Member.ChallengeWindow = newCompoundSetting[time.Duration](s.dnts_members, odaoMgr, SettingName_Member_ChallengeWindow)
	s.Member.ChallengeCost = newUintSetting(s.dnts_members, odaoMgr, SettingName_Member_ChallengeCost)

	// Minipool
	s.Minipool.ScrubPeriod = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, SettingName_Minipool_ScrubPeriod)
	s.Minipool.ScrubQuorum = newCompoundSetting[float64](s.dnts_minipool, odaoMgr, SettingName_Minipool_ScrubQuorum)
	s.Minipool.PromotionScrubPeriod = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, SettingName_Minipool_PromotionScrubPeriod)
	s.Minipool.IsScrubPenaltyEnabled = newBoolSetting(s.dnts_minipool, odaoMgr, SettingName_Minipool_IsScrubPenaltyEnabled)
	s.Minipool.BondReductionWindowStart = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, SettingName_Minipool_BondReductionWindowStart)
	s.Minipool.BondReductionWindowLength = newCompoundSetting[time.Duration](s.dnts_minipool, odaoMgr, SettingName_Minipool_BondReductionWindowLength)
	s.Minipool.BondReductionCancellationQuorum = newCompoundSetting[float64](s.dnts_minipool, odaoMgr, SettingName_Minipool_BondReductionCancellationQuorum)

	// Proposal
	s.Proposal.CooldownTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, SettingName_Proposal_CooldownTime)
	s.Proposal.VoteTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, SettingName_Proposal_VoteTime)
	s.Proposal.VoteDelayTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, SettingName_Proposal_VoteDelayTime)
	s.Proposal.ExecuteTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, SettingName_Proposal_ExecuteTime)
	s.Proposal.ActionTime = newCompoundSetting[time.Duration](s.dnts_proposals, odaoMgr, SettingName_Proposal_ActionTime)

	return s, nil
}

// =============
// === Calls ===
// =============

// Get all of the settings, organized by the type used in proposals and boostraps
func (c *OracleDaoSettings) GetSettings() map[rocketpool.ContractName]SettingsCategory {
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
				panic(fmt.Sprintf("Oracle DAO settings field named %s does not exist in the contract map.", name))
			}
			boolSettings := []IOracleDaoSetting[bool]{}
			uintSettings := []IOracleDaoSetting[*big.Int]{}

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

// Get whether or not the provided rewards network is enabled
func (c *OracleDaoSettings) GetNetworkEnabled(mc *batch.MultiCaller, enabled_Out *bool, network uint64) {
	core.AddCall(mc, c.dnts_rewards, enabled_Out, "getNetworkEnabled", big.NewInt(0).SetUint64(network))
}
