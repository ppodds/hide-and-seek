package server

import "net"

type TCPContext struct {
	App  *App
	Conn *net.Conn
	Data []byte
}
