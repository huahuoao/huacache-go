package client

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"math/rand"
)

type MyStruct struct {
	Name string
	Age  int
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 随机生成5KB的字符串

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func handleSet(c *Client, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(1 * time.Second)
	c.BindGroup("huahuo")

	// 存储单个结构体
	c.SetStruct("test", &MyStruct{Name: "huahuo", Age: 20})

	c.SetString("test_string", "hello world")

	// 存储结构体数组
	myStructArray := []MyStruct{
		{Name: "struct1", Age: 21},
		{Name: "struct2", Age: 22},
	}
	c.SetStruct("myStructArray", myStructArray) // 假设有 SetArray 方法

	// 存储字符串数组
	stringArray := []string{"string1", "string2", "string3"}
	c.SetStruct("stringArray", stringArray) // 假设有 SetStringArray 方法

	// 获取单个结构体
	var res MyStruct
	c.GetStruct("test", &res)

	// 获取字符串
	s, err := c.GetString("test_string")
	if err != nil {
		fmt.Println("获取字符串 'test_string' 出现错误:", err)
	}

	// 获取结构体数组
	var resArray []MyStruct
	c.GetStruct("myStructArray", &resArray) // 假设有 GetGetArray 方法

	// 获取字符串数组
	var strArray []string
	c.GetStruct("stringArray", &strArray) // 假设有 GetGetStringArray 方法

	// 输出结果
	fmt.Printf("res: %v\n", res)
	fmt.Printf("string: %v\n", s)
	fmt.Printf("resArray: %v\n", resArray)
	fmt.Printf("strArray: %v\n", strArray)

	fmt.Println("handleSet 函数执行完毕")
}

func singleTest(c *Client, wg *sync.WaitGroup) {
	defer wg.Done() // 在函数结束时调用 Done，表示该 goroutine 完成
	for i := 0; i < 1; i++ {
		key := randomString(10)
		value := randomString(10 * MB)
		err := c.SetString(key, value)
		if err != nil {
			fmt.Println("设置字符串出现错误:", err)
		}
		value, err = c.GetString(key)
		if err != nil {
			fmt.Println("获取字符串出现错误:", err)
		}
	}
}

func testThousandSetGetCostTime(c *Client, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	c.NewGroup("huahuo", 128*MB)
	c.BindGroup("huahuo")
	startTime := time.Now() // 记录开始时间

	var wg sync.WaitGroup // 创建 WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)             // 增加waitgroup计数
		go singleTest(c, &wg) // 传递waitgroup到goroutine
	}

	wg.Wait() // 等待所有goroutine完成

	elapsedTime := time.Since(startTime) // 计算耗时
	fmt.Printf("耗时: %s\n", elapsedTime)  // 输出耗时
}
func testAllAPI(c *Client, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(100 * time.Millisecond)
	c.NewGroup("huahuo", 128*MB)
	c.BindGroup("huahuo")
	c.SetString("test_string", "hello world")
	c.SetString("test_string1", "hello world2")
	c.SetStruct("test", &MyStruct{Name: "huahuo", Age: 20})
	c.Del("test_string")
	c.NewGroup("huahuo1", 128*MB)
}

var cli *Client

func NewClient() {
	c := New("127.0.0.1", 9000)
	cli = c
	c.Start()
}
func Test(t *testing.T) {
	go NewClient()
	time.Sleep(1000 * time.Millisecond)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go testAllAPI(cli, &wg)
	go handleSet(cli, &wg)
	go testThousandSetGetCostTime(cli, &wg)
	wg.Wait()
}
