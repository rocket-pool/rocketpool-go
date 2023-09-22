package oracle

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

type OracleDaoBoolSetting struct {
	Value bool

	settingContract *core.Contract
	odaoMgr         *OracleDaoManager
	path            string
}

func newBoolSetting(settingContract *core.Contract, odaoMgr *OracleDaoManager, path string) *OracleDaoBoolSetting {
	return &OracleDaoBoolSetting{
		settingContract: settingContract,
		odaoMgr:         odaoMgr,
		path:            path,
	}
}

func (s *OracleDaoBoolSetting) Get(mc *batch.MultiCaller) {
	core.AddCall(mc, s.settingContract, &s.Value, "getSettingBool", s.path)
}

func (s *OracleDaoBoolSetting) ProposeSet(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.odaoMgr.ProposeSetBool("", rocketpool.ContractName(s.settingContract.Name), s.path, value, opts)
}

func (s *OracleDaoBoolSetting) Bootstrap(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.odaoMgr.BootstrapBool(rocketpool.ContractName(s.settingContract.Name), s.path, value, opts)
}

/// ===================
/// === UintSetting ===
/// ===================

type OracleDaoUintSetting struct {
	Value *big.Int

	settingContract *core.Contract
	odaoMgr         *OracleDaoManager
	path            string
}

func newUintSetting(settingContract *core.Contract, odaoMgr *OracleDaoManager, path string) *OracleDaoUintSetting {
	return &OracleDaoUintSetting{
		settingContract: settingContract,
		odaoMgr:         odaoMgr,
		path:            path,
	}
}

func (s *OracleDaoUintSetting) Get(mc *batch.MultiCaller) {
	core.AddCall(mc, s.settingContract, &s.Value, "getSettingUint", s.path)
}

func (s *OracleDaoUintSetting) ProposeSet(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.odaoMgr.ProposeSetUint("", rocketpool.ContractName(s.settingContract.Name), s.path, value, opts)
}

func (s *OracleDaoUintSetting) Bootstrap(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.odaoMgr.BootstrapUint(rocketpool.ContractName(s.settingContract.Name), s.path, value, opts)
}

/// =======================
/// === CompoundSetting ===
/// =======================

type OracleDaoCompoundSetting[DataType core.FormattedType] struct {
	Value core.Parameter[DataType]

	settingContract *core.Contract
	odaoMgr         *OracleDaoManager
	path            string
}

func newCompoundSetting[DataType core.FormattedType](settingContract *core.Contract, odaoMgr *OracleDaoManager, path string) *OracleDaoCompoundSetting[DataType] {
	s := &OracleDaoCompoundSetting[DataType]{
		settingContract: settingContract,
		odaoMgr:         odaoMgr,
		path:            path,
	}

	return s
}

func (s *OracleDaoCompoundSetting[DataType]) Get(mc *batch.MultiCaller) {
	core.AddCall(mc, s.settingContract, &s.Value.RawValue, "getSettingUint", s.path)
}

func (s *OracleDaoCompoundSetting[DataType]) ProposeSet(value core.Parameter[DataType], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.odaoMgr.ProposeSetUint("", rocketpool.ContractName(s.settingContract.Name), s.path, s.Value.RawValue, opts)
}

func (s *OracleDaoCompoundSetting[DataType]) Bootstrap(value core.Parameter[DataType], opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return s.odaoMgr.BootstrapUint(rocketpool.ContractName(s.settingContract.Name), s.path, s.Value.RawValue, opts)
}
