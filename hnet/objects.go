package hnet

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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
	return common.FormatStruct(request)
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
	return common.FormatStruct(info)
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
	UserId   uint32
	Action   uint32
	Beatmap  *BeatmapInfo
	Watching string
	Mods     *Mods
}

func (status Status) HasBeatmapInfo() bool {
	return status.Action > ACTION_AWAY
}

func (status Status) String() string {
	return common.FormatStruct(status)
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
	return common.FormatStruct(beatmap)
}

type LoginResponse struct {
	Username string
	Password string
	UserId   uint32
}

func (response LoginResponse) String() string {
	return common.FormatStruct(response)
}

type UserInfo struct {
	Id   uint32
	Name string
	// TODO: Add remaining packet data
}

func (info UserInfo) String() string {
	return common.FormatStruct(info)
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
	return common.FormatStruct(stats)
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
	return common.FormatStruct(request)
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
	return common.FormatStruct(friends)
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
	return common.FormatStruct(response)
}

type RelationshipRequest struct {
	Status common.RelationshipStatus
	UserId uint32
}

func (request RelationshipRequest) String() string {
	return common.FormatStruct(request)
}

type SpectateRequest struct {
	UserId uint32
}

func (request SpectateRequest) String() string {
	return common.FormatStruct(request)
}

type HasMapRequest struct {
	HasMap bool
}

func (request HasMapRequest) String() string {
	return common.FormatStruct(request)
}

type HasMapResponse struct {
	UserId uint32
	HasMap bool
}

func (response HasMapResponse) String() string {
	return common.FormatStruct(response)
}

type ScorePack struct {
	Action uint32
	Frames []*common.ReplayFrame
}

func (pack ScorePack) String() string {
	return common.FormatStruct(pack)
}

type Mods struct {
	ArOffset  int8
	OdOffset  int8
	CsOffset  int8
	HpOffset  int8
	PlaySpeed float32
	Hidden    bool
	NoFail    bool
	Autoplay  bool
}

func (mods *Mods) String() string {
	return common.FormatStruct(mods)
}
