package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

// Bluebell 消息结构
type Bluebell struct {
	Command string
	Key     string // 键，通常是用于标识数据的字符串
	Value   []byte // 值，存储数据的字节数组
	Group   string // 组，表示消息所属的组或类别
}

func (b *Bluebell) String() string {
	return fmt.Sprintf("Bluebell{\n  Command: %s,\n  Key: %s,\n  Value: %s,\n  Group: %s\n}",
		b.Command,
		b.Key,
		string(b.Value), // 将 []byte 转换为 string
		b.Group,
	)
}

// 序列化：将 Bluebell 结构体序列化为二进制
func (b *Bluebell) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Command 字段
	if err := writeString(buf, b.Command); err != nil {
		return nil, err
	}

	// Key 字段
	if err := writeString(buf, b.Key); err != nil {
		return nil, err
	}

	// Value 字段
	if err := writeBytes(buf, b.Value); err != nil {
		return nil, err
	}

	// Group 字段
	if err := writeString(buf, b.Group); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// writeString 将字符串以长度+内容的形式写入到缓冲区
func writeString(buf *bytes.Buffer, s string) error {
	length := uint32(len(s))
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := buf.Write([]byte(s))
	return err
}

// writeBytes 将 []byte 以长度+内容的形式写入到缓冲区
func writeBytes(buf *bytes.Buffer, b []byte) error {
	length := uint32(len(b))
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := buf.Write(b)
	return err
}

func main() {
	// 服务器地址和端口
	serverAddr := "127.0.0.1:9000"

	// 创建一个连接到服务器
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 创建一个 Bluebell 消息
	msg := &Bluebell{
		Command: "SET",
		Key:     "exampleKey",
		Value:   []byte("exampleValue"),
		Group:   "exampleGroup",
	}

	// 序列化消息
	data, err := msg.Serialize()
	if err != nil {
		log.Fatalf("Failed to serialize message: %v", err)
	}

	// 发送消息到服务器
	_, err = conn.Write(data)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	fmt.Println("Message sent to server:", data)
	time.Sleep(20 * time.Second)
}
