package tcpproc

import (
	"fmt"
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server"
)

type Login struct {
}

func (login *Login) Proc(ctx *server.TCPContext) error {
	player := ctx.App.Players.AddPlayer(ctx.Conn)
	err := sendRes(ctx, &protos.Player{Id: player.ID})
	if err != nil {
		return err
	}
	return nil
}

func (login *Login) ErrorHandler(procErr error, ctx *server.TCPContext) error {
	fmt.Println(procErr)
	return nil
}
