package tcpproc

import (
	"errors"
	"fmt"
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server"
	"github.com/ppodds/hide-and-seek/server/rpc"
	"google.golang.org/protobuf/proto"
)

type StartGame struct {
}

func (startGame *StartGame) Proc(ctx *server.TCPContext) error {
	req := new(protos.StartGameRequest)
	err := unmarshalData(ctx, req)
	if err != nil {
		return err
	}
	lobby, ok := ctx.App.Lobbies.Lobbies()[req.Lobby.Id]
	if !ok {
		return errors.New("invalid lobby")
	}
	if lobby.Lead().ID != req.Player.Id {
		return errors.New("not the lobby lead")
	}
	if lobby.InGame() {
		return errors.New("game is already started")
	}
	lobby.SetInGame(true)
	game := ctx.App.Games.CreateGame(lobby.Players())
	players := make(map[uint32]*protos.GamePlayer)
	for _, p := range game.Players() {
		player, err2 := p.Player().MarshalProtoBuf()
		if err2 != nil {
			return err2
		}
		character, err3 := p.Character().MarshalProtoBuf()
		if err3 != nil {
			return err3
		}
		players[p.Player().ID] = &protos.GamePlayer{
			Player:    player,
			Character: character,
		}
	}
	// send success response to client
	res := &protos.StartGameResponse{Success: true}
	err = sendRes(ctx, res)
	if err != nil {
		return err
	}
	// broadcast
	for _, p := range game.Players() {
		broadcast := &protos.LobbyBroadcast{
			Event: protos.LobbyEvent_START,
			InitGame: &protos.InitGame{
				Game:    &protos.Game{Id: game.ID()},
				Players: players,
			},
		}
		data, err3 := proto.Marshal(broadcast)
		if err3 != nil {
			fmt.Println("skip broadcast to player", p.Player().ID, "because", err3)
			continue
		}
		err = rpc.SendUDPRes(p.Player().UDPConn(), p.Player().UDPAddr(), data)
		if err != nil {
			fmt.Println("skip broadcast to", p.Player().UDPAddr(), "because", err)
			continue
		}
	}
	return nil
}

func (startGame *StartGame) ErrorHandler(procErr error, ctx *server.TCPContext) error {
	fmt.Println(procErr)
	res := &protos.StartGameResponse{Success: false}
	err := sendRes(ctx, res)
	if err != nil {
		return err
	}
	return nil
}
