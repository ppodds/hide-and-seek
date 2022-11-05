package server

type UDPProc interface {
	Proc(ctx *UDPContext) error
	ErrorHandler(procErr error, ctx *UDPContext) error
}
