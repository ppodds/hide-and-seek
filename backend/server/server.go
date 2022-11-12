package server

import (
	"flag"
	"fmt"
	"github.com/ppodds/hide-and-seek/server/game"
	"github.com/ppodds/hide-and-seek/server/lobby"
	"github.com/ppodds/hide-and-seek/server/player"
	"github.com/ppodds/hide-and-seek/server/rpc"
	"log"
	"net"
	"os"
	"sync"
)

type App struct {
	sync.RWMutex
	tcpProcs        [256]TCPProc
	udpProcs        [256]UDPProc
	tcpProcNum      byte
	udpProcNum      byte
	tcpProcContexts map[*net.TCPConn]*rpc.TCPProcContext
	udpProcContexts map[*net.UDPConn]*rpc.UDPProcContext
	Lobbies         *lobby.Lobbies
	Players         *player.Players
	Games           *game.Games
}

func NewApp() *App {
	app := new(App)
	app.udpProcNum = 0
	app.Lobbies = lobby.NewLobbys()
	app.Players = player.NewPlayers()
	app.Games = game.NewGames()
	app.tcpProcContexts = make(map[*net.TCPConn]*rpc.TCPProcContext)
	app.udpProcContexts = make(map[*net.UDPConn]*rpc.UDPProcContext)
	return app
}

func (app *App) HandleTcpProc(conn *net.TCPConn) {
	for {
		ctx, ok := app.tcpProcContexts[conn]
		if !ok {
			buf := make([]byte, 5)
			_, err := conn.Read(buf)
			if err != nil {
				if err.Error() == "EOF" {
					continue
				}
				fmt.Println("failed to read TCP msg because of", err.Error())
				return
			}
			ctx, err = rpc.ParseTCPCall(buf)
			if err != nil {
				fmt.Printf("unsupport protocal. error: %s\n", err.Error())
				return
			}
			app.tcpProcContexts[conn] = ctx
		} else {
			delete(app.tcpProcContexts, conn)
			var buf []byte
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
			err := app.tcpProcs[ctx.ProcID].Proc(&tcpCtx)
			if err != nil {
				err := app.tcpProcs[ctx.ProcID].ErrorHandler(err, &tcpCtx)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func (app *App) HandleUdpProc(conn *net.UDPConn) {
	ctx, ok := app.udpProcContexts[conn]
	if !ok {
		buf := make([]byte, 5)
		_, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("failed to read UDP msg because of", err.Error())
			return
		}
		ctx, err = rpc.ParseUDPCall(buf, remoteAddr)
		if err != nil {
			fmt.Printf("unsupport protocal. error: %s\n", err.Error())
			return
		}
		app.udpProcContexts[conn] = ctx
	} else {
		delete(app.udpProcContexts, conn)
		var buf []byte
		if ctx.ContentLength != 0 {
			buf = make([]byte, ctx.ContentLength)
			var err error
			_, _, err = conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			buf = nil
		}
		if !(ctx.ProcID <= app.udpProcNum) {
			return
		}
		udpCtx := UDPContext{app, conn, ctx.Addr, buf}
		fmt.Println("Invoke UDP Proc", ctx.ProcID)
		err := app.udpProcs[ctx.ProcID].Proc(&udpCtx)
		if err != nil {
			err := app.udpProcs[ctx.ProcID].ErrorHandler(err, &udpCtx)
			if err != nil {
				fmt.Println(err)
				return
			}
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
