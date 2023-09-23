package protocol

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

/// ===================
/// === BoolSetting ===
/// ===================

type ProtocolDaoBoolSetting struct {
	Value bool

	settingContract *core.Contract
	pdaoMgr         *ProtocolDaoManager
	path            string
}

func newBoolSetting(settingContract *core.Contract, pdaoMgr *ProtocolDaoManager, path string) *ProtocolDaoBoolSetting {
	return &ProtocolDaoBoolSetting{
		settingContract: settingContract,
		pdaoMgr:         pdaoMgr,
		path:            path,
	}
}

func (s *ProtocolDaoBoolSetting) Get(mc *batch.MultiCaller) {
	core.AddCall(mc, s.settingContract, &s.Value, "getSettingBool", s.path)
}

// Uncomment for Houston
/*
func (s *ProtocolDaoBoolSetting) ProposeSet(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.pdaoMgr.ProposeSetBool("", rocketpool.ContractName(s.settingContract.Name), s.path, value, opts)
}
*/

func (s *ProtocolDaoBoolSetting) Bootstrap(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.pdaoMgr.BootstrapBool(rocketpool.ContractName(s.settingContract.Name), s.path, value, opts)
}

func (s *ProtocolDaoBoolSetting) GetRawValue() bool {
	return s.Value
}

func (s *ProtocolDaoBoolSetting) SetRawValue(value bool) {
	s.Value = value
}

/// ===================
/// === UintSetting ===
/// ===================

type ProtocolDaoUintSetting struct {
	Value *big.Int

	settingContract *core.Contract
	pdaoMgr         *ProtocolDaoManager
	path            string
}

func newUintSetting(settingContract *core.Contract, pdaoMgr *ProtocolDaoManager, path string) *ProtocolDaoUintSetting {
	return &ProtocolDaoUintSetting{
		settingContract: settingContract,
		pdaoMgr:         pdaoMgr,
		path:            path,
	}
}

func (s *ProtocolDaoUintSetting) Get(mc *batch.MultiCaller) {
	core.AddCall(mc, s.settingContract, &s.Value, "getSettingUint", s.path)
}

// Uncomment for Houston
/*
func (s *ProtocolDaoUintSetting) ProposeSet(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.pdaoMgr.ProposeSetUint("", rocketpool.ContractName(s.settingContract.Name), s.path, value, opts)
}
*/

func (s *ProtocolDaoUintSetting) Bootstrap(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.pdaoMgr.BootstrapUint(rocketpool.ContractName(s.settingContract.Name), s.path, value, opts)
}

func (s *ProtocolDaoUintSetting) GetRawValue() *big.Int {
	return s.Value
}

func (s *ProtocolDaoUintSetting) SetRawValue(value *big.Int) {
	s.Value = big.NewInt(0).Set(value)
}

/// =======================
/// === CompoundSetting ===
/// =======================

type ProtocolDaoCompoundSetting[DataType core.FormattedType] struct {
	Value core.Uint256Parameter[DataType]

	settingContract *core.Contract
	pdaoMgr         *ProtocolDaoManager
	path            string
}

func newCompoundSetting[DataType core.FormattedType](settingContract *core.Contract, pdaoMgr *ProtocolDaoManager, path string) *ProtocolDaoCompoundSetting[DataType] {
	s := &ProtocolDaoCompoundSetting[DataType]{
		settingContract: settingContract,
		pdaoMgr:         pdaoMgr,
		path:            path,
	}

	return s
}

func (s *ProtocolDaoCompoundSetting[DataType]) Get(mc *batch.MultiCaller) {
	core.AddCall(mc, s.settingContract, &s.Value.RawValue, "getSettingUint", s.path)
}

// Uncomment for Houston
/*
func (s *ProtocolDaoCompoundSetting[DataType]) ProposeSet(value core.Parameter[DataType], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.pdaoMgr.ProposeSetUint("", rocketpool.ContractName(s.settingContract.Name), s.path, s.Value.RawValue, opts)
}
*/

func (s *ProtocolDaoCompoundSetting[DataType]) Bootstrap(value core.Uint256Parameter[DataType], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.pdaoMgr.BootstrapUint(rocketpool.ContractName(s.settingContract.Name), s.path, s.Value.RawValue, opts)
}

func (s *ProtocolDaoCompoundSetting[DataType]) GetRawValue() *big.Int {
	return s.Value.GetRawValue()
}

func (s *ProtocolDaoCompoundSetting[DataType]) SetRawValue(value *big.Int) {
	s.Value.SetRawValue(value)
}
