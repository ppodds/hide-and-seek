package tcpproc

import (
	"errors"
	"fmt"
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server"
	"github.com/ppodds/hide-and-seek/server/rpc"
	"google.golang.org/protobuf/proto"
)

type LeaveLobby struct {
}

func (leaveLobby *LeaveLobby) Proc(ctx *server.TCPContext) error {
	req := new(protos.LeaveLobbyRequest)
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
	// check if player is not in the lobby
	inLobby := false
	for _, p := range lobby.Players() {
		if p.ID == player.ID {
			inLobby = true
			break
		}
	}
	if !inLobby {
		err = leaveLobby.leaveFailed(ctx)
		if err != nil {
			return err
		}
		return nil
	}
	lobby, err = lobby.RmPeople(player)
	if err != nil {
		err2 := leaveLobby.leaveFailed(ctx)
		if err2 != nil {
			return err2
		}
		return err
	}
	// remove lobby if lead leave or the lobby member amount = 0
	if lobby.CurPeople() == 0 || lobby.Lead().ID == player.ID {
		ctx.App.Lobbies.RmLobby(lobby.ID)
	}
	// send resp to client
	res := &protos.LeaveLobbyResponse{Success: true}
	buf, err := proto.Marshal(res)
	if err != nil {
		return err
	}
	err = rpc.SendTCPRes(ctx.Conn, buf)
	if err != nil {
		return err
	}
	// broadcast to lobby
	lobbyProto, err2 := lobby.MarshalProtoBuf()
	if err2 != nil {
		return err2
	}
	res2 := &protos.LobbyBroadcast{Event: protos.LobbyEvent_LEAVE, Lobby: lobbyProto}
	buf, err2 = proto.Marshal(res2)
	if err2 != nil {
		return err2
	}
	for _, p := range lobby.Players() {
		err = rpc.SendUDPRes(p.UDPConn(), p.UDPAddr(), buf)
		if err != nil {
			fmt.Println("skip broadcast to", p.UDPAddr(), "because", err)
			continue
		}
	}
	return nil
}

func (leaveLobby *LeaveLobby) ErrorHandler(procErr error, ctx *server.TCPContext) error {
	fmt.Println(procErr)
	return nil
}
func (leaveLobby *LeaveLobby) leaveFailed(ctx *server.TCPContext) error {
	res := &protos.LeaveLobbyResponse{Success: false}
	buf, err := proto.Marshal(res)
	if err != nil {
		return err
	}
	err = rpc.SendTCPRes(ctx.Conn, buf)
	if err != nil {
		return err
	}
	return nil
}
