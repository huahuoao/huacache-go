package client

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var timerPool = sync.Pool{
	New: func() interface{} {
		return time.NewTimer(0) // 初始化定时器
	},
}

type Client struct {
	ServerAddr string
	ServerPort int
	Conn       net.Conn
	Group      string
	ResponseCh chan []byte
	Buffer     *ByteBuffer // Use ByteBuffer instead of bytes.Buffer
}

func New(serverAddr string, serverPort int) *Client {
	return &Client{
		ServerAddr: serverAddr,
		ServerPort: serverPort,
		Buffer:     &ByteBuffer{},
	}
}

func (c *Client) Start() error {
	c.ResponseCh = make(chan []byte, 1024)
	log.Println("Starting huacache client...")
	addr := fmt.Sprintf("%s:%d", c.ServerAddr, c.ServerPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Connection failed: %v\n", err)
		return err
	}
	c.Conn = conn
	defer conn.Close()

	log.Println("Client started successfully, waiting to receive messages...")
	reader := bufio.NewReader(c.Conn)

	for {
		inBuffer := make([]byte, 1024)
		n, err := reader.Read(inBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading from server: %v\n", err)
			return err
		}

		if n > 0 {
			c.Buffer.Write(inBuffer[:n])
		}

		for c.Buffer.Len() >= 4 {
			header := c.Buffer.Bytes()[:4]
			messageLength := binary.BigEndian.Uint32(header)

			if uint32(c.Buffer.Len()) < messageLength+4 {
				break
			}

			c.Buffer.B = c.Buffer.B[4:] // Skip length bytes
			message := c.Buffer.B[:messageLength]
			c.Buffer.B = c.Buffer.B[messageLength:] // Adjust the buffer

			c.ResponseCh <- message
		}
	}

	log.Println("Client is shutting down")
	return nil
}

func (c *Client) BindGroup(group string) {
	c.Group = group
}

func (c *Client) sendRequestToServer(request *Bluebell) error {
	if c.Conn == nil {
		return errors.New("connection has not been established")
	}
	data, err := request.Encode()
	if err != nil {
		log.Printf("Failed to serialize message: %v", err)
		return err
	}
	_, err = c.Conn.Write(data)
	return err
}

func (c *Client) waitForResponseWithTimeout(timeout time.Duration) (*BluebellResponse, error) {
	timer := timerPool.Get().(*time.Timer)
	timer.Reset(timeout)
	defer func() {
		timer.Stop()
		timerPool.Put(timer)
	}()

	select {
	case response := <-c.ResponseCh:
		res, err := DeserializeResponse(response)
		if err != nil {
			log.Printf("Error during response deserialization: %v", err)
			return nil, err
		}
		return res, nil
	case <-timer.C:
		return nil, errors.New("timeout waiting for response")
	}
}
