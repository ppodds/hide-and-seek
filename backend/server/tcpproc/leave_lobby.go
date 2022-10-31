package tcpproc

import (
	"encoding/binary"
	"fmt"
	"github.com/ppodds/hide-and-seek/server"
)

func LeaveLobby(ctx server.TCPContext) {
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
	lobby = lobby.RmPeople(player)
	if lobby == nil {
		// write error to client
		_, err := (*ctx.Conn).Write([]byte{0})
		if err != nil {
			fmt.Println("unable to send message to tcp client")
			return
		}
		return
	}
	if lobby.CurPeople() == 0 || lobby.Lead().ID == player.ID {
		ctx.App.Lobbies.RmLobby(lobby.ID)
	}
	_, err := (*ctx.Conn).Write([]byte{1})
	if err != nil {
		fmt.Println("unable to send message to tcp client")
		return
	}
	for _, p := range lobby.Players() {
		_, err := (*p.Conn()).Write([]byte{1})
		if err != nil {
			fmt.Println("unable to send message to tcp client")
			return
		}
	}
	return
}
