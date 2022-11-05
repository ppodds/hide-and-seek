package lobby

import (
	"encoding/json"
	"errors"
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server/player"
	"strconv"
	"sync"
)

type Lobby struct {
	ID        uint32 `json:"id"`
	lead      *player.Player
	players   []*player.Player
	curPeople uint32
	maxPeople uint32
	sync.RWMutex
}

func NewLobby(id uint32, lead *player.Player, maxPeople uint32) *Lobby {
	lobby := new(Lobby)
	lobby.ID = id
	lobby.lead = lead
	lobby.players = []*player.Player{lead}
	lobby.maxPeople = maxPeople
	lobby.curPeople = 1
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
	return &protos.Lobby{Id: lobby.ID, Lead: lead, Players: players, CurPeople: lobby.curPeople, MaxPeople: lobby.maxPeople}, nil
}

func (lobby *Lobby) MarshalJSON() ([]byte, error) {
	lobby.RLock()
	id := strconv.Itoa(int(lobby.ID))
	lead, err := json.Marshal(lobby.lead)
	if err != nil {
		return nil, err
	}
	players, err := json.Marshal(lobby.players)
	if err != nil {
		return nil, err
	}
	curPeople := strconv.Itoa(int(lobby.curPeople))
	maxPeople := strconv.Itoa(int(lobby.maxPeople))
	lobby.RUnlock()
	if err != nil {
		return nil, err
	}
	return []byte(`{"id":` + id + `,"lead":` + string(lead) + `,"players":` + string(players) + `,"curPeople":` + curPeople + `,"maxPeople":` + maxPeople + `}`), nil
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
	return lobby.curPeople
}

func (lobby *Lobby) Lead() *player.Player {
	return lobby.lead
}

func (lobby *Lobby) Players() []*player.Player {
	return lobby.players
}
