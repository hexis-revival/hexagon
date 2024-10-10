package hnet

import (
	"fmt"

	"github.com/hexis-revival/hexagon/common"
	"golang.org/x/crypto/bcrypt"
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
		player.CloseConnection()
		player.Logger.Warning("Login attempt failed: Invalid client info")
		return nil
	}

	userObject, err := common.FetchUserByNameCaseInsensitive(
		request.Username,
		player.Server.State,
	)

	if err != nil {
		player.CloseConnection()
		player.Logger.Warning("Login attempt failed: User not found")
		return nil
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(userObject.Password),
		[]byte(request.Password),
	)

	if err != nil {
		player.CloseConnection()
		player.Logger.Warning("Login attempt failed: Incorrect password")
		return nil
	}

	if !userObject.Activated {
		player.CloseConnection()
		player.Logger.Warning("Login attempt failed: Account not activated")
		return nil
	}

	if userObject.Restricted {
		player.CloseConnection()
		player.Logger.Warning("Login attempt failed: Account restricted")
		return nil
	}

	otherUser := player.Server.Players.ByID(uint32(userObject.Id))

	if otherUser != nil {
		otherUser.CloseConnection()
	}

	// Ensure that the stats object exists
	userObject.EnsureStats(player.Server.State)

	// Populate player info & stats
	player.ApplyUserData(userObject)
	player.Server.Players.Add(player)

	player.Logger.Infof(
		"Login attempt as '%s' with version %s",
		player.Info.Name,
		player.Client.Version.String(),
	)

	player.Logger.SetName(fmt.Sprintf(
		"Player \"%s\"",
		player.Info.Name,
	))

	for _, other := range player.Server.Players.All() {
		other.SendPacket(SERVER_USER_INFO, player.Info)
		player.SendPacket(SERVER_USER_INFO, other.Info)
	}

	response := LoginResponse{
		UserId:   player.Info.Id,
		Username: player.Info.Name,
		Password: request.Password,
	}

	// Send login response
	err = player.SendPacket(SERVER_LOGIN_RESPONSE, response)
	if err != nil {
		player.CloseConnection()
		return err
	}

	friendIds, err := player.GetFriendIds()
	if err != nil {
		return err
	}

	// Send friends list
	friendsList := FriendsList{FriendIds: friendIds}
	err = player.SendPacket(SERVER_FRIENDS_LIST, friendsList)
	if err != nil {
		return err
	}

	return nil
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
	Handlers[CLIENT_LOGIN] = handleLogin
	Handlers[CLIENT_CHANGE_STATUS] = handleStatusChange
	Handlers[CLIENT_REQUEST_STATS] = handleRequestStats
	Handlers[CLIENT_RELATIONSHIP_ADD] = handleUserRelationshipAdd
	Handlers[CLIENT_RELATIONSHIP_REMOVE] = handleUserRelationshipRemove
}
