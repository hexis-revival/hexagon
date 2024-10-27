package hnet

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/hexis-revival/hexagon/common"
)

type Player struct {
	Conn       net.Conn
	Logger     *common.Logger
	Server     *HNetServer
	Client     *ClientInfo
	Info       *UserInfo
	Stats      *UserStats
	Host       *Player
	Spectators *PlayerCollection
}

func (player *Player) Send(data []byte) error {
	player.Logger.Verbosef("<- %s", common.FormatBytes(data))
	_, err := player.Conn.Write(data)
	return err
}

func (player *Player) Receive() ([]byte, error) {
	buffer := make([]byte, 1024*1024)
	n, err := player.Conn.Read(buffer)

	if err != nil {
		return nil, err
	}

	buffer = buffer[:n]
	player.Logger.Verbosef("-> %s", common.FormatBytes(buffer))
	return buffer, nil
}

func (player *Player) OnConnect() {
	player.Logger.Debug("-> Connected")
}

func (player *Player) OnDisconnect() {
	player.Logger.Infof("Disconnected -> <%s>", player.Conn.RemoteAddr())
	player.Server.Players.Remove(player)
	player.Server.Players.Broadcast(SERVER_USER_QUIT, &QuitResponse{player.Info.Id})
	player.Conn.Close()
}

func (player *Player) CloseConnection() {
	player.RevokeLogin()
	player.OnDisconnect()
}

func (player *Player) LogIncomingPacket(packetId uint32, packet Serializable) {
	player.Logger.Debugf("-> %d: %s", packetId, packet.String())
}

func (player *Player) LogOutgoingPacket(packetId uint32, packet Serializable) {
	player.Logger.Debugf("<- %d: %s", packetId, packet.String())
}

func (player *Player) SendPacketData(packetId uint32, data []byte) error {
	stream := common.NewIOStream([]byte{}, binary.BigEndian)
	stream.WriteU8(0x87)
	stream.WriteU32(packetId)
	stream.WriteU32(uint32(len(data)))
	stream.Write(data)
	return player.Send(stream.Get())
}

func (player *Player) SendPacket(packetId uint32, packet Serializable) error {
	player.LogOutgoingPacket(packetId, packet)
	stream := common.NewIOStream([]byte{}, binary.BigEndian)
	packet.Serialize(stream)
	return player.SendPacketData(packetId, stream.Get())
}

func (player *Player) IsAuthenticated() bool {
	return player.Info.Id != 0 || player.Info.Name != ""
}

func (player *Player) OnLoginSuccess(responsePassword string, userObject *common.User) error {
	otherUser := player.Server.Players.ByID(uint32(userObject.Id))

	if otherUser != nil {
		// Another user with this account is online
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
		Password: responsePassword,
	}

	// Send login response
	err := player.SendPacket(SERVER_LOGIN_RESPONSE, response)
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

func (player *Player) OnLoginFailed(reason string) {
	player.Logger.Warningf("Login attempt failed: %s", reason)
	player.CloseConnection()
}

func (player *Player) RevokeLogin() error {
	return player.SendPacketData(SERVER_LOGIN_REVOKED, []byte{})
}

func (player *Player) AddRelationship(targetId uint32, status common.RelationshipStatus) error {
	rel := &common.Relationship{
		UserId:   int(player.Info.Id),
		TargetId: int(targetId),
		Status:   status,
	}

	return common.CreateUserRelationship(rel, player.Server.State)
}

func (player *Player) RemoveRelationship(targetId uint32, status common.RelationshipStatus) error {
	rel, err := common.FetchUserRelationship(
		int(player.Info.Id),
		int(targetId),
		player.Server.State,
	)

	if err != nil {
		return err
	}

	return common.RemoveUserRelationship(rel, player.Server.State)
}

func (player *Player) GetFriendIds() ([]uint32, error) {
	relationships, err := common.FetchUserRelationships(
		int(player.Info.Id),
		common.StatusFriend,
		player.Server.State,
	)

	if err != nil {
		return nil, err
	}

	friends := make([]uint32, 0, len(relationships))

	for _, rel := range relationships {
		friends = append(friends, uint32(rel.TargetId))
	}

	return friends, nil
}

func (player *Player) ApplyUserData(user *common.User) error {
	player.Info.Name = user.Name
	player.Info.Id = uint32(user.Id)
	player.Stats.UserId = uint32(user.Id)
	player.Stats.Rank = uint32(user.Stats.Rank)
	player.Stats.RankedScore = uint64(user.Stats.RankedScore)
	player.Stats.TotalScore = uint64(user.Stats.TotalScore)
	player.Stats.Plays = uint32(user.Stats.Playcount)
	player.Stats.Accuracy = user.Stats.Accuracy
	return nil
}

func (player *Player) StartSpectating(host *Player) error {
	host.Spectators.Add(player)
	player.Host = host

	response := &SpectateRequest{
		UserId: player.Info.Id,
	}

	err := host.SendPacket(SERVER_START_SPECTATING, response)
	if err != nil {
		return err
	}

	player.Logger.Infof("Started spectating '%s'", host.Info.Name)
	return nil
}

func (player *Player) IsSpectating() bool {
	return player.Host != nil
}

func (player *Player) HasSpectators() bool {
	return player.Spectators.Count() > 0
}
