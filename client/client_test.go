package client

import (
	"testing"

	"math/rand"
)

type MyStruct struct {
	Name string
	Age  int
}

// 随机生成5KB的字符串

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Test(t *testing.T) {
	c := New("localhost", 9000)
	err := c.Start()
	if err != nil {
		t.Error(err)
	}
	c.NewGroup("huahuo", 128*MB)
	c.BindGroup("huahuo")
	random := randomString(1024)
	c.SetString("test_string", random)
	s, err := c.GetString("test_string")
	if err != nil {
		t.Error(err)
	}
	c.SetStruct("test", &MyStruct{Name: "huahuo", Age: 20})
	var obj MyStruct
	err = c.GetStruct("test", &obj)
	if err != nil {
		t.Error(err)
	}
	if s != random || obj.Name != "huahuo" || obj.Age != 20 {
		t.Error("test failed")
	}
	c.Del("test_string")
	c.NewGroup("huahuo1", 128*MB)
	c.DelGroup("huahuo1")
}
