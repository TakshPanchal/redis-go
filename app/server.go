package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

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

		go handleConnection(conn)
	}
}

// Handle Connected Clients
func handleConnection(conn net.Conn) {
	buff := make([]byte, 1024)

	for {
		// Read request msg from the client
		if _, err := conn.Read(buff); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			conn.Close()
		}
		// Send response to the client
		conn.Write([]byte("+PONG\r\n"))
	}
}
