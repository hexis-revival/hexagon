package hnet

import (
	"fmt"

	"github.com/hexis-revival/hexagon/common"
)

func ensureAuthentication(handler func(*common.IOStream, *Player) error) func(*common.IOStream, *Player) error {
	return func(stream *common.IOStream, player *Player) error {
		if !player.IsAuthenticated() {
			player.CloseConnection()
			return fmt.Errorf("unauthenticated player")
		}

		return handler(stream, player)
	}
}

func ensureUnauthenticated(handler func(*common.IOStream, *Player) error) func(*common.IOStream, *Player) error {
	return func(stream *common.IOStream, player *Player) error {
		if player.IsAuthenticated() {
			player.CloseConnection()
			return fmt.Errorf("already authenticated")
		}

		return handler(stream, player)
	}
}
