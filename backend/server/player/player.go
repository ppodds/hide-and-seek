package player

import "net"

type Player struct {
	ID   uint32 `json:"id"`
	conn *net.Conn
}

func NewPlayer(id uint32, conn *net.Conn) *Player {
	player := new(Player)
	player.ID = id
	player.conn = conn
	return player
}

func (player *Player) Conn() *net.Conn {
	return player.Conn()
}
