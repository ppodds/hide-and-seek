package player

import (
	"github.com/ppodds/hide-and-seek/protos"
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
	players.curID = 1
	return players
}

func (players *Players) AddPlayer(tcpConn *net.TCPConn) *Player {
	players.Lock()
	player := NewPlayer(players.curID, tcpConn)
	players.players[player.ID] = player
	players.curID++
	players.Unlock()
	return player
}

func (players *Players) RmPlayer(player *protos.Player) {
	players.Lock()
	defer players.Unlock()
	delete(players.players, player.Id)
}

func (players *Players) Players() map[uint32]*Player {
	players.RLock()
	defer players.RUnlock()
	return players.players
}
