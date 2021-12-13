package envutils

import (
	"fmt"
	"reflect"
)

func methodCaller(rv reflect.Value, methods ...string) error {

	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("want a struct prt, got a %v of %v", rv.Kind(), rv.Elem().Kind())
	}

	// call Method: Init and SetDefaults
	for _, method := range methods {
		mv := rv.MethodByName(method)
		if mv.IsValid() {
			mv.Call(nil)
		}
	}

	rv = reflect.Indirect(rv)
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		tag, ok := ft.Tag.Lookup("env")
		if !ok {
			continue
		}
		if tag == "-" {
			continue
		}

		if !fv.CanInterface() {
			continue
		}
		if fv.Kind() != reflect.Ptr || fv.Elem().Kind() != reflect.Struct {
			continue
		}

		if err := methodCaller(fv, methods...); err != nil {
			return err
		}

	}

	return nil
}
