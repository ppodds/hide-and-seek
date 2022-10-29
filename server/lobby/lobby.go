package lobby

import (
	"github.com/ppodds/hide-and-seek-server/server/player"
)

type Lobby struct {
	ID        uint32
	lead      *player.Player
	players   []*player.Player
	curPeople uint8
	maxPeople uint8
}

func NewLobby(id uint32, lead *player.Player, maxPeople uint8) *Lobby {
	lobby := new(Lobby)
	lobby.ID = id
	lobby.lead = lead
	lobby.players = []*player.Player{lead}
	lobby.maxPeople = maxPeople
	lobby.curPeople = 1
	return lobby
}
