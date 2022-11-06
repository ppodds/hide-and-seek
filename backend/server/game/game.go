package game

type Game struct {
	id      uint32
	players map[uint32]*GamePlayer
	ghost   *GamePlayer
}
