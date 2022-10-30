package main

import (
	"flag"
	"github.com/ppodds/hide-and-seek-server/server"
	"github.com/ppodds/hide-and-seek-server/server/tcpproc"
)

func bootstrap() *server.App {
	app := server.NewApp()
	app.AddTCPProc(tcpproc.Login)
	app.AddTCPProc(tcpproc.LobbyList)
	app.AddTCPProc(tcpproc.CreateLobby)
	return app
}

func main() {
	flag.Parse()

	app := bootstrap()
	app.Start()
}
