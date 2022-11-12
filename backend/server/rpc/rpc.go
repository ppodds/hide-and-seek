package rpc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type RPCContext struct {
	ProcID        uint8
	ContentLength uint32
}

func ParseCall(buf []byte) (*RPCContext, error) {
	if len(buf) != 5 {
		return nil, errors.New("wrong size buffer")
	}
	ctx := new(RPCContext)
	ctx.ProcID = buf[0]
	ctx.ContentLength = binary.LittleEndian.Uint32(buf[1:5])
	return ctx, nil
}

func SendTCPRes(conn *net.TCPConn, buf []byte) error {
	header := make([]byte, 4)
	binary.LittleEndian.PutUint32(header, uint32(len(buf)))
	fmt.Println("send tcp response header: ", header)
	_, err := conn.Write(header)
	if err != nil {
		return err
	}
	fmt.Println("send tcp response: ", buf)
	_, err = conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func SendUDPRes(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) error {
	fmt.Println("send udp response:", buf, "to", addr)
	_, err := conn.WriteToUDP(buf, addr)
	if err != nil {
		return err
	}
	return nil
}
