package main

import (
	"fmt"
	"net"
	"os"
)

//TODO: Build a Redis Encoder

func main() {
	fmt.Println("Logs from your program will appear here!")

	r := &Redis{}
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		c := Client{Conn: conn, Redis: r}
		go c.handle()
	}
}
