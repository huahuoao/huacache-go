package client

import (
	"testing"
)

func TestSetKey(t *testing.T) {
	c := New("127.0.0.1", 9000)
	c.BindGroup("huahuo")
	c.SetKey("test", []byte("value"))
}
