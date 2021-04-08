package envutils

import (
	"os"
	"testing"

	"log"

	"gopkg.in/yaml.v3"
)

type student struct {
	Name   string `env:""`
	Age    int
	Gender bool `env:""`
	Addr   addr `env:"addr"`
}

type addr struct {
	Home   string `env:"home"`
	School string `env:"school"`
}

func Test_marshal(t *testing.T) {
	os.Setenv("APP__Stud01_Name", "zhugeliang")
	os.Setenv("APP__Stud01_Age", "500")
	os.Setenv("APP__Stud01_Gender", "true")
	os.Setenv("APP__Stud01__Addr_home", "sichuan")
	os.Setenv("APP__Addr_home", "zhongguo,sichuan,chengdu")
	os.Setenv("APP__Addr_Home", "APP__Addr_Home")
	stu := student{
		// Name:   "zhangsan2",
		// Age:    20,
		// Gender: false,
		// Addr: addr{
		// 	Home:   "sichuan",
		// 	School: "chengdu",
		// },
	}

	// stu02 := student{}

	config := &struct {
		Stud01 *student
		// Person student
		Addr addr
		// Stud02  *student
		// Address addr
	}{
		Stud01: &stu,
		// Person: stu,
		Addr: addr{},
		// Stud02: &stu02,
	}

	// m := make(map[string]interface{})
	// _ = marshal(config, m, "APP")

	// output(m)
	var b []byte
	var err error
	// b, _ = Marshal(config, "APP")
	// _ = Output(b, os.Stdout)

	err = unmarshal(config, "APP")
	// fmt.Println(err)
	if err != nil {
		log.Fatal(err.Error())
	}
	b, err = yaml.Marshal(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	_ = Output(b, os.Stdout)
}
