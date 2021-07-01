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
	gender bool `env:""` // 小写将被忽略
	Addr   addr // struct, 将跟踪到下层
	addr2  addr
}

type addr struct {
	Home   string `env:"home"`
	School string `env:"school"`
}

func Test_marshal(t *testing.T) {
	os.Setenv("APP__Stud01_Name", "zhugeliang")
	os.Setenv("APP__Stud01_Age", "500")
	os.Setenv("APP__Stud01_Gender", "true")
	os.Setenv("APP__Stud01__addr_home", "addr: sichuan")
	os.Setenv("APP__Stud01__addr2_home", "addr2: sichuan")
	os.Setenv("APP__Addr_home", "zhongguo,sichuan,chengdu")
	os.Setenv("APP__Addr_Home", "APP__Addr_Home")
	stu := student{
		Name:   "zhangsan2",
		Age:    20,
		gender: false,
		Addr: addr{
			Home:   "sichuan",
			School: "chengdu",
		},
		addr2: addr{
			Home:   "sichuan",
			School: "chengdu",
		},
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

	/* public */
	var b []byte
	var err error

	/* marshal */
	b, _ = Marshal(config, "APP")
	_ = Output(b, os.Stdout)

	/* unmarshal */
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
