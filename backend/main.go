package main

import (
	"flag"
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
	return app
}

func main() {
	flag.Parse()

	app := bootstrap()
	app.Start()
}
