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

	// Set random player Id
	player.Presence.Id = uint32(rand.Intn(1000))
	player.Presence.Name = request.Username

	player.Version = request.Version
	player.Client = request.Client

	player.Logger.Infof(
		"Login attempt as '%s' with version %s",
		player.Presence.Name,
		player.Version.String(),
	)

	// Add to player collection
	player.Server.Players.Add(player)

	// Set placeholder stats
	player.Stats.UserId = player.Presence.Id
	player.Stats.Rank = 1
	player.Stats.Score = 300
	player.Stats.Unknown = 1
	player.Stats.Unknown2 = 2
	player.Stats.Accuracy = 0.9914
	player.Stats.Plays = 21

	for _, other := range player.Server.Players.All() {
		other.SendPacket(SERVER_USER_PRESENCE, player.Presence)
		other.SendPacket(SERVER_USER_STATS, player.Stats)

		player.SendPacket(SERVER_USER_PRESENCE, other.Presence)
		player.SendPacket(SERVER_USER_STATS, other.Stats)
	}

	response := LoginResponse{
		UserId:   player.Presence.Id,
		Username: player.Presence.Name,
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

	for _, userId := range userIds {
		stats := UserStats{
			UserId:   userId,
			Rank:     1,
			Score:    300,
			Unknown:  1,
			Unknown2: 2,
			Accuracy: 0.9914,
			Plays:    21,
		}

		player.SendPacket(SERVER_USER_STATS, stats)
	}

	return nil
}

func init() {
	Handlers[CLIENT_LOGIN] = handleLogin
	Handlers[CLIENT_CHANGE_STATUS] = handleStatusChange
	Handlers[CLIENT_REQUEST_STATS] = handleRequestStats
}
