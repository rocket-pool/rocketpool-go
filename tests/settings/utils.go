package settings_test

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao"
)

// Compares two details structs to ensure their fields all have the same values
func EnsureSameDetails[objType any](log func(string, ...any), expected *objType, actual *objType) bool {
	expectedVal := reflect.ValueOf(expected).Elem()
	actualVal := reflect.ValueOf(actual).Elem()
	return compareImpl(log, expectedVal, actualVal, expectedVal.Type().Name(), true)
}

// Compares two details structs to ensure their fields all have different values
func EnsureDifferentDetails[objType any](log func(string, ...any), expected *objType, actual *objType) bool {
	expectedVal := reflect.ValueOf(expected).Elem()
	actualVal := reflect.ValueOf(actual).Elem()
	return compareImpl(log, expectedVal, actualVal, expectedVal.Type().Name(), false)
}

// Compares two details structs to ensure their fields all have the same values
func Clone[objType any](t *testing.T, source *objType, dest *objType) {
	sourceVal := reflect.ValueOf(source).Elem()
	destVal := reflect.ValueOf(dest).Elem()
	cloneImpl(t, sourceVal, destVal, sourceVal.Type().Name())
}

// Detail comparison implementation
func compareImpl(log func(string, ...any), expected reflect.Value, actual reflect.Value, header string, checkIfEqual bool) bool {
	refType := expected.Type()
	fieldCount := refType.NumField()

	valid := true
	for i := 0; i < fieldCount; i++ {
		field := refType.Field(i)
		childExpected := expected.Field(i)
		childActual := actual.Field(i)

		// Try casting to parameters first
		expectedParam, isIParameter := childExpected.Addr().Interface().(core.IUint256Parameter)
		expectedUint8Param, isIUint8Parameter := childExpected.Addr().Interface().(core.IUint8Parameter)

		// Try casting to settings next
		expectedBoolSetting, isBoolSetting := childExpected.Interface().(dao.IBoolSetting)
		expectedUintSetting, isUintSetting := childExpected.Interface().(dao.IUintSetting)

		passedCheck := true
		if isIParameter {
			// Handle parameters
			actualParam := childActual.Addr().Interface().(core.IUint256Parameter)
			if expectedParam.GetRawValue() == nil {
				logMessage(log, "field %s.%s of type %s - expected was nil", header, field.Name, field.Type.Name())
			} else if actualParam.GetRawValue() == nil {
				logMessage(log, "field %s.%s of type %s - actual was nil", header, field.Name, field.Type.Name())
			} else {
				if checkIfEqual {
					passedCheck = expectedParam.GetRawValue().Cmp(actualParam.GetRawValue()) == 0
				} else {
					passedCheck = expectedParam.GetRawValue().Cmp(actualParam.GetRawValue()) != 0
				}
			}
		} else if isIUint8Parameter {
			// Handle uint8 parameters
			actualUint8Param := childActual.Addr().Interface().(core.IUint8Parameter)
			if checkIfEqual {
				passedCheck = expectedUint8Param.GetRawValue() == actualUint8Param.GetRawValue()
			} else {
				passedCheck = expectedUint8Param.GetRawValue() != actualUint8Param.GetRawValue()
			}
		} else if isBoolSetting {
			// Handle bool settings
			actualBoolSetting := childActual.Interface().(dao.IBoolSetting)
			if checkIfEqual {
				passedCheck = expectedBoolSetting.GetRawValue() == actualBoolSetting.GetRawValue()
			} else {
				passedCheck = expectedBoolSetting.GetRawValue() != actualBoolSetting.GetRawValue()
			}
		} else if isUintSetting {
			// Handle uint / compound settings
			actualUintSetting := childActual.Interface().(dao.IUintSetting)
			if checkIfEqual {
				passedCheck = expectedUintSetting.GetRawValue().Cmp(actualUintSetting.GetRawValue()) == 0
			} else {
				passedCheck = expectedUintSetting.GetRawValue().Cmp(actualUintSetting.GetRawValue()) != 0
			}
		} else if field.Type.Kind() == reflect.Struct {
			// Handle other nested structs
			passedCheck = compareImpl(log, childExpected, childActual, fmt.Sprintf("%s.%s", header, field.Name), checkIfEqual)
			if !passedCheck {
				valid = false
			}
			continue
		} else {
			// Handle primitives
			switch expectedVal := childExpected.Interface().(type) {
			case *big.Int:
				actualVal := childActual.Interface().(*big.Int)
				if expectedVal == nil {
					logMessage(log, "field %s.%s (big.Int) - expected was nil", header, field.Name)
				} else if actualVal == nil {
					logMessage(log, "field %s.%s (big.Int) - actual was nil", header, field.Name)
				} else {
					if checkIfEqual {
						passedCheck = expectedVal.Cmp(actualVal) == 0
					} else {
						passedCheck = expectedVal.Cmp(actualVal) != 0
					}
				}
			case bool:
				if checkIfEqual {
					passedCheck = expectedVal == childActual.Interface().(bool)
				} else {
					passedCheck = expectedVal != childActual.Interface().(bool)
				}
			default:
				logMessage(log, "cannot compare, unexpected type %s in field %s.%s", field.Type.Name(), header, field.Name)
			}
		}

		if !passedCheck {
			valid = false
			if checkIfEqual {
				logMessage(log, "%s.%s differed; expected %v but got %v", header, field.Name, childExpected.Interface(), childActual.Interface())
			} else {
				logMessage(log, "%s.%s was the same; expected not %v but got %v", header, field.Name, childExpected.Interface(), childActual.Interface())
			}
		}
	}

	return valid
}

func logMessage(log func(string, ...any), format string, args ...any) {
	if log != nil {
		log(format, args...)
	}
}

// Detail cloning implementation
func cloneImpl(t *testing.T, source reflect.Value, dest reflect.Value, header string) {
	refType := source.Type()
	fieldCount := refType.NumField()

	for i := 0; i < fieldCount; i++ {
		field := refType.Field(i)
		childSource := source.Field(i)
		childDest := dest.Field(i)

		// Try casting to parameters first
		sourceParam, isIParameter := childSource.Addr().Interface().(core.IUint256Parameter)
		sourceUint8Param, isIUint8Parameter := childSource.Addr().Interface().(core.IUint8Parameter)

		// Try casting to settings next
		sourceBoolSetting, isBoolSetting := childSource.Interface().(dao.IBoolSetting)
		sourceUintSetting, isUintSetting := childSource.Interface().(dao.IUintSetting)

		if isIParameter {
			// Handle parameters
			destParam := childDest.Addr().Interface().(core.IUint256Parameter)
			if sourceParam.GetRawValue() == nil {
				t.Errorf("field %s.%s of type %s - source was nil", header, field.Name, field.Type.Name())
			} else {
				destParam.SetRawValue(sourceParam.GetRawValue())
			}
		} else if isIUint8Parameter {
			// Handle uint8 parameters
			destUint8Param := childDest.Addr().Interface().(core.IUint8Parameter)
			destUint8Param.SetRawValue(sourceUint8Param.GetRawValue())
		} else if isBoolSetting {
			// Handle bool settings
			destBoolSetting := childDest.Interface().(dao.IBoolSetting)
			destBoolSetting.SetRawValue(sourceBoolSetting.GetRawValue())
		} else if isUintSetting {
			// Handle uint / compound settings
			destUintSetting := childDest.Interface().(dao.IUintSetting)
			destUintSetting.SetRawValue(sourceUintSetting.GetRawValue())
		} else if field.Type.Kind() == reflect.Struct {
			// Handle other nested structs
			cloneImpl(t, childSource, childDest, fmt.Sprintf("%s.%s", header, field.Name))
			continue
		} else {
			// Handle primitives
			switch sourceVal := childSource.Interface().(type) {
			case *big.Int:
				destVal := childDest.Addr().Interface().(**big.Int)
				if sourceVal == nil {
					t.Errorf("field %s.%s (big.Int) - source was nil", header, field.Name)
				} else {
					*destVal = big.NewInt(0).Set(sourceVal)
				}
			case bool:
				destVal := childDest.Addr().Interface().(*bool)
				*destVal = sourceVal
			default:
				t.Fatalf("cannot clone, unexpected type %s in field %s.%s", field.Type.Name(), header, field.Name)
			}
		}
	}
}
