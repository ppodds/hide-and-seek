package player

import (
	"net"
	"sync"
)

type Players struct {
	sync.RWMutex
	curID   uint32
	players map[uint32]*Player
}

func NewPlayers() *Players {
	players := new(Players)
	players.players = make(map[uint32]*Player)
	return players
}

func (players *Players) AddPlayer(conn *net.Conn) *Player {
	players.Lock()
	player := NewPlayer(players.curID, conn)
	players.players[player.ID] = player
	players.curID++
	players.Unlock()
	return player
}

func (players *Players) Players() map[uint32]*Player {
	players.RLock()
	defer players.RUnlock()
	return players.players
}
