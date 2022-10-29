package player

type Player struct {
	ID uint32
}

func NewPlayer(id uint32) *Player {
	player := new(Player)
	player.ID = id
	return player
}
