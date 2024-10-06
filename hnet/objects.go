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
	return strings.HasPrefix(info.Hash3, "unk")
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
	return status.Action == ACTION_PLAYING ||
		status.Action == ACTION_EDITING ||
		status.Action == ACTION_TESTING
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
	UserId   uint32
	Rank     uint32
	Score    uint64
	Unknown  uint32
	Unknown2 uint32
	Accuracy float64
	Plays    uint32
}

func (stats UserStats) String() string {
	return fmt.Sprintf(
		"UserStats{UserId: %d, Rank: %d, Score: %d, Unknown: %d, Unknown2: %d, Accuracy: %f, Plays: %d}",
		stats.UserId,
		stats.Rank,
		stats.Score,
		stats.Unknown,
		stats.Unknown2,
		stats.Accuracy*100,
		stats.Plays,
	)
}

func NewUserStats() *UserStats {
	return &UserStats{
		Rank:     0,
		Score:    0,
		Unknown:  0,
		Unknown2: 0,
		Accuracy: 0.0,
		Plays:    0,
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
