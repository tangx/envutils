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

// SetDefaults 调用 Init 和 SetDefualts 方法设置默认值。
func SetDefaults(v interface{}) error {
	rv := reflect.ValueOf(v)
	return setDefaults(rv)
}
