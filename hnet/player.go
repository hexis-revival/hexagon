package hnet

import (
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
