package tcpproc

import (
	"errors"
	"github.com/ppodds/hide-and-seek/server"
	"github.com/ppodds/hide-and-seek/server/rpc"
	"google.golang.org/protobuf/proto"
)

func unmarshalData(ctx *server.TCPContext, msg proto.Message) error {
	if ctx.Data == nil {
		return errors.New("client doesn't provide request data")
	}
	err := proto.Unmarshal(ctx.Data, msg)
	return err
}

func sendRes(ctx *server.TCPContext, msg proto.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = rpc.SendTCPRes(ctx.Conn, buf)
	return err
}
