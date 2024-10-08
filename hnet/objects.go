package hnet

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lekuruu/hexagon/common"
)

type Serializable interface {
	Serialize(stream *common.IOStream)
	String() string
}

type LoginRequest struct {
	Username string
	Password string
	Version  *VersionInfo
	Client   *ClientInfo
}

func (request LoginRequest) String() string {
	return fmt.Sprintf(
		"LoginRequest{Username: %s, Password: %s, Version: %s, %s}",
		request.Username,
		request.Password,
		request.Version.String(),
		request.Client.String(),
	)
}

type ClientInfo struct {
	ExecutableHash string
	Adapters       []string
	Hash1          string // TODO
	Hash2          string // TODO
	Hash3          string // TODO
}

func (info ClientInfo) String() string {
	return fmt.Sprintf(
		"Client{ExecutableHash: %s, Adapters: %v, Hash1: %s, Hash2: %s, Hash3: %s}",
		info.ExecutableHash,
		info.Adapters,
		info.Hash1,
		info.Hash2,
		info.Hash3,
	)
}

func (info ClientInfo) IsWine() bool {
	return strings.HasPrefix(info.Hash3, "unknown")
}

type VersionInfo struct {
	Major uint32
	Minor uint32
	Patch uint32
}

func (info VersionInfo) String() string {
	return fmt.Sprintf(
		"%d.%d.%d",
		info.Major,
		info.Minor,
		info.Patch,
	)
}

type Status struct {
	UserId  uint32
	Action  uint32
	Beatmap *BeatmapInfo
}

func (status Status) HasBeatmapInfo() bool {
	return status.Action > ACTION_AWAY
}

func (status Status) String() string {
	var beatmapString string = "nil"
	if status.Beatmap != nil {
		beatmapString = status.Beatmap.String()
	}

	return fmt.Sprintf(
		"Status{UserId: %s, Action: %s, Beatmap: %s}",
		strconv.Itoa(int(status.UserId)),
		strconv.Itoa(int(status.Action)),
		beatmapString,
	)
}

func NewStatus() *Status {
	return &Status{
		UserId:  0,
		Action:  1,
		Beatmap: nil,
	}
}

type BeatmapInfo struct {
	Checksum string
	Id       int32
	Artist   string
	Title    string
	Version  string
}

func (beatmap BeatmapInfo) String() string {
	return fmt.Sprintf(
		"BeatmapInfo{Checksum: %s, Id: %d, Artist: %s, Title: %s, Version: %s}",
		beatmap.Checksum,
		beatmap.Id,
		beatmap.Artist,
		beatmap.Title,
		beatmap.Version,
	)
}

type LoginResponse struct {
	Username string
	Password string
	UserId   uint32
}

func (response LoginResponse) String() string {
	return fmt.Sprintf(
		"LoginResponse{Username: %s, Password: %s, UserId: %d}",
		response.Username,
		response.Password,
		response.UserId,
	)
}

type UserInfo struct {
	Id   uint32
	Name string
	// TODO: Add remaining presence data
}

func (presence UserInfo) String() string {
	return fmt.Sprintf(
		"UserPresence{UserId: %d. Username: %s}",
		presence.Id,
		presence.Name,
	)
}

func NewUserInfo() *UserInfo {
	return &UserInfo{
		Id:   0,
		Name: "",
	}
}

type UserStats struct {
	UserId      uint32
	Rank        uint32
	RankedScore uint64
	TotalScore  uint64
	Accuracy    float64
	Plays       uint32
	Status      *Status
}

func (stats UserStats) String() string {
	return fmt.Sprintf(
		"UserStats{UserId: %d, Rank: %d, RankedScore: %d, TotalScore: %d, Accuracy: %f, Plays: %d, %s}",
		stats.UserId,
		stats.Rank,
		stats.RankedScore,
		stats.TotalScore,
		stats.Accuracy*100,
		stats.Plays,
		stats.Status.String(),
	)
}

func NewUserStats() *UserStats {
	return &UserStats{
		Rank:        0,
		RankedScore: 0,
		TotalScore:  0,
		Accuracy:    0.0,
		Plays:       0,
		Status:      NewStatus(),
	}
}

type StatsRequest struct {
	UserIds []uint32
}

func (request StatsRequest) String() string {
	return fmt.Sprintf(
		"StatsRequest{UserIds: %v}",
		request.UserIds,
	)
}

func NewStatsRequest() *StatsRequest {
	return &StatsRequest{
		UserIds: make([]uint32, 0),
	}
}

type FriendsList struct {
	FriendIds []uint32
}

func (friends FriendsList) String() string {
	return fmt.Sprintf(
		"FriendsList{Friends: %v}",
		friends.FriendIds,
	)
}

func NewFriendsList() *FriendsList {
	return &FriendsList{
		FriendIds: []uint32{},
	}
}
