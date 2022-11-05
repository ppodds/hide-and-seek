package tcpproc

import (
	"errors"
	"fmt"
	"github.com/ppodds/hide-and-seek/protos"
	"github.com/ppodds/hide-and-seek/server"
)

type CreateLobby struct {
}

func (createLobby *CreateLobby) Proc(ctx *server.TCPContext) error {
	req := new(protos.CreateLobbyRequest)
	err := unmarshalData(ctx, req)
	if err != nil {
		return err
	}
	lead, ok := ctx.App.Players.Players()[req.Lead.Id]
	if !ok {
		return errors.New("invalid player id")
	}
	// check if the player already create a lobby
	check := false
	for _, lobby := range ctx.App.Lobbies.Lobbies() {
		if lobby.Lead().ID == lead.ID {
			check = true
			break
		}
	}
	if check {
		return errors.New("player already created a lobby")
	}
	lobby := ctx.App.Lobbies.AddLobby(lead, 4)
	protoLobby, err := lobby.MarshalProtoBuf()
	if err != nil {
		return err
	}
	// send result to client
	res := &protos.CreateLobbyResponse{Success: true, Lobby: protoLobby}
	err = sendRes(ctx, res)
	if err != nil {
		return err
	}
	return nil
}

func (createLobby *CreateLobby) ErrorHandler(procErr error, ctx *server.TCPContext) error {
	fmt.Println(procErr)
	res := &protos.CreateLobbyResponse{Success: false}
	err := sendRes(ctx, res)
	return err
}
