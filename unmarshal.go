package envutils

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// unmarshal 从环境变量赋值结构体
func unmarshal(rv reflect.Value, prefix string) (err error) {

	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("want a struct ptr, but got a %#v", rv.Kind())
	}

	rv = reflect.Indirect(rv)
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {

		fv := reflect.Indirect(rv.Field(i))
		// 如果 fv 是 unexported, 小写,私有
		// https://golang.org/pkg/reflect/#Value.CanInterface
		if fv.IsValid() && !fv.CanInterface() {
			continue
		}

		ft := rt.Field(i)

		name, ok := ft.Tag.Lookup("env")
		// if !ok || name == "-" {
		// 	continue
		// }
		// 如果 env 的值为 - ， 则略过
		if name == "-" {
			continue
		}
		// 如果 name 为空， 则略过
		if len(name) == 0 {
			name = ft.Name
		}

		if fv.Kind() == reflect.Struct {
			subprefix := strings.Join([]string{prefix, name}, "__")
			// fmt.Println("subprefix =", subprefix)
			err = unmarshal(fv, subprefix)
			if err != nil {
				return err
			}
			continue
		}

		// 如果非结构体， 且无 env tag 则略过
		if !ok {
			continue
		}

		key := strings.Join([]string{prefix, name}, "_")
		val := os.Getenv(key)

		// fmt.Printf("key(%s) = value(%s)\n", key, val)

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(val)
		case reflect.Bool:
			b, err := strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("invalid value type key(%s), value(%s) is not", key, val)
			}
			fv.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid value type key(%s), value(%s)", key, val)
			}
			fv.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid value type key(%s), value(%s)", key, val)
			}
			fv.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			fv.SetFloat(x)
		default:
			return fmt.Errorf("unsupported type %v", fv.Type())
		}
	}

	return
}
