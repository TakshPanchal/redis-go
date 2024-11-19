package main

import (
	"fmt"
	"io"
	"net"
)

type Client struct {
	Conn  net.Conn
	Redis *Redis
}

func (c *Client) Close() {
	err := c.Conn.Close()
	if err != nil {
		_ = fmt.Errorf(err.Error())
		return
	}
}

func (c *Client) handle() {
	for {
		decoder := RESPDecoder{Reader: c.Conn}
		req, err := decoder.Decode()
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection: ", err.Error())
				c.Close()
			}
			fmt.Println(err)
			break
		}
		fmt.Println(req)
		command, ok := req.(Array)
		if !ok {
			// TODO: Write a proper error message
			c.Conn.Write([]byte("Error"))
			return
		}

		resp, err := handleRequestCommand(command)
		if err != nil {
			// TODO: Write a proper error message
			fmt.Println("Error executing command: ", err.Error())
			c.Conn.Write([]byte("Error"))
			return
		}
		fmt.Printf("Response: %v\n", resp)

		// Send response to the client
		enc := RESPEncoder{Writer: c.Conn}
		err = enc.Encode(resp)
		if err != nil {
			fmt.Println("Error encoding response: ", err.Error())
		}
	}
}
