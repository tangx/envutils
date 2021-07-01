package envutils

import (
	"fmt"
	"reflect"
)

func setDefaults(rv reflect.Value) error {

	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("want a struct prt, got a %v of %v", rv.Kind(), rv.Elem().Kind())
	}

	// call Method: Init and SetDefaults
	for _, method := range []string{"SetDefaults", "Init"} {
		mv := rv.MethodByName(method)
		if mv.IsValid() {
			mv.Call(nil)
		}
	}

	rv = reflect.Indirect(rv)
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)

		if !fv.CanInterface() {
			continue
		}
		if fv.Kind() != reflect.Ptr || fv.Elem().Kind() != reflect.Struct {
			continue
		}

		if err := setDefaults(fv); err != nil {
			return err
		}

	}

	return nil
}
