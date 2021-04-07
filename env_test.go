package envutils

import (
	"testing"
)

type student struct {
	Name   string
	Age    int
	Gender bool
	Addr   addr
}

type addr struct {
	Home   string
	School string
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
