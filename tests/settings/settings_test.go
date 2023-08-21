package settings_test

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/settings"
)

func TestBoostrapFunctions(t *testing.T) {
	// Revert to the baseline at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToBaseline()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to baseline snapshot: %w", err))
		}
	})

	// Initializers
	odao, err := settings.NewOracleDaoSettings(rp)
	if err != nil {
		t.Fatal(fmt.Errorf("error creating odao settings binding: %w", err))
	}
	pdao, err := settings.NewProtocolDaoSettings(rp)
	if err != nil {
		t.Fatal(fmt.Errorf("error creating pdao settings binding: %w", err))
	}
	err = createDefaults(mgr)
	if err != nil {
		t.Fatal("error creating defaults: %w", err)
	}

	// Get all of the current settings
	err = rp.Query(func(mc *batch.MultiCaller) error {
		odao.GetAllDetails(mc)
		pdao.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all initial details: %w", err))
	}

	// Verify details
	CompareDetails(t, &odaoDefaults, &odao.Details)
	CompareDetails(t, &pdaoDefaults, &pdao.Details)

	/*
		// Bootstrap Oracle DAO members
		dnt.BootstrapMember()

		// Bootstrap a contract upgrade
		dnt.BootstrapUpgrade()
	*/
}

// Compares two details structs to ensure their fields all have the same values
func CompareDetails[objType any](t *testing.T, expected *objType, actual *objType) {
	expectedVal := reflect.ValueOf(expected).Elem()
	actualVal := reflect.ValueOf(actual).Elem()
	compareImpl(t, expectedVal, actualVal, expectedVal.Type().Name())
}

// Detail comparison implementation
func compareImpl(t *testing.T, expected reflect.Value, actual reflect.Value, header string) {
	refType := expected.Type()
	fieldCount := refType.NumField()

	for i := 0; i < fieldCount; i++ {
		field := refType.Field(i)
		childExpected := expected.Field(i)
		childActual := actual.Field(i)

		// Try casting to parameters first
		expectedParam, isIParameter := childExpected.Interface().(core.IParameter)
		expectedUint8Param, isIUint8Parameter := childExpected.Interface().(core.IUint8Parameter)

		childSame := true
		if isIParameter {
			// Handle parameters
			actualParam := childActual.Interface().(core.IParameter)
			if expectedParam.GetRawValue() == nil {
				t.Errorf("field %s.%s of type %s - expected was nil", header, field.Name, field.Type.Name())
			} else if actualParam.GetRawValue() == nil {
				t.Errorf("field %s.%s of type %s - actual was nil", header, field.Name, field.Type.Name())
			} else {
				childSame = expectedParam.GetRawValue().Cmp(actualParam.GetRawValue()) == 0
			}
		} else if isIUint8Parameter {
			// Handle uint8 parameters
			actualUint8Param := childActual.Interface().(core.IUint8Parameter)
			childSame = expectedUint8Param.GetRawValue() == actualUint8Param.GetRawValue()
		} else if field.Type.Kind() == reflect.Struct {
			// Handle other nested structs
			compareImpl(t, childExpected, childActual, fmt.Sprintf("%s.%s", header, field.Name))
			continue
		} else {
			// Handle primitives
			switch expectedVal := childExpected.Interface().(type) {
			case *big.Int:
				actualVal := childActual.Interface().(*big.Int)
				if expectedVal == nil {
					t.Errorf("field %s.%s (big.Int) - expected was nil", header, field.Name)
				} else if actualVal == nil {
					t.Errorf("field %s.%s (big.Int) - actual was nil", header, field.Name)
				} else {
					childSame = expectedVal.Cmp(actualVal) == 0
				}
			case bool:
				childSame = expectedVal == childActual.Interface().(bool)
			default:
				t.Fatalf("unexpected type %s in field %s.%s", field.Type.Name(), header, field.Name)
			}
		}

		if !childSame {
			t.Errorf("%s.%s differed; expected %v but got %v", header, field.Name, childExpected.Interface(), childActual.Interface())
		}
	}
}
