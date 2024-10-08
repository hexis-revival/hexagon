package hnet

import (
	"encoding/binary"
	"fmt"
	"net"
	"runtime/debug"

	"github.com/hexis-revival/go-raknet"
	"github.com/hexis-revival/hexagon/common"
)

const HNET_PACKET_SIZE = 9

type HNetServer struct {
	Players  PlayerCollection
	State    *common.State
	Listener *raknet.Listener
	Logger   *common.Logger
	Host     string
	Port     int
}

func NewServer(host string, port int, logger *common.Logger, state *common.State) *HNetServer {
	return &HNetServer{
		Players: NewPlayerCollection(),
		Logger:  logger,
		State:   state,
		Host:    host,
		Port:    port,
	}
}

func (server *HNetServer) Serve() {
	// Set hexis protocol version
	raknet.SetProtocolVersion(6)

	bind := fmt.Sprintf("%s:%d", server.Host, server.Port)
	listener, err := raknet.Listen(bind)

	if err != nil {
		server.Logger.Errorf("Failed to listen on %s: '%s'", bind, err)
		return
	}

	defer listener.Close()

	server.Logger.Infof("Listening on %s", listener.Addr())
	server.Listener = listener

	for {
		conn, _ := listener.Accept()
		go server.HandleConnection(conn)
	}
}

func (server *HNetServer) HandleConnection(conn net.Conn) {
	logger := common.CreateLogger(
		conn.RemoteAddr().String(),
		server.Logger.GetLevel(),
	)

	player := &Player{
		Conn:   conn,
		Logger: logger,
		Server: server,
		Info:   NewUserInfo(),
		Stats:  NewUserStats(),
	}

	player.OnConnect()
	defer server.CloseConnection(player)

	for {
		buffer, err := player.Receive(1024 * 1024)

		if err != nil {
			player.Logger.Debugf("Failed to read data: '%s'", err)
			return
		}

		if len(buffer) < HNET_PACKET_SIZE {
			player.Logger.Errorf("Invalid packet size: %d", len(buffer))
			return
		}

		magicByte := buffer[0]
		packetId := common.ReadU32BE(buffer[1:5])
		packetSize := common.ReadU32BE(buffer[5:9])

		if magicByte != 0x87 {
			player.Logger.Errorf("Invalid magic byte: %d", magicByte)
			return
		}

		packetData := buffer[HNET_PACKET_SIZE : HNET_PACKET_SIZE+packetSize]
		handler, ok := Handlers[packetId]

		if !ok {
			player.Logger.Warningf("Unknown packetId: %d -> '%s'", packetId, string(packetData))
			continue
		}

		stream := common.NewIOStream(packetData, binary.BigEndian)
		err = handler(stream, player)

		if err != nil {
			player.Logger.Errorf("Error handling packet '%d': %s", packetId, err)
			continue
		}
	}
}

func (server *HNetServer) CloseConnection(player *Player) {
	if r := recover(); r != nil {
		server.Logger.Errorf("Panic: '%s'", r)
		server.Logger.Debug(string(debug.Stack()))
	}

	player.OnDisconnect()
}
