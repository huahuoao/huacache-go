package client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"huacache/client/serialize"
	"io"
	"log"
	"net"
)

type Client struct {
	ServerAddr string
	ServerPort int
	Conn       net.Conn
	Group      string
}

func New(serverAddr string, serverPort int) *Client {
	return &Client{
		ServerAddr: serverAddr,
		ServerPort: serverPort}
}

func (c *Client) Start() error {
	fmt.Println("Begin starting huacache client...")
	addr := fmt.Sprintf("%s:%d", c.ServerAddr, c.ServerPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return err // 返回连接错误
	}
	c.Conn = conn
	defer func() {
		fmt.Println("Closing connection")
		conn.Close() // 确保连接在函数结束时关闭
	}()
	fmt.Println("Start client successfully")
	fmt.Println("Connected to server " + addr + " wait to receive message")
	buffer := bytes.NewBuffer([]byte{})
	for {
		inBuffer := make([]byte, 1024)
		n, err := conn.Read(inBuffer)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Failed to read response: %v\n", err)
			}
			break // 发生错误时退出循环
		}

		if n > 0 {
			buffer.Write(inBuffer[:n])
		}

		// 处理缓冲区中的消息
		for buffer.Len() >= 4 {
			// 检查缓冲区的长度是否足够接收完整的消息头
			header := buffer.Bytes()[:4]
			messageLength := binary.BigEndian.Uint32(header)

			// 检查缓冲区的长度是否足够接收完整的消息
			if uint32(buffer.Len()) < messageLength+4 {
				break
			}

			// 移动缓冲区，跳过消息头
			buffer.Next(4)
			message := buffer.Next(int(messageLength))

			fmt.Printf("Received message: %s\n", string(message))
		}

		if err == io.EOF {
			break // 连接关闭，退出循环
		}

	}

	fmt.Println("Client is shutting down")
	return nil // 关闭时返回 nil
}

func (c *Client) BindGroup(group string) {
	c.Group = group
}
func (c *Client) SetKey(key string, value []byte) error {
	// 检查连接是否有效
	fmt.Print("Setting key: ")
	fmt.Println(key)
	if c.Conn == nil {
		return errors.New("connection is not established")
	}

	msg := &serialize.Bluebell{
		Command: SET_KEY,
		Key:     key,   // 使用传入的key
		Value:   value, // 使用传入的value
		Group:   c.Group,
	}

	data, err := msg.Encode()
	if err != nil {
		log.Printf("Failed to serialize message: %v", err)
		return err // 返回错误
	}
	log.Printf("Sending message: %v", data)
	_, err = c.Conn.Write(data)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err // 返回错误
	}

	return nil // 无错误时返回nil
}
