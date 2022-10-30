package player

import (
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

func (players *Players) AddPlayer() *Player {
	players.Lock()
	player := NewPlayer(players.curID)
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
