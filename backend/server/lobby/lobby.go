package lobby

import (
	"encoding/json"
	"fmt"
	"github.com/ppodds/hide-and-seek/server/player"
	"strconv"
	"sync"
)

type Lobby struct {
	ID        uint32 `json:"id"`
	lead      *player.Player
	players   []*player.Player
	curPeople uint8
	maxPeople uint8
	sync.RWMutex
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
func (lobby *Lobby) AddPlayer(player *player.Player) *Lobby {
	lobby.Lock()
	if lobby.curPeople == lobby.maxPeople {
		return nil
	}
	lobby.players = append(lobby.players, player)
	lobby.curPeople += 1
	lobby.Unlock()
	return lobby
}

// RmPeople remove a player from a lobby. Return new lobby if success, else nil.
func (lobby *Lobby) RmPeople(player *player.Player) *Lobby {
	pos := -1
	lobby.Lock()
	for i := 0; i < len(lobby.players); i++ {
		if lobby.players[i].ID == player.ID {
			pos = i
			break
		}
	}
	if pos == -1 {
		fmt.Println("Can't find the player")
		return nil
	}
	lobby.players = append(lobby.players[:pos], lobby.players[pos+1:]...)
	lobby.curPeople -= 1
	lobby.Unlock()
	return lobby
}

func (lobby *Lobby) CurPeople() uint8 {
	return lobby.curPeople
}

func (lobby *Lobby) Lead() *player.Player {
	return lobby.lead
}

func (lobby *Lobby) Players() []*player.Player {
	return lobby.players
}
