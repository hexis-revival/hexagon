package hnet

import (
	"encoding/binary"
	"fmt"
	"net"
	"runtime/debug"

	"github.com/lekuruu/go-raknet"
	"github.com/lekuruu/hexagon/common"
)

const HNET_PACKET_SIZE = 9

type HNetServer struct {
	Players  PlayerCollection
	listener *raknet.Listener
	logger   *common.Logger
	host     string
	port     int
}

func NewServer(host string, port int, logger *common.Logger) *HNetServer {
	return &HNetServer{
		Players: NewPlayerCollection(),
		logger:  logger,
		host:    host,
		port:    port,
	}
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
	logger := common.CreateLogger(
		conn.RemoteAddr().String(),
		server.logger.GetLevel(),
	)

	player := &Player{
		Conn:   conn,
		Logger: logger,
		Server: server,
	}

	player.OnConnect()
	defer server.CloseConnection(player)

	for {
		buffer, err := player.Receive(1024 * 1024)

		if err != nil {
			server.logger.Debugf("Failed to read data: '%s'", err)
			return
		}

		if len(buffer) < HNET_PACKET_SIZE {
			server.logger.Errorf("Invalid packet size: %d", len(buffer))
			return
		}

		magicByte := buffer[0]
		packetId := common.ReadU32BE(buffer[1:5])
		packetSize := common.ReadU32BE(buffer[5:9])

		if magicByte != 0x87 {
			server.logger.Errorf("Invalid magic byte: %d", magicByte)
			return
		}

		packetData := buffer[HNET_PACKET_SIZE : HNET_PACKET_SIZE+packetSize]
		handler, ok := Handlers[packetId]

		if !ok {
			server.logger.Warningf("Unknown packetId: %d -> '%s'", packetId, string(packetData))
			continue
		}

		stream := common.NewIOStream(packetData, binary.BigEndian)
		err = handler(stream, player)

		if err != nil {
			server.logger.Errorf("Error handling packet: %s", err)
			continue
		}
	}
}

func (server *HNetServer) CloseConnection(player *Player) {
	if r := recover(); r != nil {
		server.logger.Errorf("Panic: '%s'", r)
		server.logger.Debug(string(debug.Stack()))
	}

	player.OnDisconnect()
}
