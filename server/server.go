package server

import (
	"flag"
	"fmt"
	"github.com/ppodds/hide-and-seek-server/server/lobby"
	"github.com/ppodds/hide-and-seek-server/server/player"
	"log"
	"net"
	"os"
	"sync"
)

type App struct {
	sync.RWMutex
	tcpProcs    [256]func(ctx TCPContext)
	udpProcs    [256]func(ctx UDPContext)
	tcpProcNum  byte
	udpProcNum  byte
	Lobbies     *lobby.Lobbies
	tcpListener net.TCPConn
	udpConn     net.UDPConn
	Players     *player.Players
}

func NewApp() *App {
	app := new(App)
	app.udpProcNum = 0
	app.Lobbies = lobby.NewLobbys()
	app.Players = player.NewPlayers()
	return app
}

func (app *App) HandleTcpProc(conn *net.Conn) {
	for {
		ctx := TCPContext{app, conn, make([]byte, 1024)}
		_, err := (*conn).Read(ctx.Data)
		if err != nil {
			if err.Error() == "EOF" {
				continue
			}
			fmt.Println("failed to read TCP msg because of ", err.Error())
			return
		}
		procId := ctx.Data[0]
		if !(procId <= app.tcpProcNum) {
			return
		}
		app.tcpProcs[procId](ctx)
	}
}

func (app *App) HandleUdpProc(conn *net.UDPConn) {
	ctx := UDPContext{app, conn, make([]byte, 1024)}
	_, _, err := conn.ReadFromUDP(ctx.Data)
	if err != nil {
		fmt.Println("failed to read UDP msg because of ", err.Error())
		return
	}
	procId := ctx.Data[0]
	if !(procId <= app.udpProcNum) {
		return
	}
	go app.udpProcs[procId](ctx)
}

func (app *App) AddTCPProc(proc func(ctx TCPContext)) {
	app.tcpProcs[app.tcpProcNum] = proc
	app.tcpProcNum++
}

func (app *App) AddUDPProc(proc func(ctx UDPContext)) {
	app.udpProcs[app.udpProcNum] = proc
	app.udpProcNum++
}

func (app *App) Start() {
	host := flag.String("host", "localhost", "host")
	procPort := flag.String("tcpproc-port", "23455", "procedure port")
	gamePort := flag.String("game-port", "23456", "game port")

	tcpServer := startTCPServer(host, procPort)
	udpServer := startUDPServer(host, gamePort)

	defer closeServer(tcpServer)
	defer closeServer(udpServer)

	go func() {
		for {
			app.HandleUdpProc(udpServer)
		}
	}()

	for {
		conn, err := tcpServer.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go app.HandleTcpProc(&conn)
	}
}

type server interface {
	Close() error
}

func closeServer(server server) {
	err := server.Close()
	if err != nil {
		fmt.Println("Error closing:", err)
		os.Exit(1)
	}
}

func startTCPServer(host *string, port *string) *net.TCPListener {
	addr, err := net.ResolveTCPAddr("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	fmt.Printf("Start listen on %s:%s\n", *host, *port)
	return listener
}

func startUDPServer(host *string, port *string) *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", *host+":"+*port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	fmt.Printf("Start listen on %s:%s\n", *host, *port)
	return conn
}
