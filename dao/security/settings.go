package security

import (
	"fmt"
	"reflect"

	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Wrapper for a settings category, with all of its settings
type SettingsCategory struct {
	ContractName rocketpool.ContractName
	BoolSettings []ISecurityCouncilSetting[bool]
}

// Binding for security council settings
type SecurityCouncilSettings struct {
	Auction struct {
		IsCreateLotEnabled *SecurityCouncilBoolSetting `json:"isCreateLotEnabled"`
		IsBidOnLotEnabled  *SecurityCouncilBoolSetting `json:"isBidOnLotEnabled"`
	} `json:"auction"`

	Deposit struct {
		IsDepositingEnabled          *SecurityCouncilBoolSetting `json:"isDepositingEnabled"`
		AreDepositAssignmentsEnabled *SecurityCouncilBoolSetting `json:"areDepositAssignmentsEnabled"`
	} `json:"deposit"`

	Minipool struct {
		IsSubmitWithdrawableEnabled *SecurityCouncilBoolSetting `json:"isSubmitWithdrawableEnabled"`
		IsBondReductionEnabled      *SecurityCouncilBoolSetting `json:"isBondReductionEnabled"`
	} `json:"minipool"`

	Network struct {
		IsSubmitBalancesEnabled *SecurityCouncilBoolSetting `json:"isSubmitBalancesEnabled"`
		IsSubmitRewardsEnabled  *SecurityCouncilBoolSetting `json:"isSubmitRewardsEnabled"`
	} `json:"network"`

	Node struct {
		IsRegistrationEnabled              *SecurityCouncilBoolSetting `json:"isRegistrationEnabled"`
		IsSmoothingPoolRegistrationEnabled *SecurityCouncilBoolSetting `json:"isSmoothingPoolRegistrationEnabled"`
		IsDepositingEnabled                *SecurityCouncilBoolSetting `json:"isDepositingEnabled"`
		AreVacantMinipoolsEnabled          *SecurityCouncilBoolSetting `json:"areVacantMinipoolsEnabled"`
	} `json:"node"`

	// === Internal fields ===
	rp              *rocketpool.RocketPool
	secMgr          *SecurityCouncilManager
	contractNameMap map[string]rocketpool.ContractName
}

// ====================
// === Constructors ===
// ====================

// Creates a new SecurityCouncilSettings binding
func newSecurityCouncilSettings(secMgr *SecurityCouncilManager, pdaoSettings *protocol.ProtocolDaoSettings) (*SecurityCouncilSettings, error) {
	s := &SecurityCouncilSettings{
		rp:     secMgr.rp,
		secMgr: secMgr,
	}
	s.contractNameMap = map[string]rocketpool.ContractName{
		"Auction":  pdaoSettings.Auction.IsCreateLotEnabled.GetContract(),
		"Deposit":  pdaoSettings.Deposit.IsDepositingEnabled.GetContract(),
		"Minipool": pdaoSettings.Minipool.IsSubmitWithdrawableEnabled.GetContract(),
		"Network":  pdaoSettings.Network.IsSubmitBalancesEnabled.GetContract(),
		"Node":     pdaoSettings.Node.IsRegistrationEnabled.GetContract(),
	}

	// Auction
	s.Auction.IsCreateLotEnabled = newBoolSetting(secMgr, pdaoSettings.Auction.IsCreateLotEnabled)
	s.Auction.IsBidOnLotEnabled = newBoolSetting(secMgr, pdaoSettings.Auction.IsBidOnLotEnabled)

	// Deposit
	s.Deposit.IsDepositingEnabled = newBoolSetting(secMgr, pdaoSettings.Deposit.IsDepositingEnabled)
	s.Deposit.AreDepositAssignmentsEnabled = newBoolSetting(secMgr, pdaoSettings.Deposit.AreDepositAssignmentsEnabled)

	// Minipool
	s.Minipool.IsSubmitWithdrawableEnabled = newBoolSetting(secMgr, pdaoSettings.Minipool.IsSubmitWithdrawableEnabled)
	s.Minipool.IsBondReductionEnabled = newBoolSetting(secMgr, pdaoSettings.Minipool.IsBondReductionEnabled)

	// Network
	s.Network.IsSubmitBalancesEnabled = newBoolSetting(secMgr, pdaoSettings.Network.IsSubmitBalancesEnabled)
	s.Network.IsSubmitRewardsEnabled = newBoolSetting(secMgr, pdaoSettings.Network.IsSubmitRewardsEnabled)

	// Node
	s.Node.IsRegistrationEnabled = newBoolSetting(secMgr, pdaoSettings.Node.IsRegistrationEnabled)
	s.Node.IsSmoothingPoolRegistrationEnabled = newBoolSetting(secMgr, pdaoSettings.Node.IsSmoothingPoolRegistrationEnabled)
	s.Node.IsDepositingEnabled = newBoolSetting(secMgr, pdaoSettings.Node.IsDepositingEnabled)
	s.Node.AreVacantMinipoolsEnabled = newBoolSetting(secMgr, pdaoSettings.Node.AreVacantMinipoolsEnabled)

	return s, nil
}

// =============
// === Calls ===
// =============

// Get all of the settings, organized by the type used in proposals and boostraps
func (c *SecurityCouncilSettings) GetSettings() map[rocketpool.ContractName]SettingsCategory {
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
				panic(fmt.Sprintf("Security Council settings field named %s does not exist in the contract map.", name))
			}
			boolSettings := []ISecurityCouncilSetting[bool]{}

			// Get all of the settings in this cateogry
			categoryFieldVal := settingsVal.Field(i)
			settingCount := categoryFieldType.NumField()
			for j := 0; j < settingCount; j++ {
				setting := categoryFieldVal.Field(i).Interface()

				// Try bool settings
				boolSetting, isBoolSetting := setting.(ISecurityCouncilSetting[bool])
				if isBoolSetting {
					boolSettings = append(boolSettings, boolSetting)
					continue
				}
			}

			settingsCat := SettingsCategory{
				ContractName: name,
				BoolSettings: boolSettings,
			}
			catMap[name] = settingsCat
		}
	}

	return catMap
}
