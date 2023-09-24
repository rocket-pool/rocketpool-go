package settings_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/rocket-pool/rocketpool-go/core"
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
		passedCheck := true
		var firstVal string
		var secondVal string
		field := refType.Field(i)
		if !field.IsExported() {
			continue
		}
		childExpected := expected.Field(i)
		childActual := actual.Field(i)

		expectedEquatable, isEquatable := childExpected.Interface().(core.IEquatable)
		if isEquatable {
			// Handle fields
			var same bool
			actualEquatable := childActual.Interface().(core.IEquatable)
			same, firstVal, secondVal = expectedEquatable.Equals(actualEquatable)
			passedCheck = (same == checkIfEqual)
		} else if field.Type.Kind() == reflect.Struct {
			// Handle other nested structs
			passedCheck = compareImpl(log, childExpected, childActual, fmt.Sprintf("%s.%s", header, field.Name), checkIfEqual)
			if !passedCheck {
				valid = false
			}
			continue
		} else {
			logMessage(log, "cannot compare, unexpected type %s in field %s.%s", field.Type.Name(), header, field.Name)
		}

		if !passedCheck {
			valid = false
			if checkIfEqual {
				logMessage(log, "%s.%s differed; expected %v but got %v", header, field.Name, firstVal, secondVal)
			} else {
				logMessage(log, "%s.%s was the same; expected not %v but got %v", header, field.Name, firstVal, secondVal)
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
		if !field.IsExported() {
			continue
		}

		sourceCopyable, isCopyable := childSource.Interface().(core.ICopyable)
		if isCopyable {
			// Handle fields
			destCopyable := childDest.Interface().(core.ICopyable)
			destCopyable.Copy(sourceCopyable)
		} else if field.Type.Kind() == reflect.Struct {
			// Handle other nested structs
			cloneImpl(t, childSource, childDest, fmt.Sprintf("%s.%s", header, field.Name))
			continue
		} else {
			t.Fatalf("cannot clone, unexpected type %s in field %s.%s", field.Type.Name(), header, field.Name)
		}
	}
}
