package game

import "sync"

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
