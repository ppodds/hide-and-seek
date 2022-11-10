package main

import (
	"github.com/ppodds/hide-and-seek/server"
	"github.com/ppodds/hide-and-seek/server/tcpproc"
	"github.com/ppodds/hide-and-seek/server/udpproc"
)

func bootstrap() *server.App {
	app := server.NewApp()
	app.AddTCPProc(new(tcpproc.Login))
	app.AddTCPProc(new(tcpproc.LobbyList))
	app.AddTCPProc(new(tcpproc.CreateLobby))
	app.AddTCPProc(new(tcpproc.JoinLobby))
	app.AddTCPProc(new(tcpproc.LeaveLobby))
	app.AddTCPProc(new(tcpproc.Logout))
	app.AddTCPProc(new(tcpproc.StartGame))
	app.AddUDPProc(new(udpproc.ConnectLobby))
	app.AddUDPProc(new(udpproc.ConnectGame))
	app.AddUDPProc(new(udpproc.UpdatePlayer))
	return app
}

func main() {
	app := bootstrap()
	app.Start()
}
