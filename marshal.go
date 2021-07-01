package envutils

import (
	"fmt"
	"reflect"
	"strings"
)

// marshal 将结构体中的 env tag 绑定到 map 中。
func marshal(rv reflect.Value, m map[string]interface{}, prefix string) error {

	// 判断是否为所需目标
	if rv.Kind() != reflect.Struct {
		msg := fmt.Sprintf("want a struct , got a %#v", rv.Kind())
		return fmt.Errorf(msg)
	}

	// 获取实际 struct 对象
	// rt := reflect.TypeOf(v).Elem()
	// rt := Deref(reflect.TypeOf(v))
	rt := deref(rv.Type())

	/*
		遍历 struct
		Field 为 指针， 且为 nil 是， 不会进行初始化。 因为无法后端真实结构体。
	*/
	for i := 0; i < rt.NumField(); i++ {
		// Field ValueOf
		fv := reflect.Indirect(rv.Field(i))
		// 如果 fv 是 unexported, 小写,私有
		// https://golang.org/pkg/reflect/#Value.CanInterface
		if fv.IsValid() && !fv.CanInterface() {
			continue
		}

		// Field TypeOf
		ft := rt.Field(i)

		/*
			2. 再判断 field 字段的实际类型， 以免无 env tag 的字段被略过
		*/
		// 判断是否存在 env TAG， 且是否有效
		name, ok := ft.Tag.Lookup("env")
		if name == "-" {
			continue
		}

		// name 默认值
		if len(name) == 0 {
			name = ft.Name
		}

		/*
			1. 先判断 field 是否为结构体， 以便循环迭代
		*/
		// 如果 field kind 为 struct 指针， 获取真实对象
		// 如果 kind 为 struct， 循环
		if fv.Kind() == reflect.Struct {
			// struct 结构图嵌套使用 双下划线
			subprefix := strings.Join([]string{prefix, name}, "__")
			// _ = marshal(fv.Addr().Interface(), m, subprefix)
			_ = marshal(fv, m, subprefix)
		}

		if !ok {
			continue
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
