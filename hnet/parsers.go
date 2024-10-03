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
	BeatmapMD5    string
	BeatmapID     uint32
	BeatmapArtist string
	BeatmapTitle  string
	Version       string
}

func (beatmap BeatmapInfo) String() string {
	return "BeatmapInfo{" +
		"BeatmapMD5: " + beatmap.BeatmapMD5 + ", " +
		"BeatmapID: " + strconv.Itoa(int(beatmap.BeatmapID)) + ", " +
		"BeatmapArtist: " + beatmap.BeatmapArtist + ", " +
		"BeatmapTitle: " + beatmap.BeatmapTitle + ", " +
		"Version: " + beatmap.Version + "}"
}

func HasBeatmapInfo(action uint32) bool {
	return action == ACTION_PLAYING ||
		action == ACTION_EDITING ||
		action == ACTION_TESTING
}

func ReadStatusChange(stream *common.IOStream) *Status {
	defer handlePanic()

	unknown := stream.ReadU32() // TODO
	action := stream.ReadU32()

	var beatmap *BeatmapInfo = nil

	if HasBeatmapInfo(action) {
		md5 := stream.ReadString()
		id := stream.ReadU32()
		artist := stream.ReadString()
		title := stream.ReadString()
		version := stream.ReadString()

		beatmap = &BeatmapInfo{
			BeatmapMD5:    md5,
			BeatmapID:     id,
			BeatmapArtist: artist,
			BeatmapTitle:  title,
			Version:       version,
		}
	}

	status := &Status{
		Unknown: unknown,
		Action:  action,
		Beatmap: beatmap,
	}

	return status
}

func handlePanic() {
	_ = recover()
}
