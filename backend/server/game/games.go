package game

import (
	"github.com/ppodds/hide-and-seek/server/player"
	"math/rand"
	"sync"
	"time"
)

type Games struct {
	games map[uint32]*Game
	curID uint32
	sync.RWMutex
}

func NewGames() *Games {
	games := new(Games)
	games.games = make(map[uint32]*Game)
	games.curID = 1
	return games
}

func (games *Games) CreateGame(players []*player.Player) *Game {
	mapPlayers := make(map[uint32]*Player)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	picked := r.Intn(len(players))
	var ghost *Player
	for i, p := range players {
		gamePlayer := NewPlayer(p)
		mapPlayers[p.ID] = gamePlayer
		if i == picked {
			ghost = gamePlayer
			gamePlayer.character.charType = GHOST
		}
	}
	games.Lock()
	defer games.Unlock()
	game := NewGame(games.curID, mapPlayers, ghost)
	games.games[games.curID] = game
	games.curID++
	return game
}

func (games *Games) Games() map[uint32]*Game {
	games.RLock()
	defer games.RUnlock()
	return games.games
}
