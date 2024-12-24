package hnet

import (
	"encoding/hex"
	"fmt"

	"github.com/hexis-revival/hexagon/common"
)

var Handlers = map[uint32]func(*common.IOStream, *Player) error{}

func handleLogin(stream *common.IOStream, player *Player) error {
	request := ReadLoginRequest(stream)

	if request == nil {
		player.RevokeLogin()
		return fmt.Errorf("failed to read login request")
	}

	player.LogIncomingPacket(CLIENT_LOGIN, request)
	player.Client = request.Client

	if !player.Client.IsValid() {
		player.OnLoginFailed("Invalid client info")
		return nil
	}

	userObject, err := common.FetchUserByNameCaseInsensitive(
		request.Username,
		player.Server.State,
		"Stats",
	)

	if err != nil {
		player.OnLoginFailed("User not found")
		return nil
	}

	isCorrect := common.CheckPassword(
		request.Password,
		userObject.Password,
	)

	if !isCorrect {
		player.OnLoginFailed("Incorrect password")
		return nil
	}

	if !userObject.Activated {
		player.OnLoginFailed("Account not activated")
		return nil
	}

	if userObject.Restricted {
		player.OnLoginFailed("Account restricted")
		return nil
	}

	responsePasswordRaw := common.GetSHA512Hash(request.Password)
	responsePassword := hex.EncodeToString(responsePasswordRaw)

	return player.OnLoginSuccess(responsePassword, userObject)
}

func handleReconnect(stream *common.IOStream, player *Player) error {
	request := ReadLoginRequestReconnect(stream)

	if request == nil {
		player.RevokeLogin()
		return fmt.Errorf("failed to read login request")
	}

	player.LogIncomingPacket(CLIENT_LOGIN_RECONNECT, request)
	player.Client = request.Client

	if !player.Client.IsValid() {
		player.OnLoginFailed("Invalid client info")
		return nil
	}

	userObject, err := common.FetchUserByNameCaseInsensitive(
		request.Username,
		player.Server.State,
		"Stats",
	)

	if err != nil {
		player.OnLoginFailed("User not found")
		return nil
	}

	isCorrect := common.CheckPasswordHashedHex(
		request.Password,
		userObject.Password,
	)

	if !isCorrect {
		player.OnLoginFailed("Incorrect password")
		return nil
	}

	if !userObject.Activated {
		player.OnLoginFailed("Account not activated")
		return nil
	}

	if userObject.Restricted {
		player.OnLoginFailed("Account restricted")
		return nil
	}

	return player.OnLoginSuccess(request.Password, userObject)
}

func handleStatusChange(stream *common.IOStream, player *Player) error {
	status := ReadStatusChange(stream)

	if status == nil {
		return fmt.Errorf("failed to read status change")
	}

	if player.Stats.Status.Action == ACTION_PLAYING {
		time := player.Stats.Status.TimeSinceChanged()

		// Update user's playtime
		err := common.UpdatePlaytime(
			int(player.Info.Id),
			int(time.Seconds()),
			player.Server.State,
		)

		if err != nil {
			player.Logger.Errorf("Failed to update playtime: %s", err)
		}
	}

	player.LogIncomingPacket(CLIENT_CHANGE_STATUS, status)
	player.Stats.Status = status

	if player.HasSpectators() {
		player.Spectators.Broadcast(SERVER_SPECTATE_STATUS_UPDATE, player.Stats.Status)
	}

	return nil
}

func handleRequestStats(stream *common.IOStream, player *Player) error {
	statsRequest := ReadStatsRequest(stream)

	if statsRequest == nil {
		return fmt.Errorf("failed to read stats request")
	}

	player.LogIncomingPacket(CLIENT_REQUEST_STATS, statsRequest)

	for _, userId := range statsRequest.UserIds {
		user := player.Server.Players.ByID(userId)

		if user == nil {
			continue
		}

		player.SendPacket(SERVER_USER_STATS, user.Stats)
	}

	return nil
}

func handleStartSpectating(stream *common.IOStream, player *Player) error {
	request := ReadSpectateRequest(stream)

	if request == nil {
		return fmt.Errorf("failed to read spectate request")
	}

	target := player.Server.Players.ByID(request.UserId)

	if target == nil {
		return fmt.Errorf("user %d not found", request.UserId)
	}

	player.LogIncomingPacket(CLIENT_START_SPECTATING, request)
	return player.StartSpectating(target)
}

func handleStopSpectating(stream *common.IOStream, player *Player) error {
	if !player.IsSpectating() {
		return nil
	}

	request := ReadSpectateRequest(stream)

	if request == nil {
		return fmt.Errorf("failed to read spectate request")
	}

	player.LogIncomingPacket(CLIENT_STOP_SPECTATING, request)
	return player.StopSpectating()
}

func handleHasMap(stream *common.IOStream, player *Player) error {
	request := ReadHasMapRequest(stream)

	if request == nil {
		return fmt.Errorf("failed to read has map request")
	}

	player.LogIncomingPacket(CLIENT_SPECTATE_HAS_MAP, request)

	response := &HasMapResponse{
		UserId: player.Info.Id,
		HasMap: request.HasMap,
	}
	player.Host.SendPacket(SERVER_SPECTATE_HAS_MAP, response)

	return nil
}

func handleSpectateFrames(stream *common.IOStream, player *Player) error {
	if !player.HasSpectators() {
		return nil
	}

	scorePack := ReadScorePack(stream)

	if scorePack == nil {
		return fmt.Errorf("failed to read score pack")
	}

	player.LogIncomingPacket(CLIENT_SPECTATE_FRAMES, scorePack)
	player.Spectators.Broadcast(SERVER_SPECTATE_FRAMES, scorePack)
	return nil
}

func handleUserRelationshipAdd(stream *common.IOStream, player *Player) error {
	request := ReadRelationshipRequest(stream)

	if request == nil {
		return fmt.Errorf("failed to read relationship request")
	}

	target := player.Server.Players.ByID(request.UserId)

	if target == nil {
		return fmt.Errorf("user %d not found", request.UserId)
	}

	player.LogIncomingPacket(CLIENT_RELATIONSHIP_ADD, request)
	err := player.AddRelationship(request.UserId, request.Status)

	if err != nil {
		return err
	}

	player.Logger.Infof("Set relationship status to '%s' for %s", request.Status, target.Info.Name)
	return nil
}

func handleUserRelationshipRemove(stream *common.IOStream, player *Player) error {
	request := ReadRelationshipRequest(stream)

	if request == nil {
		return fmt.Errorf("failed to read relationship request")
	}

	target := player.Server.Players.ByID(request.UserId)

	if target == nil {
		return fmt.Errorf("user %d not found", request.UserId)
	}

	player.LogIncomingPacket(CLIENT_RELATIONSHIP_REMOVE, request)
	err := player.RemoveRelationship(request.UserId, request.Status)

	if err != nil {
		return err
	}

	player.Logger.Infof("Removed relationship status '%s' for %s", request.Status, target.Info.Name)
	return nil
}

func handleLeaderboardRequest(stream *common.IOStream, player *Player) error {
	request := ReadLeaderboardRequest(stream)

	if request == nil {
		return fmt.Errorf("failed to read leaderboard request")
	}

	player.LogIncomingPacket(CLIENT_LEADERBOARD_REQUEST, request)

	response := &LeaderboardResponse{
		BeatmapChecksum: request.BeatmapChecksum,
		ShowScores:      request.ShowScores,
		Unknown:         request.Unknown,
		Status:          common.BeatmapStatusNotSubmitted,
		NeedsUpdate:     false,
	}

	beatmap, err := common.FetchBeatmapById(
		int(request.BeatmapId),
		player.Server.State,
	)

	if err != nil {
		if err.Error() != "record not found" {
			return err
		}

		// Beatmap was not found, send empty response
		player.SendPacket(SERVER_LEADERBOARD_RESPONSE, response)
		return nil
	}

	response.NeedsUpdate = request.BeatmapChecksum != beatmap.Checksum
	response.Status = beatmap.Status
	return player.SendPacket(SERVER_LEADERBOARD_RESPONSE, response)
}

func handleStatsRefresh(stream *common.IOStream, player *Player) error {
	player.LogIncomingPacket(CLIENT_STATS_REFRESH, nil)

	if err := player.Refresh(); err != nil {
		return err
	}

	player.SendPacket(SERVER_USER_STATS, player.Stats)
	return nil
}

func init() {
	Handlers[CLIENT_LOGIN] = ensureUnauthenticated(handleLogin)
	Handlers[CLIENT_LOGIN_RECONNECT] = ensureUnauthenticated(handleReconnect)
	Handlers[CLIENT_CHANGE_STATUS] = ensureAuthentication(handleStatusChange)
	Handlers[CLIENT_REQUEST_STATS] = ensureAuthentication(handleRequestStats)
	Handlers[CLIENT_START_SPECTATING] = ensureAuthentication(handleStartSpectating)
	Handlers[CLIENT_STOP_SPECTATING] = ensureAuthentication(handleStopSpectating)
	Handlers[CLIENT_SPECTATE_HAS_MAP] = ensureAuthentication(handleHasMap)
	Handlers[CLIENT_SPECTATE_FRAMES] = ensureAuthentication(handleSpectateFrames)
	Handlers[CLIENT_RELATIONSHIP_ADD] = ensureAuthentication(handleUserRelationshipAdd)
	Handlers[CLIENT_RELATIONSHIP_REMOVE] = ensureAuthentication(handleUserRelationshipRemove)
	Handlers[CLIENT_LEADERBOARD_REQUEST] = ensureAuthentication(handleLeaderboardRequest)
	Handlers[CLIENT_STATS_REFRESH] = ensureAuthentication(handleStatsRefresh)
}
