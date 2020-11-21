package server

import (
	"github.com/astaxie/beego/logs"
	"github.com/soheilhy/cmux"
	"log"
	"net"
)

type Srv struct {
	Port  string
	Token string
}

func (s *Srv) Run() {
	if s.Port == "" {
		log.Panic("port is empty")
	}

	ln, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		log.Panic(err)
	}

	m := cmux.New(ln)

	// We first match the connection against HTTP2 fields. If matched, the
	// connection will be sent through the "grpcL" listener.
	//grpcL := m.Match(cmux.HTTP2HeaderFieldPrefix("content-type", "application/grpc"))
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	//Otherwise, we match it against a websocket upgrade request.
	//wsl := m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket"))

	// Otherwise, we match it against HTTP1 methods. If matched,
	// it is sent through the "httpL" listener.
	httpL := m.Match(cmux.HTTP1Fast())
	// If not matched by HTTP, we assume it is an RPC connection.
	//rpcL := m.Match(cmux.Any())

	// Then we used the muxed listeners.
	go GRPCServer(grpcL, s.Token)
	go HTTPServer(httpL, s.Token)

	logs.Info("run muxed server on port: %s", s.Port)
	if err := m.Serve(); err != cmux.ErrListenerClosed {
		panic(err)
	}
}
