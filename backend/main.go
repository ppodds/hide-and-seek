package main

import (
	"flag"
	"github.com/ppodds/hide-and-seek/server"
	"github.com/ppodds/hide-and-seek/server/tcpproc"
)

func bootstrap() *server.App {
	app := server.NewApp()
	app.AddTCPProc(tcpproc.Login)
	app.AddTCPProc(tcpproc.LobbyList)
	app.AddTCPProc(tcpproc.CreateLobby)
	app.AddTCPProc(tcpproc.JoinLobby)
	app.AddTCPProc(tcpproc.LeaveLobby)
	return app
}

func main() {
	flag.Parse()

	app := bootstrap()
	app.Start()
}
