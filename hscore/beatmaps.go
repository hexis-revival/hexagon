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

func ValidateBeatmapset(beatmapset *common.Beatmapset, user *common.User, server *ScoreServer) int {
	if beatmapset.CreatorId != user.Id {
		server.Logger.Warningf("[Beatmap Submission] Invalid owner for '%s'", user.Name)
		return BssInvalidOwner
	}

	if beatmapset.Status > common.StatusPending {
		server.Logger.Warningf("[Beatmap Submission] Beatmapset already ranked (%d)", beatmapset.Id)
		return BssAlreadyRanked
	}

	if beatmapset.AvailabilityStatus != common.BeatmapHasDownload {
		server.Logger.Warningf("[Beatmap Submission] Beatmapset not available (%d)", beatmapset.Id)
		return BssNotAvailable
	}

	return BssSuccess
}

func UpdateBeatmapIds(beatmapset *common.Beatmapset, beatmapIds []int, server *ScoreServer) (_ []int, err error) {
	if len(beatmapIds) <= 0 {
		return nil, fmt.Errorf("beatmapIds list is empty")
	}

	currentBeatmapIds := make(map[int]bool, len(beatmapset.Beatmaps))
	beatmapIdMap := make(map[int]bool, len(beatmapIds))

	for _, beatmap := range beatmapset.Beatmaps {
		currentBeatmapIds[beatmap.Id] = true
	}

	for _, beatmapId := range beatmapIds {
		beatmapIdMap[beatmapId] = true
	}

	if len(beatmapIds) < len(beatmapset.Beatmaps) {
		// Ensure every beatmap id is inside current beatmap ids
		for _, beatmapId := range beatmapIds {
			if !currentBeatmapIds[beatmapId] {
				return nil, fmt.Errorf("Beatmap '%d' not found", beatmapId)
			}
		}

		// Remove unused beatmaps
		for _, beatmap := range beatmapset.Beatmaps {
			if !beatmapIdMap[beatmap.Id] {
				continue
			}

			err := common.RemoveBeatmap(&beatmap, server.State)
			if err != nil {
				return nil, err
			}
		}

		server.Logger.Debugf(
			"Removed %d beatmaps from beatmapset",
			len(currentBeatmapIds)-len(beatmapIds),
		)
	}

	// Calculate how many beatmaps we need to create
	requiredMaps := max(0, len(beatmapIds)-len(currentBeatmapIds))
	newBeatmaps := make([]common.Beatmap, 0, requiredMaps)

	// Create new beatmaps
	for i := 0; i < requiredMaps; i++ {
		beatmap := common.Beatmap{SetId: beatmapset.Id, CreatorId: beatmapset.CreatorId}
		newBeatmaps = append(newBeatmaps, beatmap)

		err := common.CreateBeatmap(&beatmap, server.State)
		if err != nil {
			return nil, err
		}
	}

	server.Logger.Debugf(
		"Added %d beatmaps to beatmapset",
		requiredMaps,
	)

	// Append new beatmaps to the beatmapset & return new beatmap ids
	beatmapset.Beatmaps, err = common.FetchBeatmapsBySetId(
		beatmapset.Id,
		server.State,
	)

	if err != nil {
		return nil, err
	}

	beatmapIds = make([]int, 0, len(beatmapset.Beatmaps))

	for _, beatmap := range beatmapset.Beatmaps {
		beatmapIds = append(beatmapIds, beatmap.Id)
	}

	return beatmapIds, nil
}

func RemoveInactiveBeatmaps(user *common.User, server *ScoreServer) error {
	beatmapsets, err := common.FetchBeatmapsetsByStatus(
		user.Id,
		common.StatusNotSubmitted,
		server.State,
	)

	if err != nil {
		return err
	}

	server.Logger.Debugf(
		"Found %d inactive beatmapsets for '%s'",
		len(beatmapsets), user.Name,
	)
	// TODO: Remove beatmap assets from storage

	for _, beatmapset := range beatmapsets {
		err := common.RemoveBeatmapsBySetId(beatmapset.Id, server.State)
		if err != nil {
			return err
		}

		err = common.RemoveBeatmapset(&beatmapset, server.State)
		if err != nil {
			return err
		}
	}
	return nil
}

func RemainingBeatmapUpdloads(user *common.User, server *ScoreServer) (int, error) {
	unrankedBeatmaps, err := common.FetchBeatmapsetUnrankedCountByCreatorId(
		user.Id,
		server.State,
	)

	if err != nil {
		return 0, err
	}

	rankedBeatmaps, err := common.FetchBeatmapsetRankedCountByCreatorId(
		user.Id,
		server.State,
	)

	if err != nil {
		return 0, err
	}

	// Users can upload up to 6 pending maps plus
	// 1 per ranked map, up to a maximum of 10.
	return (6 - unrankedBeatmaps) + min(rankedBeatmaps, 10), nil
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

	err = RemoveInactiveBeatmaps(user, ctx.Server)
	if err != nil {
		ctx.Server.Logger.Warningf("[Beatmap Submission] Failed to remove inactive beatmaps: %s", err)
	}

	beatmapset, err := common.FetchBeatmapsetById(
		request.SetId,
		ctx.Server.State,
		"Beatmaps",
	)

	if err != nil {
		if err.Error() != "record not found" {
			// Database failure
			ctx.Server.Logger.Errorf("[Beatmap Submission] Beatmapset fetch error: %s", err)
			response.StatusCode = BssNotAvailable
			ctx.Response.Write([]byte(response.Write()))
			return
		}

		remainingUploads, err := RemainingBeatmapUpdloads(user, ctx.Server)
		if err != nil {
			ctx.Server.Logger.Warningf("[Beatmap Submission] Failed to fetch remaining uploads: %s", err)
			response.StatusCode = BssNotAvailable
			ctx.Response.Write([]byte(response.Write()))
			return
		}

		if remainingUploads <= 0 {
			ctx.Server.Logger.Warningf("[Beatmap Submission] No remaining uploads for '%s'", user.Name)
			response.StatusCode = BssNotAvailable
			ctx.Response.Write([]byte(response.Write()))
			return
		}

		beatmapset, err = CreateBeatmapset(
			request.BeatmapIds,
			user,
			ctx.Server,
		)

		if err != nil {
			ctx.Server.Logger.Errorf("[Beatmap Submission] Beatmapset creation error: %s", err)
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
			"[Beatmap Submission] Beatmapset for '%s' created with %d beatmaps (%d)",
			user.Name,
			len(beatmapset.Beatmaps),
			beatmapset.Id,
		)
		return
	}

	response.StatusCode = ValidateBeatmapset(
		beatmapset,
		user,
		ctx.Server,
	)

	if response.StatusCode != BssSuccess {
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	response.BeatmapIds, err = UpdateBeatmapIds(
		beatmapset,
		request.BeatmapIds,
		ctx.Server,
	)

	if err != nil {
		ctx.Server.Logger.Errorf("[Beatmap Submission] Beatmapset update error: %s", err)
		response.StatusCode = BssNotAvailable
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	response.SetId = beatmapset.Id
	ctx.Server.Logger.Infof(
		"[Beatmap Submission] Beatmapset for '%s' updated with %d beatmaps (%d)",
		user.Name,
		len(beatmapset.Beatmaps),
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
	response := &BeatmapUploadResponse{Success: true}

	user, success := AuthenticateUser(
		request.Username,
		request.Password,
		ctx.Server,
	)

	if !success {
		response.Success = false
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	beatmapset, err := common.FetchBeatmapsetById(
		request.SetId,
		ctx.Server.State,
		"Beatmaps",
	)

	if err != nil {
		ctx.Server.Logger.Warningf("[Beatmap Submission] Beatmapset fetch error: %s", err)
		response.Success = false
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	statusCode := ValidateBeatmapset(
		beatmapset,
		user,
		ctx.Server,
	)

	if statusCode != BssSuccess {
		response.Success = false
		ctx.Response.Write([]byte(response.Write()))
		return
	}

	// TODO: ...
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
