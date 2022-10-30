package tcpproc

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/ppodds/hide-and-seek-server/server"
)

func CreateLobby(ctx server.TCPContext) {
	leadID := binary.LittleEndian.Uint32(ctx.Data[1:])
	lead, ok := ctx.App.Players.Players()[leadID]
	if !ok {
		fmt.Println("Invalid player id")
		return
	}
	lobby := ctx.App.Lobbies.AddLobby(lead, 4)
	data, err := json.Marshal(lobby)
	if err != nil {
		fmt.Printf("JSON marshaling failed: %s", err)
		return
	}
	_, err = (*ctx.Conn).Write(data)
	if err != nil {
		fmt.Println("unable to send message to tcp client")
		return
	}
	return
}
