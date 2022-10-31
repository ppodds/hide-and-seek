package tcpproc

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/ppodds/hide-and-seek-server/server"
)

func JoinLobby(ctx server.TCPContext) {
	playerID := binary.LittleEndian.Uint32(ctx.Data[1:5])
	lobbyID := binary.LittleEndian.Uint32(ctx.Data[5:])
	player, ok := ctx.App.Players.Players()[playerID]
	if !ok {
		fmt.Println("Invalid player id")
		return
	}
	lobby, ok := ctx.App.Lobbies.Lobbies()[lobbyID]
	if !ok {
		fmt.Println("Invalid lobby id")
	}
	lobby = lobby.AddPlayer(player)
	if lobby == nil {
		// write error to client
		_, err := (*ctx.Conn).Write([]byte{0})
		if err != nil {
			fmt.Println("unable to send message to tcp client")
			return
		}
		return
	}
	data, err := json.Marshal(lobby)
	if err != nil {
		fmt.Printf("JSON marshaling failed: %s", err)
		return
	}
	for _, p := range lobby.Players() {
		_, err = (*p.Conn()).Write(data)
		if err != nil {
			fmt.Println("unable to send message to tcp client")
			return
		}
	}
	return
}
