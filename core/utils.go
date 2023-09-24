package core

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// This is a helper for adding calls to multicall that has strongly-typed output and can take in RP contracts
func AddCall[OutType CallReturnType](mc *batch.MultiCaller, contract *Contract, output *OutType, method string, args ...any) {
	mc.AddCall(*contract.Address, contract.ABI, output, method, args...)
}

// Adds a collection of IQueryable calls to a multicall
func AddQueryablesToMulticall(mc *batch.MultiCaller, queryables ...IQueryable) {
	for _, queryable := range queryables {
		queryable.AddToQuery(mc)
	}
}

// Adds all of the object's fields that implement IQueryable to the provided multicaller
func QueryAllFields(object any, mc *batch.MultiCaller) error {
	objectValue := reflect.ValueOf(object)
	objectType := reflect.TypeOf(object)
	if objectType.Kind() == reflect.Pointer {
		// If this is a pointer, switch to what it's pointing at
		objectValue = objectValue.Elem()
		objectType = objectType.Elem()
	}

	// Run through each field
	for i := 0; i < objectType.NumField(); i++ {
		field := objectValue.Field(i)
		typeField := objectType.Field(i)
		if typeField.IsExported() {
			fieldAsQueryable, isQueryable := field.Interface().(IQueryable)
			if isQueryable {
				// If it's IQueryable, run it
				fieldAsQueryable.AddToQuery(mc)
			} else if typeField.Type.Kind() == reflect.Pointer &&
				typeField.Type.Elem().Kind() == reflect.Struct {
				// If it's a pointer to a struct, recurse
				err := QueryAllFields(field.Interface(), mc)
				if err != nil {
					return err
				}
			} else if typeField.Type.Kind() == reflect.Struct {
				// If it's a struct, recurse
				err := QueryAllFields(field.Interface(), mc)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
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
