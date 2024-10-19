package hscore

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/hexis-revival/hexagon/common"
)

const (
	BssSuccess       = 0
	BssInvalidOwner  = 1
	BssNotAvailable  = 2
	BssAlreadyRanked = 3
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

	// NOTE: The client will "update" the beatmap if
	//       the same setId is responded with.
	//       Otherwise it will do a full submission.
	response := &BeatmapSubmissionResponse{
		StatusCode: BssSuccess,
		SetId:      request.SetId,
		BeatmapIds: request.BeatmapIds,
	}

	// TODO: Implement beatmap submission logic
	ctx.Response.Write([]byte(response.Write()))
}

func BeatmapUploadHandler(ctx *Context) {
	request, err := NewBeatmapUploadRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Warningf("Beatmap upload request error: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("Beatmap upload request: %s", request)
	ctx.Response.WriteHeader(http.StatusOK)

	// TODO: Implement beatmap upload logic
	response := &BeatmapUploadResponse{Success: true}
	ctx.Response.Write([]byte(response.Write()))
}

func BeatmapGenTopicHandler(ctx *Context) {
	request, err := NewBeatmapDescriptionRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Warningf("Beatmap description request error: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("Beatmap description request: %s", request)
	ctx.Response.WriteHeader(http.StatusOK)
	// TODO: Response
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

func NewBeatmapUploadRequest(request *http.Request) (*BeatmapUploadRequest, error) {
	err := request.ParseMultipartForm(10 << 20) // ~10 MB
	if err != nil {
		return nil, err
	}

	username := GetMultipartFormValue(request, "u")
	password := GetMultipartFormValue(request, "p")
	clientVersion := GetMultipartFormValue(request, "x")

	clientVersionInt, err := strconv.Atoi(clientVersion)
	if err != nil {
		return nil, err
	}

	zip, err := GetMultipartZipFile(request, "d")
	if err != nil {
		return nil, err
	}

	return &BeatmapUploadRequest{
		Username:      username,
		Password:      password,
		ClientVersion: clientVersionInt,
		Package:       zip,
	}, nil
}

func NewBeatmapDescriptionRequest(request *http.Request) (*BeatmapDescriptionRequest, error) {
	err := request.ParseMultipartForm(10 << 20) // ~10 MB
	if err != nil {
		return nil, err
	}

	username := GetMultipartFormValue(request, "u")
	password := GetMultipartFormValue(request, "p")
	clientVersion := GetMultipartFormValue(request, "x")
	setId := GetMultipartFormValue(request, "s")

	setIdInt, err := strconv.Atoi(setId)
	if err != nil {
		return nil, err
	}

	clientVersionInt, err := strconv.Atoi(clientVersion)
	if err != nil {
		return nil, err
	}

	return &BeatmapDescriptionRequest{
		Username:      username,
		Password:      password,
		ClientVersion: clientVersionInt,
		SetId:         setIdInt,
	}, nil
}
