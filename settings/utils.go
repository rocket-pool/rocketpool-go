package settings

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// The valid type definitions for a setting
type settingType interface {
	uint64 | bool | float64 | *big.Int | time.Duration
}

// Interface for a contract binding that can be used to bootstrap a parameter
type IBootstrapBinding interface {
	BootstrapBool(rocketpool.ContractName, string, bool, *bind.TransactOpts) (*core.TransactionInfo, error)
	BootstrapUint(rocketpool.ContractName, string, *big.Int, *bind.TransactOpts) (*core.TransactionInfo, error)
}

// Interface for a contract binding that can be used to propose a new parameter setting
type IProposeBinding interface {
	ProposeSetBool(string, rocketpool.ContractName, string, bool, *bind.TransactOpts) (*core.TransactionInfo, error)
	ProposeSetUint(string, rocketpool.ContractName, string, *big.Int, *bind.TransactOpts) (*core.TransactionInfo, error)
}

// Gets info for bootstrapping a new value directly
func bootstrapValue[valueType settingType](settings IBootstrapBinding, contractName rocketpool.ContractName, path string, value valueType, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	// Switch on the parameter type and convert it
	switch value := any(&value).(type) { // Go can't switch on type parameters yet so we have to do this nonsense
	case uint64:
		trueVal := big.NewInt(0).SetUint64(value)
		return settings.BootstrapUint(contractName, path, trueVal, opts)
	case bool:
		return settings.BootstrapBool(contractName, path, value, opts)
	case float64:
		trueVal := eth.EthToWei(value)
		return settings.BootstrapUint(contractName, path, trueVal, opts)
	case *big.Int:
		return settings.BootstrapUint(contractName, path, value, opts)
	case time.Duration:
		trueVal := big.NewInt(0).SetUint64(uint64(value.Seconds()))
		return settings.BootstrapUint(contractName, path, trueVal, opts)
	}

	return nil, fmt.Errorf("unexpected value or type: %v", value)
}

// Gets info for submitting a proposal to change a value
func proposeSetValue[valueType settingType](settings IProposeBinding, contractName rocketpool.ContractName, path string, value valueType, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	message := fmt.Sprintf("set %s", path)

	// Switch on the parameter type and convert it
	switch value := any(&value).(type) { // Go can't switch on type parameters yet so we have to do this nonsense
	case uint64:
		trueVal := big.NewInt(0).SetUint64(value)
		return settings.ProposeSetUint(message, contractName, path, trueVal, opts)
	case bool:
		return settings.ProposeSetBool(message, contractName, path, value, opts)
	case float64:
		trueVal := eth.EthToWei(value)
		return settings.ProposeSetUint(message, contractName, path, trueVal, opts)
	case *big.Int:
		return settings.ProposeSetUint(message, contractName, path, value, opts)
	case time.Duration:
		trueVal := big.NewInt(0).SetUint64(uint64(value.Seconds()))
		return settings.ProposeSetUint(message, contractName, path, trueVal, opts)
	}

	return nil, fmt.Errorf("unexpected value or type: %v", value)
}
