package oracle

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

/// ==================
/// === Interfaces ===
/// ==================

// A general interface for settings, parameterized by the type required for proposals and boostrapping
type IOracleDaoSetting[ProposeType core.CallReturnType] interface {
	eth.IQueryable
	GetSettingName() SettingName
	ProposeSet(value ProposeType, opts *bind.TransactOpts) (*eth.TransactionInfo, error)
	Bootstrap(value ProposeType, opts *bind.TransactOpts) (*eth.TransactionInfo, error)
}

/// ===================
/// === BoolSetting ===
/// ===================

// A simple boolean setting
type OracleDaoBoolSetting struct {
	*core.SimpleField[bool]

	// === Internal fields ===
	settingContract rocketpool.ContractName
	odaoMgr         *OracleDaoManager
	settingName     SettingName
}

// Creates a new bool setting
func newBoolSetting(settingContract *core.Contract, odaoMgr *OracleDaoManager, settingName SettingName) *OracleDaoBoolSetting {
	return &OracleDaoBoolSetting{
		SimpleField:     core.NewSimpleField[bool](settingContract, "getSettingBool", settingName),
		settingContract: rocketpool.ContractName(settingContract.Name),
		odaoMgr:         odaoMgr,
		settingName:     settingName,
	}
}

// Gets the underlying path for the setting within the contracts
func (s *OracleDaoBoolSetting) GetSettingName() SettingName {
	return s.settingName
}

// Creates a proposal to change the setting
func (s *OracleDaoBoolSetting) ProposeSet(value bool, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return s.odaoMgr.ProposeSetBool("", s.settingContract, s.settingName, value, opts)
}

// Bootstraps the setting with a new value
func (s *OracleDaoBoolSetting) Bootstrap(value bool, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return s.odaoMgr.BootstrapBool(s.settingContract, s.settingName, value, opts)
}

/// ===================
/// === UintSetting ===
/// ===================

// A simple uint setting
type OracleDaoUintSetting struct {
	*core.SimpleField[*big.Int]

	// === Internal fields ===
	settingContract rocketpool.ContractName
	odaoMgr         *OracleDaoManager
	settingName     SettingName
}

// Creates a new uint setting
func newUintSetting(settingContract *core.Contract, odaoMgr *OracleDaoManager, settingName SettingName) *OracleDaoUintSetting {
	return &OracleDaoUintSetting{
		SimpleField:     core.NewSimpleField[*big.Int](settingContract, "getSettingUint", settingName),
		settingContract: rocketpool.ContractName(settingContract.Name),
		odaoMgr:         odaoMgr,
		settingName:     settingName,
	}
}

// Gets the underlying path for the setting within the contracts
func (s *OracleDaoUintSetting) GetSettingName() SettingName {
	return s.settingName
}

// Creates a proposal to change the setting
func (s *OracleDaoUintSetting) ProposeSet(value *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return s.odaoMgr.ProposeSetUint("", s.settingContract, s.settingName, value, opts)
}

// Bootstraps the setting with a new value
func (s *OracleDaoUintSetting) Bootstrap(value *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return s.odaoMgr.BootstrapUint(s.settingContract, s.settingName, value, opts)
}

/// =======================
/// === CompoundSetting ===
/// =======================

// A uint256 setting that can be formatted to a more well-defined type
type OracleDaoCompoundSetting[DataType core.FormattedUint256Type] struct {
	*core.FormattedUint256Field[DataType]

	// === Internal fields ===
	settingContract rocketpool.ContractName
	odaoMgr         *OracleDaoManager
	settingName     SettingName
}

// Creates a new compound setting
func newCompoundSetting[DataType core.FormattedUint256Type](settingContract *core.Contract, odaoMgr *OracleDaoManager, settingName SettingName) *OracleDaoCompoundSetting[DataType] {
	s := &OracleDaoCompoundSetting[DataType]{
		FormattedUint256Field: core.NewFormattedUint256Field[DataType](settingContract, "getSettingUint", settingName),
		settingContract:       rocketpool.ContractName(settingContract.Name),
		odaoMgr:               odaoMgr,
		settingName:           settingName,
	}

	return s
}

// Gets the underlying path for the setting within the contracts
func (s *OracleDaoCompoundSetting[DataType]) GetSettingName() SettingName {
	return s.settingName
}

// Creates a proposal to change the setting
func (s *OracleDaoCompoundSetting[DataType]) ProposeSet(value *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return s.odaoMgr.ProposeSetUint("", s.settingContract, s.settingName, value, opts)
}

// Bootstraps the setting with a new value
func (s *OracleDaoCompoundSetting[DataType]) Bootstrap(value *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return s.odaoMgr.BootstrapUint(s.settingContract, s.settingName, value, opts)
}
