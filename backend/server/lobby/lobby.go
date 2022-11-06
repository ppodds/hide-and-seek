package lobby

import (
	"errors"
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server/player"
	"sync"
)

type Lobby struct {
	ID        uint32
	lead      *player.Player
	players   []*player.Player
	curPeople uint32
	maxPeople uint32
	inGame    bool
	sync.RWMutex
}

func NewLobby(id uint32, lead *player.Player, maxPeople uint32) *Lobby {
	lobby := new(Lobby)
	lobby.ID = id
	lobby.lead = lead
	lobby.players = []*player.Player{lead}
	lobby.maxPeople = maxPeople
	lobby.curPeople = 1
	lobby.inGame = false
	return lobby
}

func (lobby *Lobby) MarshalProtoBuf() (*protos.Lobby, error) {
	lobby.RLock()
	defer lobby.RUnlock()
	players := make([]*protos.Player, 0)
	for _, v := range lobby.players {
		data, err := v.MarshalProtoBuf()
		if err != nil {
			return nil, err
		}
		players = append(players, data)
	}
	lead, err := lobby.lead.MarshalProtoBuf()
	if err != nil {
		return nil, err
	}
	return &protos.Lobby{Id: lobby.ID, Lead: lead, Players: players, CurPeople: lobby.curPeople, MaxPeople: lobby.maxPeople, InGame: lobby.inGame}, nil
}

// AddPlayer Add a player into a lobby. Return new lobby if success, else nil.
func (lobby *Lobby) AddPlayer(player *player.Player) (*Lobby, error) {
	lobby.Lock()
	defer lobby.Unlock()
	if lobby.curPeople == lobby.maxPeople {
		return nil, errors.New("lobby is full")
	}
	lobby.players = append(lobby.players, player)
	lobby.curPeople += 1
	return lobby, nil
}

func (lobby *Lobby) RmPeople(player *player.Player) (*Lobby, error) {
	pos := -1
	lobby.Lock()
	defer lobby.Unlock()
	for i := 0; i < len(lobby.players); i++ {
		if lobby.players[i].ID == player.ID {
			pos = i
			break
		}
	}
	if pos == -1 {
		return lobby, errors.New("can't find the player")
	}
	lobby.players = append(lobby.players[:pos], lobby.players[pos+1:]...)
	lobby.curPeople -= 1
	return lobby, nil
}

func (lobby *Lobby) CurPeople() uint32 {
	lobby.RLock()
	defer lobby.RUnlock()
	return lobby.curPeople
}

func (lobby *Lobby) Lead() *player.Player {
	lobby.RLock()
	defer lobby.RUnlock()
	return lobby.lead
}

func (lobby *Lobby) Players() []*player.Player {
	lobby.RLock()
	defer lobby.RUnlock()
	return lobby.players
}

func (lobby *Lobby) InGame() bool {
	lobby.RLock()
	defer lobby.RUnlock()
	return lobby.inGame
}

func (lobby *Lobby) SetInGame(v bool) {
	lobby.Lock()
	defer lobby.Unlock()
	lobby.inGame = v
}
