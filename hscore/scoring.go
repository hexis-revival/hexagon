package hscore

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/hexis-revival/hexagon/common"
)

type ScoreSubmissionResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type ScoreSubmissionRequest struct {
	Replay      *common.ReplayData
	ProcessList []string
	ScoreData   *ScoreData
	Password    string
	ClientData  string
}

func (req *ScoreSubmissionRequest) String() string {
	return common.FormatStruct(req)
}

type ScoreData struct {
	BeatmapChecksum string
	ScoreChecksum   string
	Username        string
	Passed          bool
	Perfect         bool
	Time            int
	MaxCombo        int
	TotalScore      int
	Count300        int
	Count100        int
	Count50         int
	CountGeki       int
	CountKatu       int
	CountGood       int
	CountMiss       int
	ClientVersion   int
	Unknown1        string // TODO
	Unknown2        string // TODO
	Mods            *Mods
}

func (scoreData *ScoreData) String() string {
	return common.FormatStruct(scoreData)
}

func (scoreData *ScoreData) TotalHits() int {
	return scoreData.Count300 + scoreData.Count100 + scoreData.Count50
}

func (scoreData *ScoreData) Accuracy() float64 {
	return float64(scoreData.Count300*300+scoreData.Count100*100+scoreData.Count50*50) / float64(scoreData.TotalHits()*300)
}

func (scoreData *ScoreData) Grade() common.Grade {
	totalHits := scoreData.Count300 + scoreData.Count100 + scoreData.Count50 + scoreData.CountGood

	if totalHits == 0 {
		return common.GradeF
	}

	totalHitCount := float64(totalHits)
	accuracyRatio := float64(scoreData.Count300) / totalHitCount

	if !scoreData.Passed {
		return common.GradeF
	}

	if math.IsNaN(accuracyRatio) || accuracyRatio == 1.0 {
		if scoreData.Mods.Hidden {
			return common.GradeXH
		} else {
			return common.GradeSS
		}
	}

	if accuracyRatio <= 0.8 && scoreData.CountGood == 0 {
		if accuracyRatio > 0.6 {
			return common.GradeC
		}
		return common.GradeD
	}

	if accuracyRatio <= 0.9 {
		return common.GradeB
	}

	return common.GradeA
}

type Mods struct {
	ArOffset  int
	OdOffset  int
	CsOffset  int
	HpOffset  int
	PlaySpeed float32
	Hidden    bool
	NoFail    bool
	Auto      bool
}

func (mods *Mods) String() string {
	return common.FormatStruct(mods)
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
	if replayData == nil {
		// Either invalid replay or not provided
		// we will handle this later
		replayData = &common.ReplayData{}
	}

	processListData := ParseProcessList(processListDecrypted)

	// TODO: Parse process list, score data & client data
	return &ScoreSubmissionRequest{
		Replay:      replayData,
		Password:    password,
		ProcessList: processListData,
		ScoreData:   scoreDataStruct,
		ClientData:  string(clientDataDecrypted),
	}, nil
}

func ParseScoreData(scoreDataBytes []byte) (*ScoreData, error) {
	scoreData := strings.Split(string(scoreDataBytes), ";")

	if len(scoreData) != 20 {
		return nil, fmt.Errorf("invalid score data: %d fields", len(scoreData))
	}

	beatmapChecksum := scoreData[0]
	scoreChecksum := scoreData[3]
	username := scoreData[1]

	passed := scoreData[4] == "1"
	perfect := scoreData[5] == "1"

	unknown1 := scoreData[2]
	unknown2 := scoreData[19]

	collection := common.NewErrorCollection()

	time, err := strconv.Atoi(scoreData[6])
	collection.Add(err)

	maxCombo, err := strconv.Atoi(scoreData[7])
	collection.Add(err)

	totalScore, err := strconv.Atoi(scoreData[8])
	collection.Add(err)

	count300, err := strconv.Atoi(scoreData[9])
	collection.Add(err)

	count100, err := strconv.Atoi(scoreData[10])
	collection.Add(err)

	count50, err := strconv.Atoi(scoreData[11])
	collection.Add(err)

	countGeki, err := strconv.Atoi(scoreData[12])
	collection.Add(err)

	countKatu, err := strconv.Atoi(scoreData[13])
	collection.Add(err)

	countGood, err := strconv.Atoi(scoreData[14])
	collection.Add(err)

	countMiss, err := strconv.Atoi(scoreData[15])
	collection.Add(err)

	clientVersion, err := strconv.Atoi(scoreData[18])
	collection.Add(err)

	mods, err := ParseModsData(scoreData[17])
	collection.Add(err)

	if collection.HasErrors() {
		return nil, collection.Pop(0)
	}

	return &ScoreData{
		BeatmapChecksum: beatmapChecksum,
		ScoreChecksum:   scoreChecksum,
		Username:        username,
		Passed:          passed,
		Perfect:         perfect,
		Time:            time,
		MaxCombo:        maxCombo,
		TotalScore:      totalScore,
		Count300:        count300,
		Count100:        count100,
		Count50:         count50,
		CountGeki:       countGeki,
		CountKatu:       countKatu,
		CountGood:       countGood,
		CountMiss:       countMiss,
		ClientVersion:   clientVersion,
		Unknown1:        unknown1,
		Unknown2:        unknown2,
		Mods:            mods,
	}, nil
}

func ParseModsData(modsString string) (*Mods, error) {
	modData := strings.Split(modsString, ":")
	collection := common.NewErrorCollection()

	arOffset, err := strconv.Atoi(modData[0])
	collection.Add(err)

	odOffset, err := strconv.Atoi(modData[1])
	collection.Add(err)

	csOffset, err := strconv.Atoi(modData[2])
	collection.Add(err)

	hpOffset, err := strconv.Atoi(modData[3])
	collection.Add(err)

	playSpeedMultiplier, err := strconv.Atoi(modData[4])
	collection.Add(err)

	if collection.HasErrors() {
		return nil, collection.Pop(0)
	}

	playSpeed := 1 + (0.5 * float32(playSpeedMultiplier) / 10)
	hidden := modData[5] == "1"
	noFail := modData[6] == "1"
	auto := modData[7] == "1"

	return &Mods{
		ArOffset:  arOffset,
		OdOffset:  odOffset,
		CsOffset:  csOffset,
		HpOffset:  hpOffset,
		PlaySpeed: playSpeed,
		Hidden:    hidden,
		NoFail:    noFail,
		Auto:      auto,
	}, nil
}

func ParseProcessList(processListBytes []byte) []string {
	processListStr := strings.ReplaceAll(string(processListBytes), "\n", "")
	processListStr = strings.ReplaceAll(processListStr, "; ", "")
	processList := strings.Split(processListStr, "| ")
	return processList
}

func ScoreSubmissionHandler(ctx *Context) {
	// Parse score submission request
	req, err := NewScoreSubmissionRequest(ctx.Request)

	if err != nil {
		ctx.Server.Logger.Errorf("Error parsing score submission request: %v", err)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx.Server.Logger.Debugf("Score submission request: %s", req.String())
	// TODO: Process data & update player stats

	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Header().Set("Content-Type", "application/json")

	// Write response
	resp := ScoreSubmissionResponse{Success: true}
	json.NewEncoder(ctx.Response).Encode(resp)
}
