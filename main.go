package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"huacache/client/client"
	"huacache/client/serialize"
	"io"
	"log"
	"math/big"
	"net"
	"sync" // 导入 sync 包
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return ""
		}
		result[i] = letterBytes[num.Int64()]
	}
	return string(result)
}

func newClient(wg *sync.WaitGroup) {
	defer wg.Done()

	serverAddr := "127.0.0.1:9000"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer conn.Close()

	count := 10
	for i := 0; i < count; i++ {
		msg := &serialize.Bluebell{
			Command: client.SET_KEY,
			Key:     generateRandomString(10),
			Value:   []byte(generateRandomString(30)),
			Group:   "huahuo",
		}
		data, err := msg.Encode()
		if err != nil {
			log.Printf("Failed to serialize message: %v", err)
			continue
		}
		_, err = conn.Write(data)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			continue
		}
	}

	buffer := bytes.NewBuffer([]byte{})
	for {
		inBuffer := make([]byte, 1024)
		n, err := conn.Read(inBuffer)
		if err != nil && err != io.EOF {
			log.Printf("Failed to read response: %v", err)
			break
		}
		if n > 0 {
			buffer.Write(inBuffer[:n])
		}

		for buffer.Len() >= 4 {
			header := buffer.Bytes()[:4]
			messageLength := binary.BigEndian.Uint32(header)
			if uint32(buffer.Len()) < messageLength+4 {
				break
			}

			buffer.Next(4)
			message := buffer.Next(int(messageLength))

			fmt.Printf("Received message: %s\n", string(message))
		}

		if err == io.EOF {
			break
		}
	}
	fmt.Println("Client is shutting down")
}

func main() {
	var wg sync.WaitGroup // 创建 WaitGroup

	goroutines := 100
	for i := 0; i < goroutines; i++ {
		wg.Add(1)         // 增加等待计数
		go newClient(&wg) // 将 WaitGroup 传递给 Goroutine
	}

	wg.Wait() // 等待所有 Goroutine 完成
}
