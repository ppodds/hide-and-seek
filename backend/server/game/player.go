package game

import (
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server/player"
)

type Player struct {
	player    *player.Player
	character *Character
}

func NewPlayer(player *player.Player) *Player {
	p := new(Player)
	p.player = player
	p.character = NewCharacter()
	return p
}

func (player *Player) Player() *player.Player {
	return player.player
}

func (player *Player) Character() *Character {
	return player.character
}

func (player *Player) SetCharacter(character *protos.Character) {
	player.character.FromProtobuf(character)
}

func (player *Player) MarshalProtoBuf() (*protos.GamePlayer, error) {
	p, err := player.player.MarshalProtoBuf()
	if err != nil {
		return nil, err
	}
	c, err2 := player.character.MarshalProtoBuf()
	if err2 != nil {
		return nil, err
	}
	return &protos.GamePlayer{
		Player:    p,
		Character: c,
	}, nil
}
