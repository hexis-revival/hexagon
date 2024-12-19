package hscore

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hexis-revival/hexagon/common"
)

func BeatmapUpdateHandler(ctx *Context) {
	request, err := NewBeatmapUpdateRequest(*ctx.Request)
	if err != nil {
		ctx.Server.Logger.Errorf("Failed to parse beatmap update request: %s", err)
		ctx.Response.WriteHeader(400)
	}

	user, err := common.FetchUserById(
		request.UserId,
		ctx.Server.State,
	)

	if err != nil {
		ctx.Server.Logger.Errorf("Failed to fetch user: %s", err)
		ctx.Response.WriteHeader(401)
	}

	ctx.Server.Logger.Infof(
		"Beatmap update request from '%s' (%d)",
		user.Name, request.BeatmapId,
	)

	beatmapData, err := ctx.Server.State.Storage.GetBeatmapFile(request.BeatmapId)
	if err != nil {
		ctx.Server.Logger.Errorf("Failed to fetch beatmap data: %s", err)
		ctx.Response.WriteHeader(404)
	}

	ctx.Response.WriteHeader(200)
	ctx.Response.Header().Set("Content-Type", "application/xml")
	ctx.Response.Write(beatmapData)
}

func NewBeatmapUpdateRequest(request http.Request) (*BeatmapUpdateRequest, error) {
	queryParameters := request.URL.Query()

	beatmapId := queryParameters.Get("b")
	if beatmapId == "" {
		return nil, fmt.Errorf("missing beatmap id")
	}

	setId := queryParameters.Get("s")
	if setId == "" {
		return nil, fmt.Errorf("missing set id")
	}

	userId := queryParameters.Get("u")
	if userId == "" {
		return nil, fmt.Errorf("missing user id")
	}

	beatmapIdInt, err := strconv.Atoi(beatmapId)
	if err != nil {
		return nil, fmt.Errorf("invalid beatmap id: %s", beatmapId)
	}

	setIdInt, err := strconv.Atoi(setId)
	if err != nil {
		return nil, fmt.Errorf("invalid set id: %s", setId)
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %s", userId)
	}

	return &BeatmapUpdateRequest{
		BeatmapId: beatmapIdInt,
		SetId:     setIdInt,
		UserId:    userIdInt,
	}, nil
}
