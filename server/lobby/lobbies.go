package lobby

import (
	"github.com/ppodds/hide-and-seek-server/server/player"
	"sync"
)

type Lobbies struct {
	sync.RWMutex
	curID   uint32
	lobbies map[uint32]*Lobby
}

func NewLobbys() *Lobbies {
	lobbies := new(Lobbies)
	lobbies.lobbies = make(map[uint32]*Lobby)
	return lobbies
}

func (lobbies *Lobbies) AddLobby(lead *player.Player, maxNum uint8) *Lobby {
	lobbies.Lock()
	lobby := NewLobby(lobbies.curID, lead, maxNum)
	lobbies.lobbies[lobby.ID] = lobby
	lobbies.curID++
	lobbies.Unlock()
	return lobby
}

func (lobbies *Lobbies) Lobbies() map[uint32]*Lobby {
	lobbies.RLock()
	defer lobbies.RUnlock()
	return lobbies.lobbies
}
