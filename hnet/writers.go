package hnet

import (
	"strings"

	"github.com/hexis-revival/hexagon/common"
)

func (request LoginRequest) Serialize(stream *common.IOStream) {
	stream.WriteString(request.Username)
	stream.WriteString(request.Password)
	request.Client.Serialize(stream)
}

func (response LoginResponse) Serialize(stream *common.IOStream) {
	stream.WriteString(response.Username)
	stream.WriteString(response.Password)
	stream.WriteU32(response.UserId)
}

func (version VersionInfo) Serialize(stream *common.IOStream) {
	stream.WriteU32(version.Major)
	stream.WriteU32(version.Minor)
	stream.WriteU32(version.Patch)
}

func (info ClientInfo) Serialize(stream *common.IOStream) {
	info.Version.Serialize(stream)
	parts := []string{
		info.ExecutableHash,
		strings.Join(info.Adapters, ","),
		info.AdaptersHash,
		info.UninstallId,
		info.DiskSignature,
	}
	stream.WriteString(strings.Join(parts, ";"))
}

func (info UserInfo) Serialize(stream *common.IOStream) {
	stream.WriteU32(info.Id)
	stream.WriteString(info.Name)
}

func (stats UserStats) Serialize(stream *common.IOStream) {
	stream.WriteU32(stats.UserId)
	stream.WriteU32(stats.Rank)
	stream.WriteU64(stats.RankedScore)
	stream.WriteU64(stats.TotalScore)
	stream.WriteF64(stats.Accuracy)
	stream.WriteU32(stats.Plays)
	stats.Status.Serialize(stream)
}

func (request StatsRequest) Serialize(stream *common.IOStream) {
	stream.WriteIntList(request.UserIds)
}

func (status Status) Serialize(stream *common.IOStream) {
	stream.WriteU32(status.UserId)
	stream.WriteU32(status.Action)

	if !status.HasBeatmapInfo() {
		return
	}

	status.Beatmap.Serialize(stream)
}

func (info BeatmapInfo) Serialize(stream *common.IOStream) {
	stream.WriteString(info.Checksum)
	stream.WriteI32(info.Id)
	stream.WriteString(info.Artist)
	stream.WriteString(info.Title)
	stream.WriteString(info.Version)
}

func (friends FriendsList) Serialize(stream *common.IOStream) {
	stream.WriteIntList(friends.FriendIds)
}

func (response QuitResponse) Serialize(stream *common.IOStream) {
	stream.WriteU8(0)
	stream.WriteU32(response.UserId)
}

func (request RelationshipRequest) Serialize(stream *common.IOStream) {
	stream.WriteBool(request.Status == common.StatusFriend)
	stream.WriteU32(request.UserId)
}

func (request SpectateRequest) Serialize(stream *common.IOStream) {
	stream.WriteU8(1)
	stream.WriteU32(request.UserId)
}
