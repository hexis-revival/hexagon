package hnet

import (
	"encoding/binary"
	"net"

	"github.com/lekuruu/hexagon/common"
)

type Player struct {
	Conn    net.Conn
	Name    string
	Version *VersionInfo
	Client  *ClientInfo
	Status  *Status
	Logger  *common.Logger
}

func (player *Player) Send(data []byte) error {
	_, err := player.Conn.Write(data)
	return err
}

func (player *Player) SendPacket(packetId uint32, data []byte) error {
	stream := common.NewIOStream([]byte{}, binary.BigEndian)
	stream.WriteU8(0x87)
	stream.WriteU32(packetId)
	stream.WriteU32(uint32(len(data)))
	stream.Write(data)
	player.Logger.Verbosef("<- %d '%s'", packetId, string(data))
	return player.Send(stream.Get())
}
