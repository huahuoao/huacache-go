package main

import (
	"fmt"
	"huacache/client/client"
	"math/rand"
	"strings"
	"time"
)

// RandomStringGenerator 生成指定长度的随机字符串
func RandomStringGenerator(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	sb.Grow(length)

	// 创建一个新的 rand.Rand 实例，使用当前时间作为种子
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rng.Intn(len(charset))])
	}
	return sb.String()
}

// 导入 sync 包
func testClient() {
	time.Sleep(1 * time.Second)
	startTime := time.Now() // 获取开始时间
	c := CacheClient
	c.BindGroup("huahuo")
	count := 1000
	for i := 0; i < count; i++ {
		c.SetKey(RandomStringGenerator(10), []byte{1, 2, 3, 4, 5})
	}
	elapsedTime := time.Since(startTime)  // 计算耗时
	fmt.Printf("执行时间: %s\n", elapsedTime) // 打印执行时间
	client.ThousandTest()
}

var CacheClient *client.Client

func main() {
	go testClient()
	c := client.New("192.168.220.128", 9000)
	CacheClient = c
	c.Start()
}
