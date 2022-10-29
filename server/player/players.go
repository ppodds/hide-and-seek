package player

import (
	"sync"
)

type Players struct {
	sync.RWMutex
	curID   uint32
	players []*Player
}

func NewPlayers() *Players {
	players := new(Players)
	players.players = make([]*Player, 100)
	return players
}

func (players *Players) AddPlayer() *Player {
	players.Lock()
	player := NewPlayer(players.curID)
	players.players = append(players.players, player)
	players.curID++
	players.Unlock()
	return player
}
