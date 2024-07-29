//go:build testing
// +build testing

package core

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"
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
