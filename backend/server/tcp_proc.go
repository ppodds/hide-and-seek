package server

type TCPProc interface {
	Proc(ctx *TCPContext) error
	ErrorHandler(procErr error, ctx *TCPContext) error
}
