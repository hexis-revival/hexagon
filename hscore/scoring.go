package hscore

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"net/http"

	"github.com/hexis-revival/hexagon/common"
)

const (
	ValidationError     = "Invalid score submission request."
	AuthenticationError = "Could not submit your score. Please check your login credentials!"
	BeatmapError        = "This beatmap is not available for score submission."
	ServerError         = "Could not submit your score due to an internal server error."
)

func ResolveBeatmap(score *ScoreData, server *ScoreServer) (*common.Beatmap, error) {
	beatmap, err := common.FetchBeatmapByChecksum(
		score.BeatmapChecksum,
		server.State,
		"Beatmapset",
	)

	if err != nil {
		return nil, err
	}

	return beatmap, nil
}

func ValidateScore(user *common.User, request *ScoreSubmissionRequest) error {
	return nil // TODO
}

func InsertScore(user *common.User, beatmap *common.Beatmap, score *ScoreData, server *ScoreServer) (*common.Score, error) {
	return nil, nil // TODO
}

func UploadReplay(scoreId int, replay *common.ReplayData, storage *common.Storage) error {
	return nil // TODO
}

func UpdateUserStatistics(user *common.User, server *ScoreServer) error {
	return nil // TODO
}

func UpdateBeatmapStatistics(beatmap *common.Beatmap, server *ScoreServer) error {
	return nil // TODO
}

func WriteError(statusCode int, errorMessage string, ctx *Context) error {
	ctx.Response.WriteHeader(statusCode)
	encoder := json.NewEncoder(ctx.Response)
	response := ScoreSubmissionResponse{Success: false, Error: errorMessage}
	return encoder.Encode(response)
}

func ScoreSubmissionHandler(ctx *Context) {
	// Parse score submission request
	request, err := NewScoreSubmissionRequest(ctx.Request)

	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Header().Set("Content-Type", "application/json")

	if err != nil {
		ctx.Server.Logger.Errorf("Error parsing score submission request: %v", err)
		WriteError(http.StatusBadRequest, ValidationError, ctx)
		return
	}

	ctx.Server.Logger.Debugf(
		"Score submission request: %s",
		request.String(),
	)

	user, success := AuthenticateUser(
		request.ScoreData.Username,
		request.Password,
		ctx.Server,
	)

	if !success {
		ctx.Server.Logger.Warningf("Failed to authenticate user '%s'", request.ScoreData.Username)
		WriteError(http.StatusUnauthorized, AuthenticationError, ctx)
		return
	}

	beatmap, err := ResolveBeatmap(
		request.ScoreData,
		ctx.Server,
	)

	if err != nil {
		ctx.Server.Logger.Warningf("Error resolving beatmap: %v", err)
		WriteError(http.StatusNotFound, BeatmapError, ctx)
		return
	}

	if err = ValidateScore(user, request); err != nil {
		ctx.Server.Logger.Warningf("Error validating score: %v", err)
		WriteError(http.StatusBadRequest, ValidationError, ctx)
		return
	}

	score, err := InsertScore(
		user, beatmap,
		request.ScoreData,
		ctx.Server,
	)

	if err != nil {
		ctx.Server.Logger.Warningf("Error inserting score: %v", err)
		WriteError(http.StatusInternalServerError, ServerError, ctx)
		return
	}

	if err = UploadReplay(score.Id, request.Replay, &ctx.Server.State.Storage); err != nil {
		ctx.Server.Logger.Warningf("Error uploading replay: %v", err)
		WriteError(http.StatusInternalServerError, ServerError, ctx)
		return
	}

	if err = UpdateUserStatistics(user, ctx.Server); err != nil {
		ctx.Server.Logger.Warningf("Error updating user statistics: %v", err)
		WriteError(http.StatusInternalServerError, ServerError, ctx)
		return
	}

	if err = UpdateBeatmapStatistics(beatmap, ctx.Server); err != nil {
		ctx.Server.Logger.Warningf("Error updating beatmap statistics: %v", err)
		WriteError(http.StatusInternalServerError, ServerError, ctx)
		return
	}

	response := ScoreSubmissionResponse{Success: true}
	json.NewEncoder(ctx.Response).Encode(response)
}

func NewScoreSubmissionRequest(request *http.Request) (*ScoreSubmissionRequest, error) {
	err := request.ParseMultipartForm(10 << 20) // ~10 MB
	if err != nil {
		return nil, err
	}

	replay := GetMultipartFormFile(request, "r")
	password := GetMultipartFormValue(request, "p")
	ivBase64 := GetMultipartFormValue(request, "iv")
	processListBase64 := GetMultipartFormValue(request, "pl")
	scoreDataBase64 := GetMultipartFormValue(request, "s")
	clientDataBase64 := GetMultipartFormValue(request, "i")

	collection := common.NewErrorCollection()

	iv, err := base64.StdEncoding.DecodeString(ivBase64)
	collection.Add(err)

	processList, err := base64.StdEncoding.DecodeString(processListBase64)
	collection.Add(err)

	scoreData, err := base64.StdEncoding.DecodeString(scoreDataBase64)
	collection.Add(err)

	clientData, err := base64.StdEncoding.DecodeString(clientDataBase64)
	collection.Add(err)

	processListDecrypted, err := common.DecryptScoreData(iv, processList)
	collection.Add(err)

	scoreDataDecrypted, err := common.DecryptScoreData(iv, scoreData)
	collection.Add(err)

	clientDataDecryptedBase64, err := common.DecryptScoreData(iv, clientData)
	collection.Add(err)

	clientDataDecrypted, err := base64.StdEncoding.DecodeString(string(clientDataDecryptedBase64))
	collection.Add(err)

	scoreDataStruct, err := ParseScoreData(scoreDataDecrypted)
	collection.Add(err)

	if collection.HasErrors() {
		return nil, collection.Pop(0)
	}

	replayStream := common.NewIOStream(replay, binary.BigEndian)
	replayData, _ := common.ReadCompressedReplay(replayStream)
	processListData := ParseProcessList(processListDecrypted)

	return &ScoreSubmissionRequest{
		Replay:      replayData,
		Password:    password,
		ProcessList: processListData,
		ScoreData:   scoreDataStruct,
		ClientData:  string(clientDataDecrypted),
	}, nil
}
