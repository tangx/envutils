package envutils

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func unmarshalFile(rv reflect.Value, prefix string, data []byte) error {
	m := make(map[string]string)

	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	return unmarshal(rv, prefix, m)
}

func unmarshalEnv(rv reflect.Value, prefix string) error {
	m := make(map[string]string)
	envs := os.Environ()

	for _, pair := range envs {
		kv := strings.Split(pair, "=")
		m[kv[0]] = kv[1]
	}

	return unmarshal(rv, prefix, m)
}

// unmarshal 从环境变量赋值结构体
func unmarshal(rv reflect.Value, prefix string, m map[string]string) (err error) {
	rv = reflect.Indirect(rv)

	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("want a struct, but got a %#v", rv.Kind().String())
	}

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
			err = unmarshal(fv, subprefix, m)
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
		// val := os.Getenv(key)
		val, ok := m[key]
		if !ok {
			continue
		}

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
