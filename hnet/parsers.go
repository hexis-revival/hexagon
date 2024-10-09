package hnet

import (
	"strings"

	"github.com/hexis-revival/hexagon/common"
)

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

	client := ParseClientInfo(clientInfo)
	client.Version = version

	return &LoginRequest{
		Username: username,
		Password: password,
		Client:   client,
	}
}

func ParseClientInfo(clientInfoString string) *ClientInfo {
	parts := strings.Split(clientInfoString, ";")
	adapters := strings.Split(parts[1], ",")

	return &ClientInfo{
		Adapters:       adapters,
		ExecutableHash: parts[0],
		AdaptersHash:   parts[2],
		UninstallId:    parts[3],
		DiskSignature:  parts[4],
	}
}

func ReadStatusChange(stream *common.IOStream) *Status {
	defer handlePanic()

	status := &Status{
		UserId:  stream.ReadU32(),
		Action:  stream.ReadU32(),
		Beatmap: nil,
	}

	if !status.HasBeatmapInfo() {
		return status
	}

	status.Beatmap = &BeatmapInfo{
		Checksum: stream.ReadString(),
		Id:       stream.ReadI32(),
		Artist:   stream.ReadString(),
		Title:    stream.ReadString(),
		Version:  stream.ReadString(),
	}

	return status
}

func ReadStatsRequest(stream *common.IOStream) *StatsRequest {
	defer handlePanic()

	return &StatsRequest{
		UserIds: stream.ReadIntList(),
	}
}

func handlePanic() {
	_ = recover()
}
