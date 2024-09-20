
# 节点部署

首先部署huacache节点

link：https://github.com/huahuoao/huacache

## Docker
```shell
git clone https://github.com/huahuoao/huacache
```
```shell
cd huacache
```
```shell
docker build -t huacache .
```
```shell
docker run -itd -p 9000:9000 huacache
```
# 快速入门

## 包导入

```shell
go get github.com/huahuoao/huacache-go
```

## 引入项目

huacache.go

```go
package main

import h "github.com/huahuoao/huacache-go/client"

func InitHuaCache() {
	c := h.New("localhost", 8080)
	err := c.Start()
	if err != nil {
		panic(err)
	}
}

```

main.go

```go
package main

func main() {
	InitHuaCache()
}
```

## API

### New

```go
func New(serverAddr string, serverPort int) *Client
```

创建一个新的 `Client` 实例，使用指定的服务器地址和端口。

**参数：**

- `serverAddr`: 服务器地址（如 "127.0.0.1"）。
- `serverPort`: 服务器端口（如 8080）。

**返回值：**

- 返回一个 `Client` 的指针。

### Start

```go
func (c *Client) Start() error
```

启动客户端，建立与服务器的连接并开始读取消息。

**返回值：**

- 返回一个 `error`，如果连接失败或启动过程中发生错误，返回相应的错误信息。

### BindGroup

```go
func (c *Client) BindGroup(group string)
```

将客户端绑定到指定的组。

**参数：**

- `group`: 要绑定的组名。

### SetStruct

```go
func (c *Client) SetStruct(key string, value interface{}) error
```

设置与给定键关联的结构化值。

**参数：**

- `key`: 要设置的键。
- `value`: 与键关联的值，可以是任何类型的结构。

**返回值：**

- 返回一个 `error`，如果操作失败，返回相应的错误信息。

### SetString

```go
func (c *Client) SetString(key string, value string) error
```

设置与给定键关联的字符串值。

**参数：**

- `key`: 要设置的键。
- `value`: 要设置的字符串值。

**返回值：**

- 返回一个 `error`，如果操作失败，返回相应的错误信息。

### GetStruct

```go
func (c *Client) GetStruct(key string, obj interface{}) error
```

检索与给定键关联的结构化值。

**参数：**

- `key`: 要检索的键。
- `obj`: 用于接收返回值的对象，必须是指针类型。

**返回值：**

- 返回一个 `error`，如果操作失败，返回相应的错误信息。

### GetString

```go
func (c *Client) GetString(key string) (string, error)
```

检索与给定键关联的字符串值。

**参数：**

- `key`: 要检索的键。

**返回值：**

- 返回对应的字符串值和一个 `error`，如果操作失败，返回相应的错误信息。

### Del

```go
func (c *Client) Del(key string) error
```

删除与给定键关联的值。

**参数：**

- `key`: 要删除的键。

**返回值：**

- 返回一个 `error`，如果操作失败，返回相应的错误信息。

### NewGroup

```go
func (c *Client) NewGroup(groupName string, size int64) error
```

创建一个新的组，指定组名和大小。

**参数：**

- `groupName`: 新组的名称。
- `size`: 新组的大小。

**返回值：**

- 返回一个 `error`，如果操作失败，返回相应的错误信息。

### DelGroup

```go
func (c *Client) DelGroup(groupName string) error
```

删除指定的组。

**参数：**

- `groupName`: 要删除的组名。

**返回值：**

- 返回一个 `error`，如果操作失败，返回相应的错误信息。

## 示例

```go
package main

import (
	"fmt"
	"testing"

	h "github.com/huahuoao/huacache-go/client"
)

type MyStruct struct {
	Name string
	Age  int
}

func TestApi(t *testing.T) {
	client := h.New("localhost", 9000)
	err := client.Start()
	if err!= nil {
		t.Error(err)
	}
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
	s, err := client.GetString("key")
	if err != nil {
		t.Error(err)
	}
	var targetStruct MyStruct
	//获取struct类型值
	err = client.GetStruct("key2", &targetStruct)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(s)
	fmt.Println(targetStruct)
	//删除key
	client.Del("key")
	//删除group
	client.DelGroup("test_group")
}

```

