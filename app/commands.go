package main

import (
	"errors"
	"fmt"
)

// Supported Commands
const (
	PING = "PING"
	ECHO = "ECHO"
)

func handleRequestCommand(cmds Array) (interface{}, error) {
	cmd, ok := cmds[0].([]byte)
	if !ok {
		return nil, errors.New("unable to resolve command")
	}
	fmt.Printf("Handle Command: %v\n", cmd)

	switch string(cmd) {
	case PING:
		return "PONG", nil

	case ECHO:
		return echo(cmds[1:])
	default:
		return nil, errors.New("c")
	}
}

func echo(args Array) (BulkString, error) {
	if len(args) == 0 {
		return nil, errors.New("argument required")
	}
	arg, ok := args[0].(BulkString)
	if !ok {
		return nil, errors.New("string Argument required")
	}
	return arg, nil
}
