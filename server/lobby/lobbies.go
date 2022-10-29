package lobby

import (
	"github.com/ppodds/hide-and-seek-server/server/player"
	"sync"
)

type Lobbies struct {
	sync.RWMutex
	curID   uint32
	lobbies []*Lobby
}

func NewLobbys() *Lobbies {
	lobbies := new(Lobbies)
	lobbies.lobbies = make([]*Lobby, 10)
	return lobbies
}

func (lobbies *Lobbies) AddLobby(lead *player.Player, maxNum uint8) *Lobby {
	lobbies.Lock()
	lobby := NewLobby(lobbies.curID, lead, maxNum)
	lobbies.lobbies = append(lobbies.lobbies, lobby)
	lobbies.curID++
	lobbies.Unlock()
	return lobby
}
