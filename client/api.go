package client

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
)

func (c *Client) SetStruct(key string, value interface{}) error {
	// Log the key and value being processed
	log.Printf("Setting key: %s, value: %v", key, value)

	// Serialize key and value to calculate total size
	valueBytes := SonicSerialize(value)

	// Check if the total size exceeds the limit

	request := &Bluebell{
		Command: SET_KEY,
		Key:     key,
		Value:   valueBytes,
		Group:   c.Group,
	}

	log.Printf("Sending request to server: %+v", request)
	go c.sendRequestToServer(request)

	res, err := c.waitForResponseWithTimeout(5 * time.Second) // 等待响应，设置超时
	if err != nil {
		return err
	}
	if res.Code != SUCCESS {
		return errors.New(string(res.Result))
	}

	log.Println("Set success")
	return nil
}

func (c *Client) SetString(key string, value string) error {
	// Serialize key and value to calculate total size
	b := []byte(value)

	request := &Bluebell{
		Command: SET_KEY,
		Key:     key,
		Value:   b,
		Group:   c.Group,
	}

	go c.sendRequestToServer(request)
	res, err := c.waitForResponseWithTimeout(5 * time.Second) // 等待响应，设置超时
	if err != nil {
		return err
	}
	if res.Code != SUCCESS {
		return errors.New(string(res.Result))
	}
	return nil
}

func (c *Client) GetStruct(key string, obj interface{}) error {
	request := &Bluebell{
		Command: GET_KEY,
		Key:     key,
		Value:   nil,
		Group:   c.Group,
	}

	go c.sendRequestToServer(request)
	res, err := c.waitForResponseWithTimeout(5 * time.Second) // 等待响应，设置超时
	if err != nil {
		return err
	}

	if res.Code != SUCCESS {
		return errors.New(string(res.Result))
	}
	sonic.Unmarshal(res.Result, obj)
	return nil
}

func (c *Client) GetString(key string) (string, error) {
	request := &Bluebell{
		Command: GET_KEY,
		Key:     key,
		Value:   nil,
		Group:   c.Group,
	}

	go c.sendRequestToServer(request)
	res, err := c.waitForResponseWithTimeout(5 * time.Second) // 等待响应，设置超时
	if err != nil {
		return "", err
	}

	if res.Code != SUCCESS {
		return "", errors.New(string(res.Result))
	}
	return string(res.Result), nil
}

func (c *Client) Del(key string) error {
	request := &Bluebell{
		Command: DEL_KEY,
		Key:     key,
		Value:   nil,
		Group:   c.Group,
	}

	go c.sendRequestToServer(request)
	res, err := c.waitForResponseWithTimeout(5 * time.Second) // 等待响应，设置超时
	if err != nil {
		return err
	}

	if res.Code != SUCCESS {
		return errors.New(string(res.Result))
	}
	return nil
}

func (c *Client) NewGroup(groupName string, size int64) error {
	s := strconv.FormatInt(size, 10)
	request := &Bluebell{
		Command: NEW_GROUP,
		Key:     s,
		Value:   nil,
		Group:   groupName,
	}

	go c.sendRequestToServer(request)

	res, err := c.waitForResponseWithTimeout(5 * time.Second) // 等待响应，设置超时
	if err != nil {
		log.Printf("等待响应时发生错误: %v", err) // 日志：记录错误
		return err
	}

	if res.Code != SUCCESS {
		errMsg := string(res.Result)
		fmt.Print(errMsg)
		return errors.New(errMsg)
	}

	return nil
}

func (c *Client) DelGroup(groupName string) error {
	request := &Bluebell{
		Command: DEL_GROUP,
		Key:     groupName,
		Value:   nil,
		Group:   "",
	}

	go c.sendRequestToServer(request)
	res, err := c.waitForResponseWithTimeout(5 * time.Second) // 等待响应，设置超时
	if err != nil {
		return err
	}

	if res.Code != SUCCESS {
		return errors.New(string(res.Result))
	}
	return nil
}
