package envutils

import (
	"testing"
)

type student struct {
	Name   string `env:""`
	Age    int
	Gender bool `env:""`
	Addr   addr
}

type addr struct {
	Home   string
	School string `env:""`
}

func TestParseEnv(t *testing.T) {
	stu := student{
		Name:   "zhangsan",
		Age:    20,
		Gender: false,
		Addr: addr{
			Home:   "sichuan",
			School: "chengdu",
		},
	}

	config := &struct {
		Student *student
	}{
		Student: &stu,
	}

	m := make(map[string]interface{})
	_ = ParseEnv(config, m, "APP")

	output(m)
}
