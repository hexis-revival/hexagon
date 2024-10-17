package hscore

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/hexis-revival/hexagon/common"
)

func BeatmapGenIdHandler(ctx *Context) {
	request, err := NewBeatmapSubmissionRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Warningf("Beatmap submission request error: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("Beatmap submission request: %s", request)
	ctx.Response.WriteHeader(http.StatusOK)
	// TODO: Implement logic
}

func NewBeatmapSubmissionRequest(request *http.Request) (*BeatmapSubmissionRequest, error) {
	err := request.ParseMultipartForm(10 << 20) // ~10 MB
	if err != nil {
		return nil, err
	}

	username := GetMultipartFormValue(request, "u")
	password := GetMultipartFormValue(request, "p")
	setId := GetMultipartFormValue(request, "s")
	beatmapIdsString := GetMultipartFormValue(request, "v")
	clientVersion := GetMultipartFormValue(request, "x")

	collection := common.NewErrorCollection()

	setIdInt, err := strconv.Atoi(setId)
	collection.Add(err)

	clientVersionInt, err := strconv.Atoi(clientVersion)
	collection.Add(err)

	beatmapIdsList := strings.Split(beatmapIdsString, ":")
	beatmapIds := make([]int, 0, len(beatmapIdsList))

	for _, beatmapId := range beatmapIdsList {
		beatmapIdInt, err := strconv.Atoi(beatmapId)
		beatmapIds = append(beatmapIds, beatmapIdInt)
		collection.Add(err)
	}

	if collection.HasErrors() {
		return nil, collection.Pop(0)
	}

	return &BeatmapSubmissionRequest{
		Username:      username,
		Password:      password,
		BeatmapIds:    beatmapIds,
		SetId:         setIdInt,
		ClientVersion: clientVersionInt,
	}, nil
}
