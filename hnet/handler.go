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
	player.Version = request.Version
	player.Client = request.Client

	// Set random player Id
	player.Info.Id = uint32(rand.Intn(1000))
	player.Info.Name = request.Username

	player.Logger.Infof(
		"Login attempt as '%s' with version %s",
		player.Info.Name,
		player.Version.String(),
	)

	// Add to player collection
	player.Server.Players.Add(player)

	// Set placeholder stats
	player.Stats.UserId = player.Info.Id
	player.Stats.Rank = 1
	player.Stats.Score = 300
	player.Stats.Unknown = 1
	player.Stats.Unknown2 = 2
	player.Stats.Accuracy = 0.9914
	player.Stats.Plays = 21

	for _, other := range player.Server.Players.All() {
		other.SendPacket(SERVER_USER_INFO, player.Info)
		player.SendPacket(SERVER_USER_INFO, other.Info)
	}

	response := LoginResponse{
		UserId:   player.Info.Id,
		Username: player.Info.Name,
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
	statsRequest := ReadStatsRequest(stream)

	if statsRequest == nil {
		return fmt.Errorf("failed to read stats request")
	}

	player.Logger.Debugf("-> %s", statsRequest.String())

	for _, userId := range statsRequest.UserIds {
		user := player.Server.Players.ByID(userId)

		if user == nil {
			continue
		}

		player.SendPacket(SERVER_USER_STATS, user.Stats)
	}

	return nil
}

func init() {
	Handlers[CLIENT_LOGIN] = handleLogin
	Handlers[CLIENT_CHANGE_STATUS] = handleStatusChange
	Handlers[CLIENT_REQUEST_STATS] = handleRequestStats
}
