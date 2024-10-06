package hnet

import (
	"strings"

	"github.com/lekuruu/hexagon/common"
)

func (request LoginRequest) Serialize(stream *common.IOStream) {
	stream.WriteString(request.Username)
	stream.WriteString(request.Password)
	request.Version.Serialize(stream)
	request.Client.Serialize(stream)
}

func (version VersionInfo) Serialize(stream *common.IOStream) {
	stream.WriteU32(version.Major)
	stream.WriteU32(version.Minor)
	stream.WriteU32(version.Patch)
}

func (info ClientInfo) Serialize(stream *common.IOStream) {
	parts := []string{
		info.ExecutableHash,
		strings.Join(info.Adapters, ","),
		info.Hash1,
		info.Hash2,
		info.Hash3,
	}
	stream.WriteString(strings.Join(parts, ";"))
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

func (response LoginResponse) Serialize(stream *common.IOStream) {
	stream.WriteString(response.Username)
	stream.WriteString(response.Password)
	stream.WriteU32(response.UserId)
}

func (presence UserInfo) Serialize(stream *common.IOStream) {
	stream.WriteU32(presence.Id)
	stream.WriteString(presence.Name)
}

func (stats UserStats) Serialize(stream *common.IOStream) {
	stream.WriteU32(stats.UserId)
	stream.WriteU32(stats.Rank)
	stream.WriteU64(stats.Score)
	stream.WriteU32(stats.Unknown)
	stream.WriteU32(stats.Unknown2)
	stream.WriteF64(stats.Accuracy)
	stream.WriteU32(stats.Plays)
}
