package game

import (
	"github.com/ppodds/hide-and-seek/protos"
)

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func ProtobufToVector3(v *protos.Vector3) *Vector3 {
	vec := new(Vector3)
	vec.X = v.X
	vec.Y = v.Y
	vec.Z = v.Z
	return vec
}

func (v *Vector3) MarshalProtoBuf() (*protos.Vector3, error) {
	return &protos.Vector3{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}, nil
}
