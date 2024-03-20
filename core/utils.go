package core

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
)

// This is a helper for adding calls to multicall that has strongly-typed output and can take in RP contracts
func AddCall[OutType CallReturnType](mc *batch.MultiCaller, contract *Contract, output *OutType, method string, args ...any) {
	eth.AddCallToMulticaller(mc, contract.Contract, output, method, args...)
}

// This is a helper for adding calls to multicall that has untyped output and can take in RP contracts
// Only use this in situations where the output type is unique (such as a struct for a specific contract view)
func AddCallRaw(mc *batch.MultiCaller, contract *Contract, output any, method string, args ...any) {
	mc.AddCall(contract.Address, contract.ABI, output, method, args...)
}

// Converts the value to a *big.Int; useful for transactions that require formattable types
func GetValueForUint256[ValueType FormattedUint256Type](value ValueType) *big.Int {
	switch v := any(&value).(type) {
	case *time.Time:
		return big.NewInt(v.Unix())
	case *uint64:
		return big.NewInt(0).SetUint64(*v)
	case *int64:
		return big.NewInt(*v)
	case *float64:
		return eth.EthToWei(*v)
	case *time.Duration:
		return big.NewInt(int64(v.Seconds()))
	default:
		panic(fmt.Sprintf("unexpected type: %s", reflect.TypeOf(value).Name()))
	}
}
