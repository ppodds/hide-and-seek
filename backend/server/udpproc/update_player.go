package udpproc

import (
	"errors"
	"fmt"
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server"
)

type UpdatePlayer struct {
}

func (updatePlayer *UpdatePlayer) Proc(ctx *server.UDPContext) error {
	req := new(protos.UpdatePlayerRequest)
	err := unmarshalData(ctx, req)
	if err != nil {
		return err
	}
	game, ok := ctx.App.Games.Games()[req.Game.Id]
	if !ok {
		return errors.New("invalid game id")
	}
	player, ok := game.Players()[req.Player.Player.Id]
	if !ok {
		return errors.New("invalid player id")
	}
	player.SetCharacter(req.Player.Character)
	playerProto, err2 := player.MarshalProtoBuf()
	if err2 != nil {
		return err2
	}
	data := &protos.UpdatePlayerBroadcast{
		Player: playerProto,
	}
	// broadcast
	for _, p := range game.Players() {
		if p.Player().ID == player.Player().ID {
			continue
		}
		err2 = sendRes(ctx, p.Player().UDPAddr(), data)
		if err2 != nil {
			fmt.Println("skip broadcast to player", p.Player().ID, "because", err2)
			continue
		}
	}
	return nil
}

func (updatePlayer *UpdatePlayer) ErrorHandler(procErr error, ctx *server.UDPContext) error {
	fmt.Println(procErr)
	return nil
}
