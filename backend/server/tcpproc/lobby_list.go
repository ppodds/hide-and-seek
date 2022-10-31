package tcpproc

import (
	"encoding/json"
	"fmt"
	"github.com/ppodds/hide-and-seek/server"
)

func LobbyList(ctx server.TCPContext) {
	data, err := json.Marshal(ctx.App.Lobbies.Lobbies())
	if err != nil {
		fmt.Printf("JSON marshaling failed: %s", err)
		return
	}
	_, err = (*ctx.Conn).Write(data)
	if err != nil {
		fmt.Println("unable to send message to tcp client")
		return
	}
}
