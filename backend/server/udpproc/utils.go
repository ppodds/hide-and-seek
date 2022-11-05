package udpproc

import (
	"errors"
	"github.com/ppodds/hide-and-seek/server"
	"github.com/ppodds/hide-and-seek/server/rpc"
	"google.golang.org/protobuf/proto"
	"net"
)

func unmarshalData(ctx *server.UDPContext, msg proto.Message) error {
	if ctx.Data == nil {
		return errors.New("client doesn't provide request data")
	}
	err := proto.Unmarshal(ctx.Data, msg)
	return err
}

func sendRes(ctx *server.UDPContext, addr *net.UDPAddr, msg proto.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = rpc.SendUDPRes(ctx.Conn, addr, buf)
	return err
}
