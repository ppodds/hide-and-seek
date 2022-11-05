package player

import (
	"github.com/ppodds/hide-and-seek/protos"
	"net"
	"sync"
)

type Player struct {
	ID      uint32 `json:"id"`
	tcpConn *net.TCPConn
	udpConn *net.UDPConn
	udpAddr *net.UDPAddr
	sync.RWMutex
}

func NewPlayer(id uint32, tcpConn *net.TCPConn) *Player {
	player := new(Player)
	player.ID = id
	player.tcpConn = tcpConn
	return player
}

func (player *Player) TCPConn() *net.TCPConn {
	player.RLock()
	defer player.RUnlock()
	return player.tcpConn
}

func (player *Player) UDPConn() *net.UDPConn {
	player.RLock()
	defer player.RUnlock()
	return player.udpConn
}

func (player *Player) UDPAddr() *net.UDPAddr {
	player.RLock()
	defer player.RUnlock()
	return player.udpAddr
}

func (player *Player) SetUDPConn(conn *net.UDPConn) {
	player.Lock()
	defer player.Unlock()
	player.udpConn = conn
}

func (player *Player) SetUDPAddr(addr *net.UDPAddr) {
	player.Lock()
	defer player.Unlock()
	player.udpAddr = addr
}

func (player *Player) MarshalProtoBuf() (*protos.Player, error) {
	player.RLock()
	defer player.RUnlock()
	return &protos.Player{Id: player.ID}, nil
}
