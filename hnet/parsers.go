package hnet

import (
	"fmt"
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

func handlePanic() {
	_ = recover()
}
