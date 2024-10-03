package hnet

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lekuruu/hexagon/common"
)

type LoginRequest struct {
	Username string
	Password string
	Version  *VersionInfo
	Client   *ClientInfo
}

func (request *LoginRequest) String() string {
	return "LoginRequest{" +
		"Username: " + request.Username + ", " +
		"Password: " + request.Password + ", " +
		"Version: " + request.Version.String() + ", " +
		request.Client.String() + "}"
}

type ClientInfo struct {
	ExecutableHash string
	Adapters       []string
	Hash1          string // TODO
	Hash2          string // TODO
	Hash3          string // TODO
}

func (info *ClientInfo) String() string {
	return "ClientInfo{" +
		"ExecutableHash: " + info.ExecutableHash + ", " +
		"Adapters: " + strings.Join(info.Adapters, ":") + ", " +
		"Hash1: " + info.Hash1 + ", " +
		"Hash2: " + info.Hash2 + ", " +
		"Hash3: " + info.Hash3 + "}"
}

func (info *ClientInfo) IsWine() bool {
	return strings.HasPrefix(info.Hash3, "unk")
}

type VersionInfo struct {
	Major uint32
	Minor uint32
	Patch uint32
}

func (info *VersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d", info.Major, info.Minor, info.Patch)
}

func ReadLoginRequest(stream *common.IOStream) *LoginRequest {
	defer handlePanic()

	username := stream.ReadString()
	password := stream.ReadString()
	majorVersion := stream.ReadU32()
	minorVersion := stream.ReadU32()
	patchVersion := stream.ReadU32()
	clientInfo := stream.ReadString()

	version := &VersionInfo{
		Major: majorVersion,
		Minor: minorVersion,
		Patch: patchVersion,
	}

	return &LoginRequest{
		Username: username,
		Password: password,
		Version:  version,
		Client:   ParseClientInfo(clientInfo),
	}
}

func ParseClientInfo(clientInfoString string) *ClientInfo {
	parts := strings.Split(clientInfoString, ";")
	adapters := strings.Split(parts[1], ",")

	return &ClientInfo{
		Adapters:       adapters,
		ExecutableHash: parts[0],
		Hash1:          parts[2],
		Hash2:          parts[3],
		Hash3:          parts[4],
	}
}

type Status struct {
	Unknown uint32 // TODO
	Action  uint32
	Beatmap *BeatmapInfo
}

func (status *Status) HasBeatmapInfo() bool {
	return status.Action == ACTION_PLAYING ||
		status.Action == ACTION_EDITING ||
		status.Action == ACTION_TESTING
}

func (status Status) String() string {
	var beatmapString string = "nil"
	if status.Beatmap != nil {
		beatmapString = status.Beatmap.String()
	}

	return "Status{" +
		"Unknown: " + strconv.Itoa(int(status.Unknown)) + ", " +
		"Action: " + strconv.Itoa(int(status.Action)) + ", " +
		"Beatmap: " + beatmapString + "}"
}

type BeatmapInfo struct {
	Checksum string
	Id       uint32
	Artist   string
	Title    string
	Version  string
}

func (beatmap *BeatmapInfo) String() string {
	return "BeatmapInfo{" +
		"BeatmapMD5: " + beatmap.Checksum + ", " +
		"BeatmapID: " + strconv.Itoa(int(beatmap.Id)) + ", " +
		"BeatmapArtist: " + beatmap.Artist + ", " +
		"BeatmapTitle: " + beatmap.Title + ", " +
		"Version: " + beatmap.Version + "}"
}

func ReadStatusChange(stream *common.IOStream) *Status {
	defer handlePanic()

	status := &Status{
		Unknown: stream.ReadU32(),
		Action:  stream.ReadU32(),
		Beatmap: nil,
	}

	if !status.HasBeatmapInfo() {
		return status
	}

	status.Beatmap = &BeatmapInfo{
		Checksum: stream.ReadString(),
		Id:       stream.ReadU32(),
		Artist:   stream.ReadString(),
		Title:    stream.ReadString(),
		Version:  stream.ReadString(),
	}

	return status
}

func handlePanic() {
	_ = recover()
}
