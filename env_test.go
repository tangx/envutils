package envutils

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

type student struct {
	Name   string `env:""`
	Age    int
	gender bool  `env:""` // 小写将被忽略
	Addr   *addr // struct, 将跟踪到下层
	addr2  addr
}

func (s *student) Init() {
	s.Name = "caocao"
	s.Age = 100
}

type addr struct {
	Home   string `env:"home"`
	School string `env:"school"`
}

func (a *addr) Init() {
	a.Home = "changshaaaaaaaaa"
}

var (
	_ = os.Setenv("APP__Stud01_Name", "zhugeliang")
	_ = os.Setenv("APP__Stud01_Age", "500")
	_ = os.Setenv("APP__Stud01_Gender", "true")
	_ = os.Setenv("APP__Stud01__addr_home", "addr: sichuan")
	_ = os.Setenv("APP__Stud01__addr2_home", "addr2: sichuan")
	_ = os.Setenv("APP__Addr_home", "zhongguo,sichuan,chengdu")
	_ = os.Setenv("APP__Addr_Home", "APP__Addr_Home")

	APPNAME     = "APP"
	CONFIG_FILE = `config.yml`
)

func Test_marshal(t *testing.T) {
	stu := &student{
		Name:   "zhangsan2",
		Age:    20,
		gender: false,
		Addr: &addr{
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
		Stud01: stu,
		Addr:   addr{},
	}

	/* marshal */
	b, _ := Marshal(config, "APP")
	_ = output(b, os.Stdout)
	os.WriteFile(CONFIG_FILE, b, 0644)
}

func Test_UnmarshalEnv(t *testing.T) {
	stu := student{}
	config := &struct {
		Stud01 *student
		Addr   addr
	}{
		Stud01: &stu,
		Addr:   addr{},
	}

	/* unmarshal */
	err := UnmarshalEnv(config, APPNAME)
	if err != nil {
		log.Fatal(err)
	}

	b, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	_ = output(b, os.Stdout)
}

func Test_UnmarshalFile(t *testing.T) {
	stu := student{}
	config := &struct {
		Stud01 *student
		Addr   addr
	}{
		Stud01: &stu,
		Addr:   addr{},
	}

	/* unmarshal */
	err := UnmarshalFile(config, APPNAME, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	b, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	_ = output(b, os.Stdout)
}

func Test_CallMethod(t *testing.T) {
	stu := &student{
		Addr: &addr{},
	}

	config := &struct {
		Student *student
	}{
		Student: stu,
	}

	err := SetDefaults(config)
	if err != nil {
		panic(err)
	}

	// fmt.Println(config.Student)
	// fmt.Println(config.Student.Addr)

	b, err := Marshal(config, "APP")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}

func output(data []byte, w io.Writer) error {
	_, err := w.Write(data)
	return err
}
