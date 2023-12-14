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
