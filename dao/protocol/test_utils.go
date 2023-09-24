//go:build testing
// +build testing

package protocol

import (
	"fmt"
	"reflect"

	"github.com/rocket-pool/rocketpool-go/core"
)

func (s *ProtocolDaoBoolSetting) Equals(other core.IEquatable) (bool, string, string) {
	castedOther, ok := other.(*ProtocolDaoBoolSetting)
	if !ok {
		return false, fmt.Sprintf("[type %s]", reflect.TypeOf(s)), fmt.Sprintf("[type %s]", reflect.TypeOf(other))
	}
	return s.SimpleField.Equals(castedOther.SimpleField)
}

func (s *ProtocolDaoUintSetting) Equals(other core.IEquatable) (bool, string, string) {
	castedOther, ok := other.(*ProtocolDaoUintSetting)
	if !ok {
		return false, fmt.Sprintf("[type %s]", reflect.TypeOf(s)), fmt.Sprintf("[type %s]", reflect.TypeOf(other))
	}
	return s.SimpleField.Equals(castedOther.SimpleField)
}

func (s *ProtocolDaoCompoundSetting[DataType]) Equals(other core.IEquatable) (bool, string, string) {
	castedOther, ok := other.(*ProtocolDaoCompoundSetting[DataType])
	if !ok {
		return false, fmt.Sprintf("[type %s]", reflect.TypeOf(s)), fmt.Sprintf("[type %s]", reflect.TypeOf(other))
	}
	return s.FormattedUint256Field.Equals(castedOther.FormattedUint256Field)
}

func (s *ProtocolDaoBoolSetting) Copy(other core.ICopyable) {
	castedOther, ok := other.(*ProtocolDaoBoolSetting)
	if !ok {
		panic(fmt.Sprintf("wrong type; expected %s, got %s", reflect.TypeOf(s), reflect.TypeOf(other)))
	}
	castedOther.SimpleField.Copy(s.SimpleField)
}

func (s *ProtocolDaoUintSetting) Copy(other core.ICopyable) {
	castedOther, ok := other.(*ProtocolDaoUintSetting)
	if !ok {
		panic(fmt.Sprintf("wrong type; expected %s, got %s", reflect.TypeOf(s), reflect.TypeOf(other)))
	}
	castedOther.SimpleField.Copy(s.SimpleField)
}

func (s *ProtocolDaoCompoundSetting[DataType]) Copy(other core.ICopyable) {
	castedOther, ok := other.(*ProtocolDaoCompoundSetting[DataType])
	if !ok {
		panic(fmt.Sprintf("wrong type; expected %s, got %s", reflect.TypeOf(s), reflect.TypeOf(other)))
	}
	castedOther.FormattedUint256Field.Copy(s.FormattedUint256Field)
}
