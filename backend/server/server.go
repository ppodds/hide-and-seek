package server

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/ppodds/hide-and-seek/server/game"
	"github.com/ppodds/hide-and-seek/server/lobby"
	"github.com/ppodds/hide-and-seek/server/player"
	"github.com/ppodds/hide-and-seek/server/rpc"
)

type App struct {
	sync.RWMutex
	tcpProcs   [256]TCPProc
	udpProcs   [256]UDPProc
	tcpProcNum byte
	udpProcNum byte
	Lobbies    *lobby.Lobbies
	Players    *player.Players
	Games      *game.Games
}

func NewApp() *App {
	app := new(App)
	app.udpProcNum = 0
	app.Lobbies = lobby.NewLobbys()
	app.Players = player.NewPlayers()
	app.Games = game.NewGames()
	return app
}

func (app *App) HandleTcpProc(conn *net.TCPConn) {
	buf := make([]byte, 5)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx, err := rpc.ParseCall(buf)
	if err != nil {
		fmt.Printf("unsupport protocal. error: %s\n", err.Error())
		return
	}
	if ctx.ContentLength != 0 {
		buf = make([]byte, ctx.ContentLength)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		buf = nil
	}
	if !(ctx.ProcID < app.tcpProcNum) {
		return
	}
	tcpCtx := TCPContext{app, conn, buf}
	fmt.Println("Invoke TCP Proc", ctx.ProcID)
	err = app.tcpProcs[ctx.ProcID].Proc(&tcpCtx)
	if err != nil {
		err := app.tcpProcs[ctx.ProcID].ErrorHandler(err, &tcpCtx)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (app *App) HandleUdpProc(conn *net.UDPConn) {
	buf := make([]byte, 4096)
	_, udpAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("failed to read UDP msg because of", err)
		return
	}
	ctx, err := rpc.ParseCall(buf[:5])
	if err != nil {
		fmt.Println(err)
		return
	}
	if !(ctx.ProcID <= app.udpProcNum) {
		return
	}
	var data []byte
	if ctx.ContentLength != 0 {
		data = buf[5 : 5+ctx.ContentLength]
	} else {
		data = nil
	}
	udpCtx := UDPContext{app, conn, udpAddr, data}
	fmt.Println("Invoke UDP Proc", ctx.ProcID)
	err = app.udpProcs[ctx.ProcID].Proc(&udpCtx)
	if err != nil {
		err := app.udpProcs[ctx.ProcID].ErrorHandler(err, &udpCtx)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (app *App) AddTCPProc(proc TCPProc) {
	app.tcpProcs[app.tcpProcNum] = proc
	app.tcpProcNum++
}

func (app *App) AddUDPProc(proc UDPProc) {
	app.udpProcs[app.udpProcNum] = proc
	app.udpProcNum++
}

func (app *App) Start() {
	host := flag.String("host", "localhost", "host")
	procPort := flag.String("tcpproc-port", "23455", "procedure port")
	gamePort := flag.String("game-port", "23456", "game port")

	flag.Parse()

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
		conn, err := tcpServer.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		go app.HandleTcpProc(conn)
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
