package envutils

import (
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

// Marshal 将结构体转换成对应 []byte
func Marshal(v interface{}, prefix string) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("want a struct ptr, got a %#v", rv.Kind())
	}

	m := make(map[string]interface{})

	// 获取 v 底层数据结构
	rv = reflect.Indirect(rv)

	err := marshal(rv, m, prefix)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(m)
}

// UnmarshalEnv 从环境变量中赋值结构体
func UnmarshalEnv(v interface{}, prefix string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr && rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("want a struct ptr, got a %v", rv.Kind())
	}

	return unmarshalEnv(rv, prefix)
}

// UnmarshalFile 从环境变量中赋值结构体
func UnmarshalFile(v interface{}, prefix string, file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr && rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("want a struct ptr, got a %v", rv.Kind())
	}

	return unmarshalFile(rv, prefix, data)
}

// CallSetDefaults 调用 SetDefualts 设置对象的字段默认值。
func CallSetDefaults(v interface{}) error {
	rv := reflect.ValueOf(v)
	return callSetDefaults(rv)
}

func callSetDefaults(rv reflect.Value) error {
	return methodCaller(rv, "SetDefaults")
}

// CallInitialize to initilize the target which already have variables
// 初始化对象。
// 要求对象字段赋值。 在初始化时，可以进行对象的健康检查， 例如数据库连接。
func CallInitialize(v interface{}) error {
	rv := reflect.ValueOf(v)
	return methodCaller(rv, "Initialize")
}

// CallMethods 调用自定义方法名字
func CallMethods(v interface{}, methods ...string) error {
	rv := reflect.ValueOf(v)
	return methodCaller(rv, methods...)
}

/* 现在 SetDefaults 和 Init 交给用户自己调用*/
// CallInit 调用 SetDefaults 和 Init 方法
// func CallInit(v interface{}) error {
// 	rv := reflect.ValueOf(v)
// 	return callInit(rv)
// }

// func callInit(rv reflect.Value) error {
// 	return methodCaller(rv, "SetDefaults", "Init")
// }
