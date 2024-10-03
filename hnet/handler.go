package hnet

import (
	"fmt"

	"github.com/lekuruu/hexagon/common"
)

var Handlers = map[uint32]func(*common.IOStream, *Player) error{}

func handleLogin(stream *common.IOStream, player *Player) error {
	request := ReadLoginRequest(stream)

	if request == nil {
		return fmt.Errorf("failed to read login request")
	}

	player.Logger.Debug("-> %s", request.String())
	player.Name = request.Username
	player.Version = request.Version
	player.Client = request.Client

	// TODO: Username & Password validation
	return nil
}

func init() {
	Handlers[CLIENT_LOGIN] = handleLogin
}
