package tcpproc

import (
	"errors"
	"fmt"
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server"
	"github.com/ppodds/hide-and-seek/server/rpc"
	"google.golang.org/protobuf/proto"
)

type JoinLobby struct {
}

func (joinLobby *JoinLobby) Proc(ctx *server.TCPContext) error {
	req := new(protos.JoinLobbyRequest)
	err := unmarshalData(ctx, req)
	if err != nil {
		return err
	}
	player, ok := ctx.App.Players.Players()[req.Player.Id]
	if !ok {
		return errors.New("invalid player id")
	}
	lobby, ok := ctx.App.Lobbies.Lobbies()[req.Lobby.Id]
	if !ok {
		return errors.New("invalid lobby id")
	}
	// check if player is already in the lobby
	for _, p := range lobby.Players() {
		if p.ID == player.ID {
			return errors.New("player is already in the lobby")
		}
	}
	lobby, err = lobby.AddPlayer(player)
	if err != nil {
		return err
	}
	protoLobby, err2 := lobby.MarshalProtoBuf()
	if err2 != nil {
		return err2
	}
	// send success to client
	t2 := &protos.JoinLobbyResponse{Success: true, Lobby: protoLobby}
	err = sendRes(ctx, t2)
	if err != nil {
		return err
	}
	// broadcast to lobby
	t := &protos.LobbyBroadcast{Event: protos.LobbyEvent_JOIN, Lobby: protoLobby}
	data, err := proto.Marshal(t)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, p := range lobby.Players() {
		if p.ID == player.ID {
			continue
		}
		err = rpc.SendUDPRes(p.UDPConn(), p.UDPAddr(), data)
		if err != nil {
			fmt.Println("skip broadcast to", p.UDPAddr(), "because", err)
			continue
		}
	}
	return nil
}

func (joinLobby *JoinLobby) ErrorHandler(procErr error, ctx *server.TCPContext) error {
	fmt.Println(procErr)
	err := sendRes(ctx, &protos.JoinLobbyResponse{Success: false})
	if err != nil {
		return err
	}
	return nil
}
