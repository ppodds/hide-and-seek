package udpproc

import (
	"errors"
	"fmt"
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server"
)

type ConnectLobby struct {
}

func (connectLobby *ConnectLobby) Proc(ctx *server.UDPContext) error {
	req := new(protos.ConnectLobbyRequest)
	err := unmarshalData(ctx, req)
	if err != nil {
		return err
	}
	player, ok := ctx.App.Players.Players()[req.Player.Id]
	if !ok {
		return errors.New("invalid player id")
	}
	player.SetUDPConn(ctx.Conn)
	player.SetUDPAddr(ctx.Addr)
	err = sendRes(ctx, ctx.Addr, &protos.ConnectLobbyResponse{Success: true})
	if err != nil {
		return err
	}
	return nil
}

func (connectLobby *ConnectLobby) ErrorHandler(procErr error, ctx *server.UDPContext) error {
	fmt.Println(procErr)
	err := sendRes(ctx, ctx.Addr, &protos.ConnectLobbyResponse{Success: false})
	if err != nil {
		return err
	}
	return nil
}
