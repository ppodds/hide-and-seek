package server

import (
	"net"
)

type TCPContext struct {
	App  *App
	Conn *net.TCPConn
	Data []byte
}
