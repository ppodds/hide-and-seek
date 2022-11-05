package server

import (
	"net"
)

type UDPContext struct {
	App  *App
	Conn *net.UDPConn
	Addr *net.UDPAddr
	Data []byte
}
