package main

type Redis struct {
}

func (r *Redis) Pong() []byte {
	return []byte("+PONG\r\n")
}

//func (r *Redis) Ping() []byte {}
