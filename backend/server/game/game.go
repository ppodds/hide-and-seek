package game

import "time"

type Game struct {
	id        uint32
	players   map[uint32]*Player
	ghost     *Player
	startFrom time.Time
}

func NewGame(id uint32, players map[uint32]*Player, ghost *Player) *Game {
	game := new(Game)
	game.id = id
	game.players = players
	game.ghost = ghost
	game.startFrom = time.Now()
	return game
}

func (game *Game) ID() uint32 {
	return game.id
}

func (game *Game) Players() map[uint32]*Player {
	return game.players
}

func (game *Game) Ghost() *Player {
	return game.ghost
}

func (game *Game) StartFrom() time.Time {
	return game.startFrom
}
