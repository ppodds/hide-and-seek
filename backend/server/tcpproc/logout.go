package tcpproc

import (
	"fmt"
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server"
)

type Logout struct {
}

func (logout *Logout) Proc(ctx *server.TCPContext) error {
	req := new(protos.LogoutRequest)
	err := unmarshalData(ctx, req)
	if err != nil {
		return err
	}
	ctx.App.Players.RmPlayer(req.Player)
	return nil
}

func (logout *Logout) ErrorHandler(procErr error, ctx *server.TCPContext) error {
	fmt.Println(procErr)
	return nil
}
