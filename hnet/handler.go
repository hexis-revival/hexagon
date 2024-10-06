package hnet

import (
	"fmt"
	"math/rand"

	"github.com/lekuruu/hexagon/common"
)

var Handlers = map[uint32]func(*common.IOStream, *Player) error{}

func handleLogin(stream *common.IOStream, player *Player) error {
	request := ReadLoginRequest(stream)

	if request == nil {
		return fmt.Errorf("failed to read login request")
	}

	player.LogIncomingPacket(CLIENT_LOGIN, request)
	player.Name = request.Username
	player.Version = request.Version
	player.Client = request.Client

	player.Logger.Infof(
		"Login attempt as '%s' with version %s",
		player.Name,
		player.Version.String(),
	)

	// Set random player Id for now
	player.Id = uint32(rand.Intn(1000))
	player.Server.Players.Add(player)

	// TODO: Username & Password validation
	// TODO: Pull data from database

	response := LoginResponse{
		UserId:   player.Id,
		Username: player.Name,
		Password: request.Password,
	}

	return player.SendPacket(SERVER_LOGIN_RESPONSE, response)
}

func handleStatusChange(stream *common.IOStream, player *Player) error {
	status := ReadStatusChange(stream)

	if status == nil {
		return fmt.Errorf("failed to read status change")
	}

	player.Status = status
	player.LogIncomingPacket(CLIENT_CHANGE_STATUS, status)
	return nil
}

func handleRequestStats(stream *common.IOStream, player *Player) error {
	var userIds = stream.ReadIntList()

	player.Logger.Infof("Requested stats of %d users", len(userIds))
	// TODO: pull stats, and enqueue them

	return nil
}

func init() {
	Handlers[CLIENT_LOGIN] = handleLogin
	Handlers[CLIENT_CHANGE_STATUS] = handleStatusChange
	Handlers[CLIENT_REQUEST_STATS] = handleRequestStats
}
