package security

import (
	"reflect"

	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

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
	rp     *rocketpool.RocketPool
	secMgr *SecurityCouncilManager
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
func (c *SecurityCouncilSettings) GetSettings() []ISecurityCouncilSetting[bool] {
	boolSettings := []ISecurityCouncilSetting[bool]{}

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
				boolSetting, isBoolSetting := setting.(ISecurityCouncilSetting[bool])
				if isBoolSetting {
					boolSettings = append(boolSettings, boolSetting)
					continue
				}
			}

		}
	}

	return boolSettings
}
