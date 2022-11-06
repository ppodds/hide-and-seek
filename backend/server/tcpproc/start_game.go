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
	// send success response to client
	res := &protos.StartGameResponse{Success: true}
	err = sendRes(ctx, res)
	if err != nil {
		return err
	}
	// broadcast
	broadcast := &protos.LobbyBroadcast{Event: protos.LobbyEvent_START}
	data, err := proto.Marshal(broadcast)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, p := range lobby.Players() {
		err = rpc.SendUDPRes(p.UDPConn(), p.UDPAddr(), data)
		if err != nil {
			fmt.Println("skip broadcast to", p.UDPAddr(), "because", err)
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
