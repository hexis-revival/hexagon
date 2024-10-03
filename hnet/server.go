package hnet

import (
	"fmt"
	"net"
	"runtime/debug"

	"github.com/lekuruu/go-raknet"
	"github.com/lekuruu/hexagon/common"
)

type HNetServer struct {
	listener *raknet.Listener
	logger   *common.Logger
	host     string
	port     int
}

func NewServer(host string, port int, logger *common.Logger) (*HNetServer, error) {
	return &HNetServer{
		logger: logger,
		host:   host,
		port:   port,
	}, nil
}

func (server *HNetServer) Serve() {
	// Set hexis protocol version
	raknet.SetProtocolVersion(6)

	bind := fmt.Sprintf("%s:%d", server.host, server.port)
	listener, err := raknet.Listen(bind)

	if err != nil {
		server.logger.Errorf("Failed to listen on %s: '%s'", bind, err)
		return
	}

	defer listener.Close()

	server.logger.Infof("Listening on %s", listener.Addr())
	server.listener = listener

	for {
		conn, _ := listener.Accept()
		go server.HandleConnection(conn)
	}
}

func (server *HNetServer) HandleConnection(conn net.Conn) {
	defer server.CloseConnection(conn)

	buffer := make([]byte, 1024*1024)
	n, err := conn.Read(buffer)

	if err != nil {
		server.logger.Error(err)
		return
	}

	buffer = buffer[:n]
	server.logger.Info("Received %s", string(buffer))
	// TODO: Add Handler
}

func (server *HNetServer) CloseConnection(conn net.Conn) {
	if r := recover(); r != nil {
		server.logger.Error("Panic: %s", r)
		debug.PrintStack()
	}
}
