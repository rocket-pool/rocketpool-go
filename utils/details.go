package utils

import (
	"fmt"
	"reflect"

	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

func GetAllDetails(contract any, details any, mc *multicall.MultiCaller) error {
	// Get the value and type of the contract
	contractValue := reflect.ValueOf(contract)
	contractType := reflect.TypeOf(contract)

	// Get the value and type of the details container
	detailsValue := reflect.ValueOf(details).Elem()
	detailsType := reflect.TypeOf(details).Elem()

	// Get a binding for the multicaller
	mcPtr := reflect.ValueOf(mc)

	// Run through each field
	for i := 0; i < detailsType.NumField(); i++ {
		field := detailsType.Field(i)

		// Get the RP struct tag
		getterName := fmt.Sprintf("Get%s", field.Name)
		methodProto, found := contractType.MethodByName(getterName)
		if !found {
			return fmt.Errorf("error getting retriever method for type %s, field %s: struct %s does not have a method named '%s'", detailsType.String(), field.Name, contractType.String(), getterName)
		}

		// Invoke the method
		method := contractValue.Method(methodProto.Index)
		pointerToField := detailsValue.Field(i).Addr()
		method.Call([]reflect.Value{mcPtr, pointerToField})
	}

	return nil
}
