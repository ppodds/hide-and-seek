package rpc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type TCPProcContext struct {
	ProcID        uint8
	ContentLength uint32
}

type UDPProcContext struct {
	ProcID        uint8
	ContentLength uint32
	Addr          *net.UDPAddr
}

func ParseTCPCall(buf []byte) (*TCPProcContext, error) {
	if len(buf) != 5 {
		return nil, errors.New("wrong size buffer")
	}
	ctx := new(TCPProcContext)
	ctx.ProcID = buf[0]
	ctx.ContentLength = binary.LittleEndian.Uint32(buf[1:5])
	return ctx, nil
}

func ParseUDPCall(buf []byte, addr *net.UDPAddr) (*UDPProcContext, error) {
	if len(buf) != 5 {
		return nil, errors.New("wrong size buffer")
	}
	ctx := new(UDPProcContext)
	ctx.ProcID = buf[0]
	ctx.ContentLength = binary.LittleEndian.Uint32(buf[1:5])
	ctx.Addr = addr
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
