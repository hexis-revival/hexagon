package hnet

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/hexis-revival/hexagon/common"
)

type Serializable interface {
	Serialize(stream *common.IOStream)
	String() string
}

type LoginRequest struct {
	Username string
	Password string
	Client   *ClientInfo
}

func (request LoginRequest) String() string {
	return fmt.Sprintf(
		"LoginRequest{Username: %s, Password: %s, %s}",
		request.Username,
		request.Password,
		request.Client.String(),
	)
}

type ClientInfo struct {
	Version        *VersionInfo
	ExecutableHash string
	Adapters       []string
	AdaptersHash   string
	UninstallId    string
	DiskSignature  string
}

func (info ClientInfo) String() string {
	return fmt.Sprintf(
		"Client{Version: %s, ExecutableHash: %s, Adapters: %v, AdaptersHash: %s, UninstallId: %s, DiskSignature: %s}",
		info.Version.String(),
		info.ExecutableHash,
		info.Adapters,
		info.AdaptersHash,
		info.UninstallId,
		info.DiskSignature,
	)
}

func (info ClientInfo) IsWine() bool {
	return info.DiskSignature == "unknown"
}

func (info ClientInfo) IsValid() bool {
	if len(info.Adapters) == 0 {
		return false
	}

	adaptersString := strings.Join(info.Adapters, ",")
	adaptersHash := md5.Sum([]byte(adaptersString))
	adaptersHashHex := hex.EncodeToString(adaptersHash[:])

	if adaptersHashHex != info.AdaptersHash {
		return false
	}

	if len(info.ExecutableHash) != 32 {
		return false
	}

	if len(info.AdaptersHash) != 32 {
		return false
	}

	if len(info.UninstallId) != 32 {
		return false
	}

	if len(info.DiskSignature) != 32 && !info.IsWine() {
		return false
	}

	return true
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
	// TODO: Add remaining packet data
}

func (info UserInfo) String() string {
	return fmt.Sprintf(
		"UserInfo{UserId: %d. Username: %s}",
		info.Id,
		info.Name,
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

type QuitResponse struct {
	UserId uint32
}

func (response QuitResponse) String() string {
	return fmt.Sprintf(
		"QuitResponse{UserId: %d}",
		response.UserId,
	)
}
