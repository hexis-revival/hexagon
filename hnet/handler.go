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

	player.LogIncomingPacket(CLIENT_CHANGE_STATUS, status)
	player.Stats.Status = status
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

func init() {
	Handlers[CLIENT_LOGIN] = ensureUnauthenticated(handleLogin)
	Handlers[CLIENT_LOGIN_RECONNECT] = ensureUnauthenticated(handleReconnect)
	Handlers[CLIENT_CHANGE_STATUS] = ensureAuthentication(handleStatusChange)
	Handlers[CLIENT_REQUEST_STATS] = ensureAuthentication(handleRequestStats)
	Handlers[CLIENT_START_SPECTATING] = ensureAuthentication(handleStartSpectating)
	Handlers[CLIENT_RELATIONSHIP_ADD] = ensureAuthentication(handleUserRelationshipAdd)
	Handlers[CLIENT_RELATIONSHIP_REMOVE] = ensureAuthentication(handleUserRelationshipRemove)
}
