//go:build testing
// +build testing

package core

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/beacon"
)

// ==================
// === Interfaces ===
// ==================

type IEquatable interface {
	Equals(other IEquatable) (bool, string, string)
}

type ICopyable interface {
	Copy(other ICopyable)
}

// ===================
// === SimpleField ===
// ===================

func (f *SimpleField[ValueType]) Set(value ValueType) {
	switch val := any(&f.value).(type) {
	case **big.Int:
		castedValue := any(value).(*big.Int)
		*val = big.NewInt(0).Set(castedValue)
	case *common.Address:
		castedValue := any(value).(common.Address)
		copy((*val).Bytes(), castedValue.Bytes())
	case *common.Hash:
		castedValue := any(value).(common.Hash)
		copy((*val).Bytes(), castedValue.Bytes())
	case *beacon.ValidatorPubkey:
		castedValue := any(value).(beacon.ValidatorPubkey)
		copy((*val)[:], castedValue[:])
	case *[]byte:
		castedValue := any(value).([]byte)
		copy(*val, castedValue)
	default:
		f.value = value
	}
}

func (f *SimpleField[ValueType]) Equals(other IEquatable) (bool, string, string) {
	castedOther, ok := other.(*SimpleField[ValueType])
	if !ok {
		return false, fmt.Sprintf("[type %s]", reflect.TypeOf(f)), fmt.Sprintf("[type %s]", reflect.TypeOf(other))
	}

	switch val := any(f.value).(type) {
	case *big.Int:
		otherVal := any(castedOther.value).(*big.Int)
		return val.Cmp(otherVal) == 0, val.String(), otherVal.String()
	case []byte:
		otherVal := any(castedOther.value).([]byte)
		return bytes.Equal(val, otherVal), fmt.Sprint(val), fmt.Sprint(otherVal)
	default:
		return val == any(castedOther.value), fmt.Sprint(val), fmt.Sprint(castedOther.value)
	}
}

func (f *SimpleField[ValueType]) Copy(other ICopyable) {
	castedOther, ok := other.(*SimpleField[ValueType])
	if !ok {
		return
	}
	f.Set(castedOther.value)
}

// =============================
// === FormattedUint256Field ===
// =============================

func (f *FormattedUint256Field[ValueType]) Set(value ValueType) {
	f.value = GetValueForUint256(value)
}

func (f *FormattedUint256Field[ValueType]) SetRawValue(value *big.Int) {
	f.value = big.NewInt(0).Set(value)
}

func (f *FormattedUint256Field[ValueType]) Equals(other IEquatable) (bool, string, string) {
	castedOther, ok := other.(*FormattedUint256Field[ValueType])
	if !ok {
		return false, fmt.Sprintf("[type %s]", reflect.TypeOf(f)), fmt.Sprintf("[type %s]", reflect.TypeOf(other))
	}
	return f.value.Cmp(castedOther.value) == 0, f.value.String(), castedOther.value.String()
}

func (f *FormattedUint256Field[ValueType]) Copy(other ICopyable) {
	castedOther, ok := other.(*FormattedUint256Field[ValueType])
	if !ok {
		return
	}
	f.SetRawValue(castedOther.value)
}

// ===========================
// === FormattedUint8Field ===
// ===========================

// Gets the raw value after it's been queried
func (f *FormattedUint8Field[ValueType]) Set(value ValueType) {
	f.value = uint8(value)
}

func (f *FormattedUint8Field[ValueType]) Equals(other IEquatable) (bool, string, string) {
	castedOther, ok := other.(*FormattedUint8Field[ValueType])
	if !ok {
		return false, fmt.Sprintf("[type %s]", reflect.TypeOf(f)), fmt.Sprintf("[type %s]", reflect.TypeOf(other))
	}
	return f.value == castedOther.value, fmt.Sprint(f.value), fmt.Sprint(castedOther.value)
}

func (f *FormattedUint8Field[ValueType]) Copy(other ICopyable) {
	castedOther, ok := other.(*FormattedUint8Field[ValueType])
	if !ok {
		return
	}
	f.value = castedOther.value
}
