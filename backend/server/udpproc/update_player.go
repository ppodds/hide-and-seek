package udpproc

import (
	"errors"
	"fmt"
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server"
	"time"
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
	// check if ghost win
	liveCount := 0
	for _, p := range game.Players() {
		if !p.Character().Dead() {
			liveCount++
		}
	}
	ghostWin := liveCount == 1
	playerWin := time.Since(game.StartFrom()).Minutes() > 3
	var data *protos.GameBroadcast
	if ghostWin || playerWin {
		var winner protos.CharacterType
		if ghostWin {
			winner = protos.CharacterType_GHOST
		} else {
			winner = protos.CharacterType_PLAYER
		}
		data = &protos.GameBroadcast{
			Event:  protos.GameEvent_GAME_OVER,
			Winner: &winner,
		}
		lobby, ok := ctx.App.Lobbies.Lobbies()[game.LobbyID()]
		if !ok {
			fmt.Println("failed to get game lobby")
			return nil
		}
		lobby.SetInGame(false)
		ctx.App.Games.RmGame(game.ID())
	} else {
		playerProto, err := player.MarshalProtoBuf()
		if err != nil {
			return err
		}
		data = &protos.GameBroadcast{
			Event:  protos.GameEvent_UPDATE_PLAYER,
			Player: playerProto,
		}
	}
	// broadcast
	for _, p := range game.Players() {
		err = sendRes(ctx, p.Player().UDPAddr(), data)
		if err != nil {
			fmt.Println("skip broadcast to player", p.Player().ID, "because", err)
			continue
		}
	}
	return nil
}

func (updatePlayer *UpdatePlayer) ErrorHandler(procErr error, ctx *server.UDPContext) error {
	fmt.Println(procErr)
	return nil
}
