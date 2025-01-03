package hscore

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"

	"github.com/hexis-revival/hexagon/common"
)

var processRegexString = `^(.*?)(?: \((.*?)\))?;$`
var processRegex = regexp.MustCompile(processRegexString)

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
		"Set",
	)

	if err != nil {
		return nil, err
	}

	return beatmap, nil
}

func ValidateScore(user *common.User, beatmap *common.Beatmap, request *ScoreSubmissionRequest) (bool, error) {
	if request.ScoreData.Mods.Auto {
		return true, errors.New("submitted score with auto mod")
	}

	if request.ScoreData.MaxCombo <= 0 {
		return true, fmt.Errorf("submitted score with invalid max combo '%d'", request.ScoreData.MaxCombo)
	}

	if request.ScoreData.PassedObjects() <= 0 {
		return true, errors.New("submitted score with no passed objects")
	}

	if request.ScoreData.PassedObjects() > beatmap.TotalObjects() {
		return true, fmt.Errorf("submitted score with too many passed objects '%d'", request.ScoreData.PassedObjects())
	}

	if request.ScoreData.TotalScore <= 0 {
		return true, fmt.Errorf("submitted score with invalid total score '%d'", request.ScoreData.TotalScore)
	}

	if request.ScoreData.Passed && len(request.Replay.Frames) <= 100 {
		return true, fmt.Errorf("submitted score with too few replay frames '%d'", len(request.Replay.Frames))
	}

	if request.ScoreData.MaxCombo > beatmap.MaxCombo {
		comboDifference := request.ScoreData.MaxCombo - beatmap.MaxCombo

		if comboDifference > 5 {
			return false, fmt.Errorf("submitted score with invalid max combo '%d'", request.ScoreData.MaxCombo)
		}
	}

	if request.ScoreData.ClientBuildDate != 20140304 {
		return false, fmt.Errorf("submitted score with unknown client build date '%d'", request.ScoreData.ClientBuildDate)
	}

	if request.ScoreData.ClientVersion != 105 {
		return false, fmt.Errorf("submitted score with unknown client version '%d'", request.ScoreData.ClientVersion)
	}

	if !request.ScoreData.Passed {
		return false, nil
	}

	for _, process := range request.ProcessList {
		match := processRegex.FindStringSubmatch(process)

		if match == nil {
			return false, fmt.Errorf("submitted score with invalid process list '%s'", request.ProcessList)
		}

		if match[1] != "Hexis.exe" && match[2] != "Hexis" {
			continue
		}

		return false, nil
	}

	return true, errors.New("could not find hexis inside process list")
}

func InsertScore(user *common.User, beatmap *common.Beatmap, scoreData *ScoreData, server *ScoreServer) (*common.Score, error) {
	score := &common.Score{
		BeatmapId:     beatmap.Id,
		UserId:        user.Id,
		Checksum:      scoreData.ScoreChecksum,
		Status:        common.ScoreStatusUnranked,
		ClientVersion: scoreData.ClientVersion,
		TotalScore:    int64(scoreData.TotalScore),
		MaxCombo:      scoreData.MaxCombo,
		Accuracy:      scoreData.Accuracy(),
		Grade:         scoreData.Grade(),
		FullCombo:     scoreData.Perfect,
		Passed:        scoreData.Passed,
		Count300:      scoreData.Count300,
		Count100:      scoreData.Count100,
		Count50:       scoreData.Count50,
		CountGeki:     scoreData.CountGeki,
		CountKatu:     scoreData.CountKatu,
		CountGood:     scoreData.CountGood,
		CountMiss:     scoreData.CountMiss,
		AROffset:      scoreData.Mods.ArOffset,
		ODOffset:      scoreData.Mods.OdOffset,
		CSOffset:      scoreData.Mods.CsOffset,
		HPOffset:      scoreData.Mods.HpOffset,
		PSOffset:      scoreData.Mods.PsOffset,
		ModHidden:     scoreData.Mods.Hidden,
		ModNoFail:     scoreData.Mods.NoFail,
	}

	if beatmap.Status < common.BeatmapStatusRanked {
		// Beatmap is unranked, insert unranked score
		return score, common.CreateScore(score, server.State)
	}

	if !score.Passed {
		// User did not pass this map, insert failed score
		score.Status = common.ScoreStatusFailed
		return score, common.CreateScore(score, server.State)
	}

	personalBest, err := common.FetchPersonalBest(
		user.Id,
		beatmap.Id,
		server.State,
	)

	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	if personalBest == nil {
		// No personal best, insert score as PB
		score.Status = common.ScoreStatusPB
		return score, common.CreateScore(score, server.State)
	}

	if personalBest.TotalScore > score.TotalScore {
		// New score is lower than PB, insert as submitted
		score.Status = common.ScoreStatusSubmitted
		return score, common.CreateScore(score, server.State)
	}

	// A new personal best has been achieved
	score.Status = common.ScoreStatusPB
	personalBest.Status = common.ScoreStatusSubmitted

	if err = common.UpdateScore(personalBest, server.State); err != nil {
		return nil, err
	}

	return score, common.CreateScore(score, server.State)
}

func UploadReplay(scoreId int, replay *common.ReplayData, storage common.Storage) error {
	stream := common.NewIOStream([]byte{}, binary.BigEndian)
	replay.SerializeFrames(stream)
	return storage.SaveReplayFile(scoreId, stream.Get())
}

func UpdateUserStatistics(scoreData *ScoreData, user *common.User, server *ScoreServer) (err error) {
	user.Stats.TotalScore += int64(scoreData.TotalScore)
	user.Stats.TotalHits += int64(scoreData.TotalHits())
	user.Stats.Playcount += 1

	bestScoresRanked, err := common.FetchBestScores(
		user.Id,
		int(common.BeatmapStatusRanked),
		server.State,
	)

	if err != nil {
		return err
	}

	bestScoresApproved, err := common.FetchBestScores(
		user.Id,
		int(common.BeatmapStatusApproved),
		server.State,
	)

	if err != nil {
		return err
	}

	totalScores := len(bestScoresRanked) + len(bestScoresApproved)

	if totalScores <= 0 {
		// User has not set any scores yet
		return nil
	}

	user.Stats.RankedScore = 0

	for _, score := range bestScoresRanked {
		user.Stats.RankedScore += score.TotalScore
	}

	gradeMap := map[common.Grade]int{
		common.GradeD:  0,
		common.GradeC:  0,
		common.GradeB:  0,
		common.GradeA:  0,
		common.GradeS:  0,
		common.GradeSH: 0,
		common.GradeX:  0,
		common.GradeXH: 0,
	}

	accuracySum := 0.0
	maxCombo := 0

	for _, score := range bestScoresRanked {
		accuracySum += score.Accuracy
		gradeMap[score.Grade]++

		if score.MaxCombo > maxCombo {
			maxCombo = score.MaxCombo
		}
	}

	for _, score := range bestScoresApproved {
		accuracySum += score.Accuracy
		gradeMap[score.Grade]++

		if score.MaxCombo > maxCombo {
			maxCombo = score.MaxCombo
		}
	}

	user.Stats.MaxCombo = maxCombo
	user.Stats.Accuracy = math.Round(accuracySum/float64(totalScores)*1e6) / 1e6
	user.Stats.XHCount = gradeMap[common.GradeXH]
	user.Stats.XCount = gradeMap[common.GradeX]
	user.Stats.SHCount = gradeMap[common.GradeSH]
	user.Stats.SCount = gradeMap[common.GradeS]
	user.Stats.ACount = gradeMap[common.GradeA]
	user.Stats.BCount = gradeMap[common.GradeB]
	user.Stats.CCount = gradeMap[common.GradeC]
	user.Stats.DCount = gradeMap[common.GradeD]

	err = common.UpdateRankingsEntry(&user.Stats, user.Country, server.State)
	if err != nil {
		server.Logger.Errorf("Failed to update rankings entry: %v", err)
	}

	user.Stats.Rank, err = common.GetScoreRank(user.Id, server.State)
	if err != nil {
		server.Logger.Errorf("Failed to get user rank: %v", err)
	}

	return common.UpdateStats(&user.Stats, server.State)
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
		ctx.Server.Logger.Anomalyf("Error parsing score submission request: %v", err)
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

	if reject, err := ValidateScore(user, beatmap, request); err != nil {
		ctx.Server.Logger.Anomalyf(
			"(%s) Score validation error: %v",
			user.Name, err,
		)

		if reject {
			WriteError(http.StatusBadRequest, ValidationError, ctx)
			ctx.Server.Logger.Warningf("Rejected score submission from '%s'.", user.Name)
			return
		}
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

	if score.Passed {
		err = UploadReplay(
			score.Id,
			request.Replay,
			ctx.Server.State.Storage,
		)

		if err != nil {
			ctx.Server.Logger.Warningf("Error uploading replay: %v", err)
		}
	}

	if err = UpdateUserStatistics(request.ScoreData, user, ctx.Server); err != nil {
		ctx.Server.Logger.Warningf("Error updating user statistics: %v", err)
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
