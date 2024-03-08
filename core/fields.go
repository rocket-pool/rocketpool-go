package core

import (
	"math/big"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/types"
)

// ===================
// === SimpleField ===
// ===================

// A field where the underlying value type is the same as the type stored in the contracts
type SimpleField[ValueType CallReturnType] struct {
	contract   *Contract
	getterName string
	args       []any
	value      ValueType
}

// Creates a new SimpleField instance
func NewSimpleField[ValueType CallReturnType](contract *Contract, getterName string, args ...any) *SimpleField[ValueType] {
	return &SimpleField[ValueType]{
		contract:   contract,
		getterName: getterName,
		args:       args,
	}
}

// Adds a query to the field's value to the multicall
func (f *SimpleField[ValueType]) AddToQuery(mc *batch.MultiCaller) {
	AddCall(mc, f.contract, &f.value, f.getterName, f.args...)
}

// Gets the field's value after it's been queried
func (f *SimpleField[ValueType]) Get() ValueType {
	return f.value
}

// =============================
// === FormattedUint256Field ===
// =============================

// A collection of legal types for FormattedUint256Field
type FormattedUint256Type interface {
	time.Time | uint32 | uint64 | int64 | float64 | time.Duration
}

// A field that is stored as a uint256 in the contracts, but represents a more well-defined type
type FormattedUint256Field[ValueType FormattedUint256Type] struct {
	contract   *Contract
	getterName string
	args       []any
	value      *big.Int
}

// Creates a new FormattedUint256Type instance
func NewFormattedUint256Field[ValueType FormattedUint256Type](contract *Contract, getterName string, args ...any) *FormattedUint256Field[ValueType] {
	return &FormattedUint256Field[ValueType]{
		contract:   contract,
		getterName: getterName,
		args:       args,
	}
}

// Adds a query to the field's value to the multicall
func (f *FormattedUint256Field[ValueType]) AddToQuery(mc *batch.MultiCaller) {
	AddCall(mc, f.contract, &f.value, f.getterName, f.args...)
}

// Gets the raw value after it's been queried
func (f *FormattedUint256Field[ValueType]) Raw() *big.Int {
	return f.value
}

// Gets the value after it's been queried, converted to the more well-defined type
func (f *FormattedUint256Field[ValueType]) Formatted() ValueType {
	// Switch on the parameter type and convert it
	var out ValueType
	switch outPtr := any(&out).(type) {
	case *time.Time:
		*outPtr = time.Unix(f.value.Int64(), 0)
	case *uint32:
		*outPtr = uint32(f.value.Uint64())
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

// ===========================
// === FormattedUint8Field ===
// ===========================

// A collection of legal types for FormattedUint8Field
type FormattedUint8Type interface {
	types.MinipoolStatus | types.MinipoolDeposit | types.ProposalState | types.ProtocolDaoProposalState
}

// A field that is stored as a uint8 in the contracts, but represents a more well-defined type
type FormattedUint8Field[ValueType FormattedUint8Type] struct {
	contract   *Contract
	getterName string
	args       []any
	value      uint8
}

// Creates a new FormattedUint8Type instance
func NewFormattedUint8Field[ValueType FormattedUint8Type](contract *Contract, getterName string, args ...any) *FormattedUint8Field[ValueType] {
	return &FormattedUint8Field[ValueType]{
		contract:   contract,
		getterName: getterName,
		args:       args,
	}
}

// Adds a query to the field's value to the multicall
func (f *FormattedUint8Field[ValueType]) AddToQuery(mc *batch.MultiCaller) {
	AddCall(mc, f.contract, &f.value, f.getterName, f.args...)
}

// Gets the raw value after it's been queried
func (f *FormattedUint8Field[ValueType]) Raw() uint8 {
	return f.value
}

// Gets the value after it's been queried, converted to the more well-defined type
func (f *FormattedUint8Field[ValueType]) Formatted() ValueType {
	// Switch on the parameter type and convert it
	var out ValueType
	switch outPtr := any(&out).(type) {
	case *types.MinipoolStatus:
		*outPtr = types.MinipoolStatus(f.value)
	case *types.MinipoolDeposit:
		*outPtr = types.MinipoolDeposit(f.value)
	case *types.ProposalState:
		*outPtr = types.ProposalState(f.value)
	case *types.ProtocolDaoProposalState:
		*outPtr = types.ProtocolDaoProposalState(f.value)
	}
	return out
}
