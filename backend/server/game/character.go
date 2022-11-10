package game

import (
	"github.com/ppodds/hide-and-seek/protos"
	"sync"
)

type CharacterType int

const (
	GHOST CharacterType = iota
	PLAYER
)

type Character struct {
	charType CharacterType
	dead     bool
	pos      *Vector3
	velocity *Vector3
	sync.RWMutex
}

func NewCharacter() *Character {
	character := new(Character)
	character.charType = PLAYER
	character.pos = &Vector3{X: 78.68, Y: 23.71, Z: 44.98}
	character.velocity = new(Vector3)
	return character
}

func (character *Character) FromProtobuf(v *protos.Character) {
	character.pos = ProtobufToVector3(v.Pos)
	character.velocity = ProtobufToVector3(v.Velocity)
}

func (character *Character) MarshalProtoBuf() (*protos.Character, error) {
	pos, err := character.pos.MarshalProtoBuf()
	if err != nil {
		return nil, err
	}
	velocity, err2 := character.velocity.MarshalProtoBuf()
	if err != nil {
		return nil, err2
	}
	charType := protos.CharacterType_PLAYER
	if character.charType == GHOST {
		charType = protos.CharacterType_GHOST
	}
	return &protos.Character{
		Type:     charType,
		Dead:     character.dead,
		Pos:      pos,
		Velocity: velocity,
	}, nil
}
