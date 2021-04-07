package envutils

import (
	"os"
	"testing"
)

type student struct {
	Name   string `env:""`
	Age    int
	Gender bool `env:""`
	Addr   addr
}

type addr struct {
	Home   string `env:"home"`
	School string `env:"school"`
}

func Test_marshal(t *testing.T) {
	stu := student{
		Name:   "zhangsan2",
		Age:    20,
		Gender: false,
		Addr: addr{
			Home:   "sichuan",
			School: "chengdu",
		},
	}

	// stu02 := student{}

	config := &struct {
		Stud01  student
		Stud02  *student
		Address addr
	}{
		Stud01: stu,
		// Stud02: &stu02,
	}

	// m := make(map[string]interface{})
	// _ = marshal(config, m, "APP")

	// output(m)
	b, _ := Marshal(config, "APP")
	_ = Output(b, os.Stdout)
}
