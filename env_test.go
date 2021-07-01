package envutils

import (
	"log"
	"os"
	"testing"

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

var (
	_ = os.Setenv("APP__Stud01_Name", "zhugeliang")
	_ = os.Setenv("APP__Stud01_Age", "500")
	_ = os.Setenv("APP__Stud01_Gender", "true")
	_ = os.Setenv("APP__Stud01__addr_home", "addr: sichuan")
	_ = os.Setenv("APP__Stud01__addr2_home", "addr2: sichuan")
	_ = os.Setenv("APP__Addr_home", "zhongguo,sichuan,chengdu")
	_ = os.Setenv("APP__Addr_Home", "APP__Addr_Home")

	APPNAME = "APP"
)

func Test_marshal(t *testing.T) {
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
	config := &struct {
		Stud01 *student
		Addr   addr
	}{
		Stud01: &stu,
		Addr:   addr{},
	}

	/* marshal */
	b, _ := Marshal(config, "APP")
	_ = Output(b, os.Stdout)

}

func Test_LoadEnv(t *testing.T) {
	stu := student{}
	config := &struct {
		Stud01 *student
		Addr   addr
	}{
		Stud01: &stu,
		Addr:   addr{},
	}

	LoadEnv(config, APPNAME)
	/* unmarshal */
	err := LoadEnv(config, APPNAME)
	if err != nil {
		log.Fatal(err)
	}

	b, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	_ = Output(b, os.Stdout)
}
