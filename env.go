package envutils

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Marshal 将结构体转换成对应 []byte
func Marshal(v interface{}, prefix string) ([]byte, error) {
	m := make(map[string]interface{})
	err := marshal(v, m, prefix)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(m)
}

// UnmarshalFromEnv 从环境变量中赋值结构体
func UnmarshalFromEnv(v interface{}, prefix string) (err error) {

	// 获取所有 key
	m := make(map[string]interface{})
	err = marshal(v, m, prefix)
	if err != nil {
		return err
	}

	// 获取所有环境变量
	for key := range m {
		m[key] = os.Getenv(key)
	}

	return
}

func Output(data []byte, w io.Writer) error {
	_, err := w.Write(data)
	return err
}

// marshal 将结构体中的 env tag 绑定到 map 中。
func marshal(v interface{}, m map[string]interface{}, prefix string) error {

	// 获取 v 底层数据结构
	rv := reflect.Indirect(reflect.ValueOf(v))

	// 判断是否为所需目标
	if rv.Kind() != reflect.Struct {
		msg := fmt.Sprintf("want a struct , got a %#v", rv.Kind())
		return fmt.Errorf(msg)
	}

	// 获取实际 struct 对象
	// rt := reflect.TypeOf(v).Elem()
	rt := Deref(reflect.TypeOf(v))

	/*
		遍历 struct
		Field 为 指针， 且为 nil 是， 不会进行初始化。 因为无法后端真实结构体。
	*/
	for i := 0; i < rt.NumField(); i++ {
		// Field ValueOf
		fv := rv.Field(i)
		// Field TypeOf
		ft := rt.Field(i)

		/*
			1. 先判断 field 是否为结构体， 以便循环迭代
		*/
		// 如果 field kind 为 struct 指针， 获取真实对象
		// 如果 kind 为 struct， 循环
		if fv = reflect.Indirect(fv); fv.Kind() == reflect.Struct {
			// struct 结构图嵌套使用 双下划线
			subprefix := strings.Join([]string{prefix, ft.Name}, "__")
			_ = marshal(fv.Addr().Interface(), m, subprefix)
		}

		/*
			2. 再判断 field 字段的实际类型， 以免无 env tag 的字段被略过
		*/
		// 判断是否存在 env TAG， 且是否有效
		var name string
		var ok bool
		if name, ok = ft.Tag.Lookup("env"); !ok || name == "-" {
			continue
		}

		// name 默认值
		if len(name) == 0 {
			name = ft.Name
		}

		// struct 中 field 嵌套使用 单下划线
		key := strings.Join([]string{prefix, name}, "_")

		// 根据实际类型处理
		switch val := fv.Interface().(type) {
		case string:
			m[key] = val
		case int, int8, int16, int32, int64:
			m[key] = val
		case uint, uint8, uint16, uint32, uint64:
			m[key] = val
		case bool:
			m[key] = val
		}
	}

	return nil
}

// unmarshal 从环境变量赋值结构体
func unmarshal(v interface{}, prefix string) (err error) {
	rvPtr := reflect.ValueOf(v)

	if rvPtr.Kind() != reflect.Ptr && rvPtr.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("want a struct ptr, but got a %#v", rvPtr.Kind())
	}

	rv := reflect.Indirect(rvPtr)
	rt := Deref(reflect.TypeOf(v))

	for i := 0; i < rv.NumField(); i++ {

		fv := reflect.Indirect(rv.Field(i))
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
			err = unmarshal(fv.Addr().Interface(), subprefix)
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
