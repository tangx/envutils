package envutils

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseEnv(v interface{}, m map[string]interface{}, prefix string) error {

	rvPtr := reflect.ValueOf(v)

	// 判断是否为所需目标
	if rvPtr.Kind() != reflect.Ptr && rvPtr.Elem().Kind() != reflect.Struct {
		msg := fmt.Sprintf("want a struct prt , got a %#v", rvPtr.Kind())
		return fmt.Errorf(msg)
	}

	// 获取实际 struct 对象
	rv := reflect.Indirect(rvPtr)
	rt := reflect.TypeOf(v).Elem()

	// 遍历 struct
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		/*
			1. 先判断 field 是否为结构体， 以便循环迭代
		*/
		// 如果 field kind 为 struct 指针， 获取真实对象
		if fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}
		// 如果 kind 为 struct， 循环
		if fv.Kind() == reflect.Struct {
			base := strings.Join([]string{prefix, ft.Name}, "__")
			_ = ParseEnv(fv.Addr().Interface(), m, base)
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
		key := join(prefix, name)

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

func join(basename string, fieldname string) string {
	return strings.Join([]string{basename, fieldname}, "_")
}

func output(m map[string]interface{}) {
	data, err := yaml.Marshal(m)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(data)
	_, _ = buf.WriteTo(os.Stdout)
}
