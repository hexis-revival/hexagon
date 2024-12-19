package hscore

import (
	"fmt"
	"net/http"
	"strconv"
)

func BeatmapUpdateHandler(ctx *Context) {}

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
