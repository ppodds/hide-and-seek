package lobby

import (
	"encoding/json"
	"github.com/ppodds/hide-and-seek-server/server/player"
	"strconv"
)

type Lobby struct {
	ID        uint32 `json:"id"`
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

func (lobby *Lobby) MarshalJSON() ([]byte, error) {
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
	if err != nil {
		return nil, err
	}
	return []byte(`{"id":` + id + `,"lead":` + string(lead) + `,"players":` + string(players) + `,"curPeople":` + curPeople + `,"maxPeople":` + maxPeople + `}`), nil
}
