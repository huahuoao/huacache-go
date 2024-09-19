# 快速入门
link：https://github.com/huahuoao/huacache
## 包导入

```shell
go get github.com/huahuoao/huacache-go
```

## 引入项目

huacache.go

```go
package main

import h "github.com/huahuoao/huacache-go/client"

var HCache *h.Client

func InitHuaCache() {
	c := h.New("localhost", 8080)
	HCache = c
	c.Start()
}

```

main.go

```go
package main

func main() {
	InitHuaCache()
}
```

## 具体使用样例

```go
package main

import (
	"fmt"
	"testing"
	"time"

	h "github.com/huahuoao/huacache-go/client"
)

var client *h.Client

func NewClient() {
	client = h.New("localhost", 9000)
	client.Start()
}

type MyStruct struct {
	Name string
	Age  int
}

func TestApi(t *testing.T) {
	go NewClient()
	//等待客户端连接
	time.Sleep(1 * time.Second)
	//新建组
	client.NewGroup("test_group", 128*h.MB)
	//绑定组
	client.BindGroup("test_group")
	//设置string类型值
	client.SetString("key", "string_value")
	myStruct := &MyStruct{Name: "huahuoao", Age: 25}
	//设置struct类型值
	client.SetStruct("key2", myStruct)
	//获取string类型值
	s, _ := client.GetString("key")
	var targetStruct MyStruct
	//获取struct类型值
	client.GetStruct("key2", &targetStruct)
	fmt.Println(s)
	fmt.Println(targetStruct)
	//删除key
	client.Del("key")
	//删除group
	client.DelGroup("test_group")
}

```

