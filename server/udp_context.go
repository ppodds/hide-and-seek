package server

import "net"

type UDPContext struct {
	App  *App
	Conn *net.UDPConn
	Data []byte
}
