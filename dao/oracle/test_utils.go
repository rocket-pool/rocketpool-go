//go:build testing
// +build testing

package oracle

import (
	"fmt"
	"reflect"

	"github.com/rocket-pool/rocketpool-go/core"
)

func (s *OracleDaoBoolSetting) Equals(other core.IEquatable) (bool, string, string) {
	castedOther, ok := other.(*OracleDaoBoolSetting)
	if !ok {
		return false, fmt.Sprintf("[type %s]", reflect.TypeOf(s)), fmt.Sprintf("[type %s]", reflect.TypeOf(other))
	}
	return s.SimpleField.Equals(castedOther.SimpleField)
}

func (s *OracleDaoUintSetting) Equals(other core.IEquatable) (bool, string, string) {
	castedOther, ok := other.(*OracleDaoUintSetting)
	if !ok {
		return false, fmt.Sprintf("[type %s]", reflect.TypeOf(s)), fmt.Sprintf("[type %s]", reflect.TypeOf(other))
	}
	return s.SimpleField.Equals(castedOther.SimpleField)
}

func (s *OracleDaoCompoundSetting[DataType]) Equals(other core.IEquatable) (bool, string, string) {
	castedOther, ok := other.(*OracleDaoCompoundSetting[DataType])
	if !ok {
		return false, fmt.Sprintf("[type %s]", reflect.TypeOf(s)), fmt.Sprintf("[type %s]", reflect.TypeOf(other))
	}
	return s.FormattedUint256Field.Equals(castedOther.FormattedUint256Field)
}

func (s *OracleDaoBoolSetting) Copy(other core.ICopyable) {
	castedOther, ok := other.(*OracleDaoBoolSetting)
	if !ok {
		panic(fmt.Sprintf("wrong type; expected %s, got %s", reflect.TypeOf(s), reflect.TypeOf(other)))
	}
	castedOther.SimpleField.Copy(s.SimpleField)
}

func (s *OracleDaoUintSetting) Copy(other core.ICopyable) {
	castedOther, ok := other.(*OracleDaoUintSetting)
	if !ok {
		panic(fmt.Sprintf("wrong type; expected %s, got %s", reflect.TypeOf(s), reflect.TypeOf(other)))
	}
	castedOther.SimpleField.Copy(s.SimpleField)
}

func (s *OracleDaoCompoundSetting[DataType]) Copy(other core.ICopyable) {
	castedOther, ok := other.(*OracleDaoCompoundSetting[DataType])
	if !ok {
		panic(fmt.Sprintf("wrong type; expected %s, got %s", reflect.TypeOf(s), reflect.TypeOf(other)))
	}
	castedOther.FormattedUint256Field.Copy(s.FormattedUint256Field)
}
