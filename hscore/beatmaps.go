package hscore

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hexis-revival/hexagon/common"
)

const (
	BssSuccess = 0
	// The username or password you entered is incorrect.
	BssAuthenticationError = 1
	// The beatmap you are trying to submit is no longer available.
	BssNotAvailable = 2
	// The beatmap you are trying to submit has already been ranked.
	BssAlreadyRanked = 3
	// The beatmap you are trying to submit is not owned by you.
	BssInvalidOwner = 4
)

func FormatBeatmapsetName(beatmapset *common.Beatmapset) string {
	return fmt.Sprintf(
		"%d %s - %s (%s)",
		beatmapset.Id,
		beatmapset.Artist,
		beatmapset.Title,
		beatmapset.Creator.Name,
	)
}

func AuthenticateUser(username string, password string, server *ScoreServer) (*common.User, bool) {
	userObject, err := common.FetchUserByNameCaseInsensitive(
		username,
		server.State,
		"Stats",
	)

	if err != nil {
		server.Logger.Warningf("[Beatmap Submission] User '%s' not found", username)
		return nil, false
	}

	decodedPassword, err := hex.DecodeString(password)

	if err != nil {
		server.Logger.Warningf("[Beatmap Submission] Password decoding error: %s", err)
		return nil, false
	}

	isCorrect := common.CheckPasswordHashed(
		decodedPassword,
		userObject.Password,
	)

	if !isCorrect {
		server.Logger.Warningf("[Beatmap Submission] Incorrect password for '%s'", username)
		return nil, false
	}

	if !userObject.Activated {
		server.Logger.Warningf("[Beatmap Submission] Account not activated for '%s'", username)
		return nil, false
	}

	if userObject.Restricted {
		server.Logger.Warningf("[Beatmap Submission] Account restricted for '%s'", username)
		return nil, false
	}

	return userObject, true
}

func CreateBeatmapset(beatmapIds []int, user *common.User, server *ScoreServer) (*common.Beatmapset, error) {
	beatmapset := &common.Beatmapset{
		Title:     "",
		Artist:    "",
		Source:    "",
		CreatorId: user.Id,
	}

	err := common.CreateBeatmapset(beatmapset, server.State)
	if err != nil {
		return nil, err
	}

	beatmaps := make([]common.Beatmap, 0, len(beatmapIds))

	for i := 0; i < len(beatmapIds); i++ {
		beatmap := common.Beatmap{SetId: beatmapset.Id, CreatorId: user.Id}
		beatmaps = append(beatmaps, beatmap)
	}

	err = common.CreateBeatmaps(beatmaps, server.State)
	if err != nil {
		return nil, err
	}

	beatmapset.Beatmaps = beatmaps
	return beatmapset, nil
}

func BeatmapGenIdHandler(ctx *Context) {
	request, err := NewBeatmapSubmissionRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Warningf("[Beatmap Submission] Request error: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("[Beatmap Submission] Request: %s", request)
	ctx.Response.WriteHeader(http.StatusOK)

	user, success := AuthenticateUser(
		request.Username,
		request.Password,
		ctx.Server,
	)

	response := &BeatmapSubmissionResponse{
		StatusCode: BssSuccess,
		SetId:      -1,
		BeatmapIds: []int{},
	}

	if !success {
		response.StatusCode = BssAuthenticationError
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	beatmapset, err := common.FetchBeatmapsetById(
		request.SetId,
		ctx.Server.State,
		"Beatmaps",
	)

	if err != nil {
		beatmapset, err = CreateBeatmapset(
			request.BeatmapIds,
			user,
			ctx.Server,
		)

		if err != nil {
			ctx.Server.Logger.Warningf("[Beatmap Submission] Beatmapset creation error: %s", err)
			response.StatusCode = BssNotAvailable
			ctx.Response.Write([]byte(response.Write()))
			return
		}

		for _, beatmap := range beatmapset.Beatmaps {
			response.BeatmapIds = append(response.BeatmapIds, beatmap.Id)
		}

		response.SetId = beatmapset.Id
		response.StatusCode = BssSuccess
		ctx.Response.Write([]byte(response.Write()))
		ctx.Server.Logger.Infof(
			"[Beatmap Submission] Beatmapset for '%s' created (%d)",
			user.Name,
			beatmapset.Id,
		)
		return
	}

	// TODO: Generate/Update beatmapset
	ctx.Server.Logger.Infof(
		"[Beatmap Submission] Beatmapset for '%s' updated (%d)",
		user.Name,
		beatmapset.Id,
	)
	ctx.Response.Write([]byte(response.Write()))
}

func BeatmapUploadHandler(ctx *Context) {
	request, err := NewBeatmapUploadRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Warningf("[Beatmap Submission] Upload request error: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("[Beatmap Submission] Upload request: %s", request)
	ctx.Response.WriteHeader(http.StatusOK)

	// TODO: Implement beatmap upload logic
	response := &BeatmapUploadResponse{Success: true}
	ctx.Response.Write([]byte(response.Write()))
}

func BeatmapGenTopicHandler(ctx *Context) {
	request, err := NewBeatmapDescriptionRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Warningf("[Beatmap Submission] Description request error: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("[Beatmap Submission] Description request: %s", request)
	ctx.Response.WriteHeader(http.StatusOK)

	response := &BeatmapDescriptionResponse{
		TopicId: 1,
		Content: "The quick brown fox jumps over the lazy dog.",
	}

	// TODO: Implement beatmap description logic
	ctx.Response.Write([]byte(response.Write()))
}

func BeatmapPostHandler(ctx *Context) {
	request, err := NewBeatmapPostRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Warningf("[Beatmap Submission] Post request error: %s", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("[Beatmap Submission] Post request: %s", request)
	ctx.Response.WriteHeader(http.StatusOK)
	// TODO: Implement beatmap post logic
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

func NewBeatmapPostRequest(request *http.Request) (*BeatmapPostRequest, error) {
	err := request.ParseMultipartForm(10 << 20) // ~10 MB
	if err != nil {
		return nil, err
	}

	username := GetMultipartFormValue(request, "u")
	password := GetMultipartFormValue(request, "p")
	clientVersion := GetMultipartFormValue(request, "x")
	setId := GetMultipartFormValue(request, "s")
	content := GetMultipartFormValue(request, "t")

	setIdInt, err := strconv.Atoi(setId)
	if err != nil {
		return nil, err
	}

	clientVersionInt, err := strconv.Atoi(clientVersion)
	if err != nil {
		return nil, err
	}

	return &BeatmapPostRequest{
		Username:      username,
		Password:      password,
		ClientVersion: clientVersionInt,
		SetId:         setIdInt,
		Content:       content,
	}, nil
}
