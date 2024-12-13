package hnet

import (
	"strings"

	"github.com/hexis-revival/hexagon/common"
)

func ReadLoginRequest(stream *common.IOStream) *LoginRequest {
	defer recover()

	username := stream.ReadString()
	password := stream.ReadString()
	majorVersion := stream.ReadU32()
	minorVersion := stream.ReadU32()
	patchVersion := stream.ReadU32()
	clientInfo := stream.ReadString()
	_ = stream.ReadU8() // TODO

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

func ReadLoginRequestReconnect(stream *common.IOStream) *LoginRequest {
	defer recover()

	username := stream.ReadString()
	password := stream.ReadString()

	// Reconnect packet contains extra 4 bytes of null bytes
	_ = stream.ReadU32()

	majorVersion := stream.ReadU32()
	minorVersion := stream.ReadU32()
	patchVersion := stream.ReadU32()
	clientInfo := stream.ReadString()
	_ = stream.ReadU8() // TODO

	// At the end there are an extra 4 bytes
	// {0xFF, 0xFF, 0xFF, 0xFF}
	_ = stream.ReadU32()

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
	defer recover()

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

	status.Watching = stream.ReadString()

	status.Mods = &Mods{
		ArOffset:  stream.ReadI8(),
		OdOffset:  stream.ReadI8(),
		CsOffset:  stream.ReadI8(),
		HpOffset:  stream.ReadI8(),
		PlaySpeed: 1 + (0.5 * float32(stream.ReadI8()) / 10),
		Hidden:    stream.ReadBool(),
		NoFail:    stream.ReadBool(),
		Autoplay:  stream.ReadBool(),
	}

	return status
}

func ReadStatsRequest(stream *common.IOStream) *StatsRequest {
	defer recover()

	return &StatsRequest{
		UserIds: stream.ReadIntList(),
	}
}

func ReadRelationshipRequest(stream *common.IOStream) *RelationshipRequest {
	defer recover()

	status := common.StatusBlocked
	isFriend := stream.ReadBool()
	userId := stream.ReadU32()

	if isFriend {
		status = common.StatusFriend
	}

	return &RelationshipRequest{
		Status: status,
		UserId: userId,
	}
}

func ReadSpectateRequest(stream *common.IOStream) *SpectateRequest {
	defer recover()

	_ = stream.ReadBool() // TODO: this seems to be always 1?
	userId := stream.ReadU32()

	return &SpectateRequest{
		UserId: userId,
	}
}

func ReadHasMapRequest(stream *common.IOStream) *HasMapRequest {
	defer recover()

	_ = stream.ReadU8()
	hasMap := stream.ReadU32Bool()

	return &HasMapRequest{
		HasMap: hasMap,
	}
}

func ReadScorePack(stream *common.IOStream) *ScorePack {
	defer recover()

	action := stream.ReadU32()
	frames := make([]*common.ReplayFrame, stream.ReadU32())

	for i := range frames {
		frames[i] = common.ReadReplayFrame(stream)
	}

	return &ScorePack{
		Action: action,
		Frames: frames,
	}
}
