package tcpproc

import (
	"fmt"
	"github.com/ppodds/hide-and-seek/server"
)

type LobbyList struct{}

func (lobbyList *LobbyList) Proc(ctx *server.TCPContext) error {
	data, err := ctx.App.Lobbies.MarshalProtoBuf()
	if err != nil {
		return err
	}
	err = sendRes(ctx, data)
	return err
}

func (lobbyList *LobbyList) ErrorHandler(procErr error, ctx *server.TCPContext) error {
	fmt.Println(procErr)
	return nil
}
