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
		Conn:       conn,
		Logger:     logger,
		Server:     server,
		Info:       NewUserInfo(),
		Stats:      NewUserStats(),
		Spectators: NewPlayerCollection(),
	}

	player.OnConnect()
	defer server.CloseConnection(player)

	for {
		buffer, err := player.Receive()

		if err != nil {
			player.Logger.Debugf("Error receiving data: %s", err)
			break
		}

		if len(buffer) < HNET_PACKET_SIZE {
			player.Logger.Warningf("Invalid packet size: %d", len(buffer))
			break
		}

		magicByte := buffer[0]
		packetId := common.ReadU32BE(buffer[1:5])
		packetSize := common.ReadU32BE(buffer[5:9])

		if magicByte != 0x87 {
			player.Logger.Warningf("Invalid magic byte: %d", magicByte)
			break
		}

		if packetSize > uint32(len(buffer)-HNET_PACKET_SIZE) {
			player.Logger.Warningf("Invalid packet size: %d", packetSize)
			break
		}

		packetData := buffer[HNET_PACKET_SIZE : HNET_PACKET_SIZE+packetSize]
		handler, ok := Handlers[packetId]

		if !ok {
			player.Logger.Warningf("Unknown packetId: %d -> '%s'", packetId, common.FormatBytes(packetData))
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
