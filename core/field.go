package core

import (
	"math/big"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// ==================
// === Interfaces ===
// ==================

// Represents structs that can have their values queried during a multicall
type IQueryable interface {
	// Adds the struct's values to the provided multicall query before it runs
	AddToQuery(mc *batch.MultiCaller)
}

// ===================
// === SimpleField ===
// ===================

type SimpleField[DataType CallReturnType] struct {
	Contract   *Contract
	GetterName string

	value DataType
}

func (f *SimpleField[DataType]) AddToQuery(mc *batch.MultiCaller) {
	AddCall(mc, f.Contract, &f.value, f.GetterName)
}

func (f *SimpleField[DataType]) Get() DataType {
	return f.value
}

// =============================
// === FormattedUint256Field ===
// =============================

// A collection of legal types for FormattedUint256Field
type FormattedUint256Type interface {
	time.Time | uint64 | int64 | float64 | time.Duration
}

// A field that is stored as a uint256 in the contracts, but represents a more well-defined type
type FormattedUint256Field[ValueType FormattedUint256Type] struct {
	contract   *Contract
	getterName string
	value      *big.Int
}

// Creates a new FormattedUint256Type instance
func NewFormattedUint256Field[ValueType FormattedUint256Type](contract *Contract, getterName string) FormattedUint256Field[ValueType] {
	return FormattedUint256Field[ValueType]{
		contract:   contract,
		getterName: getterName,
	}
}

// Adds a query to the field's value to the multicall
func (f *FormattedUint256Field[DataType]) AddToQuery(mc *batch.MultiCaller) {
	AddCall(mc, f.contract, &f.value, f.getterName)
}

// Gets the raw value after it's been queried
func (f *FormattedUint256Field[DataType]) Raw() *big.Int {
	return f.value
}

// Gets the value after it's been queried, converted to the more well-defined type
func (f *FormattedUint256Field[DataType]) Formatted() DataType {
	// Switch on the parameter type and convert it
	var out DataType
	switch outPtr := any(&out).(type) { // Go can't switch on type parameters yet so we have to do this nonsense
	case *time.Time:
		*outPtr = time.Unix(f.value.Int64(), 0)
	case *uint64:
		*outPtr = f.value.Uint64()
	case *int64:
		*outPtr = f.value.Int64()
	case *float64:
		*outPtr = eth.WeiToEth(f.value)
	case *time.Duration:
		*outPtr = time.Duration(f.value.Int64()) * time.Second
	}
	return out
}

// =============================
// === FormattedUint8Field ===
// =============================

type FormattedUint8Field[DataType FormattedUint8Type] struct {
	Contract   *Contract
	GetterName string

	value uint8
}

func (f *FormattedUint8Field[DataType]) AddToQuery(mc *batch.MultiCaller) {
	AddCall(mc, f.Contract, &f.value, f.GetterName)
}

func (f *FormattedUint8Field[DataType]) Raw() uint8 {
	return f.value
}

func (f *FormattedUint8Field[DataType]) Formatted() DataType {
	// Switch on the parameter type and convert it
	var out DataType
	switch outPtr := any(&out).(type) { // Go can't switch on type parameters yet so we have to do this nonsense
	case *types.MinipoolStatus:
		*outPtr = types.MinipoolStatus(f.value)
	case *types.MinipoolDeposit:
		*outPtr = types.MinipoolDeposit(f.value)
	}
	return out
}
