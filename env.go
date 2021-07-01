package envutils

import (
	"fmt"
	"io"
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

// LoadEnv 从环境变量中赋值结构体
func LoadEnv(v interface{}, prefix string) (err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("want a struct ptr, got a %v", rv.Kind())
	}

	// 获取所有 key
	m := make(map[string]interface{})

	// 获取 v 底层数据结构
	rv = reflect.Indirect(rv)

	err = marshal(rv, m, prefix)
	if err != nil {
		return err
	}

	// 获取所有环境变量
	for key := range m {
		m[key] = os.Getenv(key)
	}

	return unmarshal(rv, prefix)
}

func output(data []byte, w io.Writer) error {
	_, err := w.Write(data)
	return err
}
