package security

import (
	"fmt"
	"reflect"

	"github.com/rocket-pool/rocketpool-go/v2/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
)

// =================
// === Constants ===
// =================

const (
	AuctionNamespace  string = "auction"
	DepositNamespace  string = "deposit"
	MinipoolNamespace string = "minipool"
	NetworkNamespace  string = "network"
	NodeNamespace     string = "node"
)

// ===============
// === Structs ===
// ===============

// Wrapper for a settings category, with all of its settings
type SettingsCategory struct {
	ContractName rocketpool.ContractName
	BoolSettings []*SecurityCouncilBoolSetting
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
	s.Auction.IsCreateLotEnabled = newBoolSetting(secMgr, pdaoSettings.Auction.IsCreateLotEnabled, AuctionNamespace)
	s.Auction.IsBidOnLotEnabled = newBoolSetting(secMgr, pdaoSettings.Auction.IsBidOnLotEnabled, AuctionNamespace)

	// Deposit
	s.Deposit.IsDepositingEnabled = newBoolSetting(secMgr, pdaoSettings.Deposit.IsDepositingEnabled, DepositNamespace)
	s.Deposit.AreDepositAssignmentsEnabled = newBoolSetting(secMgr, pdaoSettings.Deposit.AreDepositAssignmentsEnabled, DepositNamespace)

	// Minipool
	s.Minipool.IsSubmitWithdrawableEnabled = newBoolSetting(secMgr, pdaoSettings.Minipool.IsSubmitWithdrawableEnabled, MinipoolNamespace)
	s.Minipool.IsBondReductionEnabled = newBoolSetting(secMgr, pdaoSettings.Minipool.IsBondReductionEnabled, MinipoolNamespace)

	// Network
	s.Network.IsSubmitBalancesEnabled = newBoolSetting(secMgr, pdaoSettings.Network.IsSubmitBalancesEnabled, NetworkNamespace)
	s.Network.IsSubmitRewardsEnabled = newBoolSetting(secMgr, pdaoSettings.Network.IsSubmitRewardsEnabled, NetworkNamespace)

	// Node
	s.Node.IsRegistrationEnabled = newBoolSetting(secMgr, pdaoSettings.Node.IsRegistrationEnabled, NodeNamespace)
	s.Node.IsSmoothingPoolRegistrationEnabled = newBoolSetting(secMgr, pdaoSettings.Node.IsSmoothingPoolRegistrationEnabled, NodeNamespace)
	s.Node.IsDepositingEnabled = newBoolSetting(secMgr, pdaoSettings.Node.IsDepositingEnabled, NodeNamespace)
	s.Node.AreVacantMinipoolsEnabled = newBoolSetting(secMgr, pdaoSettings.Node.AreVacantMinipoolsEnabled, NodeNamespace)

	return s, nil
}

// =============
// === Calls ===
// =============

// Get all of the settings, organized by the type used in proposals and boostraps
func (c *SecurityCouncilSettings) GetSettings() map[rocketpool.ContractName]SettingsCategory {
	catMap := map[rocketpool.ContractName]SettingsCategory{}

	settingsVal := reflect.ValueOf(c).Elem()
	settingsType := reflect.TypeOf(settingsVal.Interface())
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
			boolSettings := []*SecurityCouncilBoolSetting{}

			// Get all of the settings in this cateogry
			categoryFieldVal := settingsVal.Field(i)
			settingCount := categoryFieldType.NumField()
			for j := 0; j < settingCount; j++ {
				setting := categoryFieldVal.Field(j).Interface()

				// Try bool settings
				boolSetting, isBoolSetting := setting.(*SecurityCouncilBoolSetting)
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
