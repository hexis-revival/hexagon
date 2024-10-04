package hnet

import (
	"fmt"
	"strconv"
	"strings"
)

type LoginRequest struct {
	Username string
	Password string
	Version  *VersionInfo
	Client   *ClientInfo
}

func (request *LoginRequest) String() string {
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

func (info *ClientInfo) String() string {
	return fmt.Sprintf(
		"Client{ExecutableHash: %s, Adapters: %v, Hash1: %s, Hash2: %s, Hash3: %s}",
		info.ExecutableHash,
		info.Adapters,
		info.Hash1,
		info.Hash2,
		info.Hash3,
	)
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

func (status *Status) HasBeatmapInfo() bool {
	return status.Action == ACTION_PLAYING ||
		status.Action == ACTION_EDITING ||
		status.Action == ACTION_TESTING
}

func (status *Status) String() string {
	var beatmapString string = "nil"
	if status.Beatmap != nil {
		beatmapString = status.Beatmap.String()
	}

	return fmt.Sprintf(
		"Status{Unknown: %s, Action: %s, Beatmap: %s}",
		strconv.Itoa(int(status.UserId)),
		strconv.Itoa(int(status.Action)),
		beatmapString,
	)
}

type BeatmapInfo struct {
	Checksum string
	Id       uint32
	Artist   string
	Title    string
	Version  string
}

func (beatmap *BeatmapInfo) String() string {
	return fmt.Sprintf(
		"BeatmapInfo{Checksum: %s, Id: %d, Artist: %s, Title: %s, Version: %s}",
		beatmap.Checksum,
		beatmap.Id,
		beatmap.Artist,
		beatmap.Title,
		beatmap.Version,
	)
}
