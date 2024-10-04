package hnet

import (
	"encoding/binary"
	"fmt"

	"github.com/lekuruu/hexagon/common"
)

var Handlers = map[uint32]func(*common.IOStream, *Player) error{}

func handleLogin(stream *common.IOStream, player *Player) error {
	request := ReadLoginRequest(stream)

	if request == nil {
		return fmt.Errorf("failed to read login request")
	}

	player.Logger.Debugf("-> %s", request.String())
	player.Name = request.Username
	player.Version = request.Version
	player.Client = request.Client

	player.Logger.Infof(
		"Login attempt as '%s' with version %s",
		player.Name,
		player.Version.String(),
	)

	// TODO: Username & Password validation
	// TODO: Pull data from database

	responseStream := common.NewIOStream(
		[]byte{},
		binary.BigEndian,
	)

	response := LoginResponse{
		Username: player.Name,
		Unknown:  "",
		UserId:   1,
	}
	response.Serialize(responseStream)

	return player.SendPacket(SERVER_LOGIN_RESPONSE, responseStream.Get())
}

func handleStatusChange(stream *common.IOStream, player *Player) error {
	status := ReadStatusChange(stream)

	if status == nil {
		return fmt.Errorf("failed to read status change")
	}

	player.Status = status
	player.Logger.Debugf("-> %s", status.String())
	return nil
}

func init() {
	Handlers[CLIENT_LOGIN] = handleLogin
	Handlers[CLIENT_CHANGE_STATUS] = handleStatusChange
}
