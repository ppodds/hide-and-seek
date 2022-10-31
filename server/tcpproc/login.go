package tcpproc

import (
	"encoding/binary"
	"fmt"
	"github.com/ppodds/hide-and-seek-server/server"
)

func Login(ctx server.TCPContext) {
	player := ctx.App.Players.AddPlayer(ctx.Conn)
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, player.ID)
	_, err := (*ctx.Conn).Write(buf)
	if err != nil {
		fmt.Println("unable to send message to tcp client")
		return
	}
}
