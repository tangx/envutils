# envutils
go env utils


## Usage

1. 读取 struct 并将配置文件保存在 config.yml 中

```go

func dump() {

	server := &Server{
		Address: "192.168.100.100",
	}

	config := &struct {
		Server *Server
	}{
		Server: server,
	}

	err := envutils.CallSetDefaults(config)
	if err != nil {
		panic(err)
	}

	b, err := envutils.Marshal(config, appname)
	if err != nil {
		panic(err)
	}
	_ = os.WriteFile(cfgfile, b, os.ModePerm)
}
```

2. 查看保存文件

```bash
# cat config.yml 
AppName__Server_address: 192.168.100.100
AppName__Server_port: 80
```

3. 从文件中读取配置

```go

func read() {

	server := &Server{
		Address: "0.0.0.0",
	}

	config := &struct {
		Server *Server
	}{
		Server: server,
	}

	err := envutils.UnmarshalFile(config, appname, cfgfile)
	if err != nil {
		panic(err)
	}

	fmt.Println("addr=", config.Server.Address)
	fmt.Println("port=", config.Server.Port)
}
// addr= 192.168.100.100
// port= 80
```

## Todo

+ [x] 将结构体的 tag 转换为 config.yml
+ [x] 从环境变量赋值结构体
+ [x] 结构体字段支持 tag 名称
