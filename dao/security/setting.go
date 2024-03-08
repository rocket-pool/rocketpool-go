package security

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
)

/// ==================
/// === Interfaces ===
/// ==================

// A general interface for settings, parameterized by the type required for proposals
type ISecurityCouncilSetting[ProposeType core.CallReturnType] interface {
	eth.IQueryable
	GetProtocolDaoSetting() protocol.IProtocolDaoSetting[ProposeType]
	ProposeSet(value ProposeType, opts *bind.TransactOpts) (*eth.TransactionInfo, error)
}

/// ===================
/// === BoolSetting ===
/// ===================

// A simple boolean setting
type SecurityCouncilBoolSetting struct {
	// === Internal fields ===
	setting *protocol.ProtocolDaoBoolSetting
	secMgr  *SecurityCouncilManager
}

// Creates a new bool setting
func newBoolSetting(secMgr *SecurityCouncilManager, setting *protocol.ProtocolDaoBoolSetting) *SecurityCouncilBoolSetting {
	return &SecurityCouncilBoolSetting{
		secMgr:  secMgr,
		setting: setting,
	}
}

// Gets the underlying path for the setting within the contracts
func (s *SecurityCouncilBoolSetting) GetProtocolDaoSetting() *protocol.ProtocolDaoBoolSetting {
	return s.setting
}

// Creates a proposal to change the setting
func (s *SecurityCouncilBoolSetting) ProposeSet(value bool, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return s.secMgr.ProposeSetBool("", s.setting.GetContract(), s.setting.GetSettingName(), value, opts)
}
