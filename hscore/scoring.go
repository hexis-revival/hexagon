package hscore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/lekuruu/hexagon/common"
)

type ScoreSubmissionResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type ScoreSubmissionRequest struct {
	Replay      *ReplayData
	ProcessList []string
	ScoreData   *ScoreData
	Password    string
	ClientData  string
}

func (req *ScoreSubmissionRequest) String() string {
	return fmt.Sprintf(
		"ScoreSubmissionRequest{%s, ProcessList: %d processes, %s, Password: '%s', ClientData: %s}",
		req.Replay.String(),
		len(req.ProcessList),
		req.ScoreData.String(),
		req.Password,
		req.ClientData,
	)
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
	return fmt.Sprintf(
		"ScoreData{BeatmapChecksum: %s, ScoreChecksum: %s, Username: %s, Passed: %t, Perfect: %t, Time: %d, MaxCombo: %d, TotalScore: %d, Count300: %d, Count100: %d, Count50: %d, CountGeki: %d, CountKatu: %d, CountGood: %d, CountMiss: %d, ClientVersion: %d, Unknown1: %s, Unknown2: %s, %s}",
		scoreData.BeatmapChecksum,
		scoreData.ScoreChecksum,
		scoreData.Username,
		scoreData.Passed,
		scoreData.Perfect,
		scoreData.Time,
		scoreData.MaxCombo,
		scoreData.TotalScore,
		scoreData.Count300,
		scoreData.Count100,
		scoreData.Count50,
		scoreData.CountGeki,
		scoreData.CountKatu,
		scoreData.CountGood,
		scoreData.CountMiss,
		scoreData.ClientVersion,
		scoreData.Unknown1,
		scoreData.Unknown2,
		scoreData.Mods.String(),
	)
}

func (scoreData *ScoreData) Accuracy() float64 {
	totalHits := scoreData.Count300 + scoreData.Count100 + scoreData.Count50
	return float64(scoreData.Count300*300+scoreData.Count100*100+scoreData.Count50*50) / float64(totalHits*300)
}

type Mods struct {
	ArChange  int
	OdChange  int
	CsChange  int
	HpChange  int
	PlaySpeed float32
	Hidden    bool
	NoFail    bool
	Auto      bool
}

func (mods *Mods) String() string {
	return fmt.Sprintf(
		"Mods{ArChange: %d, OdChange: %d, CsChange: %d, HpChange: %d, PlaySpeed: %v, Hidden: %t, NoFail: %t, Auto: %t}",
		mods.ArChange,
		mods.OdChange,
		mods.CsChange,
		mods.HpChange,
		mods.PlaySpeed,
		mods.Hidden,
		mods.NoFail,
		mods.Auto,
	)
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

	iv, err := base64.StdEncoding.DecodeString(ivBase64)
	if err != nil {
		return nil, err
	}

	processList, err := base64.StdEncoding.DecodeString(processListBase64)
	if err != nil {
		return nil, err
	}

	scoreData, err := base64.StdEncoding.DecodeString(scoreDataBase64)
	if err != nil {
		return nil, err
	}

	clientData, err := base64.StdEncoding.DecodeString(clientDataBase64)
	if err != nil {
		return nil, err
	}

	processListDecrypted, err := common.DecryptScoreData(iv, processList)
	if err != nil {
		return nil, err
	}

	scoreDataDecrypted, err := common.DecryptScoreData(iv, scoreData)
	if err != nil {
		return nil, err
	}

	clientDataDecryptedBase64, err := common.DecryptScoreData(iv, clientData)
	if err != nil {
		return nil, err
	}

	clientDataDecrypted, err := base64.StdEncoding.DecodeString(string(clientDataDecryptedBase64))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	scoreDataStruct, err := ParseScoreData(scoreDataDecrypted)
	if err != nil {
		return nil, err
	}

	replayData, _ := ReadCompressedReplay(replay)
	if replayData == nil {
		// Either invalid replay or not provided
		// we will handle this later
		replayData = &ReplayData{}
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

	beatmapChecksum := scoreData[0]
	scoreChecksum := scoreData[3]
	username := scoreData[1]

	passed := scoreData[4] == "1"
	perfect := scoreData[5] == "1"

	unknown1 := scoreData[2]
	unknown2 := scoreData[19]

	time, err := strconv.Atoi(scoreData[6])
	if err != nil {
		return nil, err
	}

	maxCombo, err := strconv.Atoi(scoreData[7])
	if err != nil {
		return nil, err
	}

	totalScore, err := strconv.Atoi(scoreData[8])
	if err != nil {
		return nil, err
	}

	count300, err := strconv.Atoi(scoreData[9])
	if err != nil {
		return nil, err
	}

	count100, err := strconv.Atoi(scoreData[10])
	if err != nil {
		return nil, err
	}

	count50, err := strconv.Atoi(scoreData[11])
	if err != nil {
		return nil, err
	}

	countGeki, err := strconv.Atoi(scoreData[12])
	if err != nil {
		return nil, err
	}

	countKatu, err := strconv.Atoi(scoreData[13])
	if err != nil {
		return nil, err
	}

	countGood, err := strconv.Atoi(scoreData[14])
	if err != nil {
		return nil, err
	}

	countMiss, err := strconv.Atoi(scoreData[15])
	if err != nil {
		return nil, err
	}

	clientVersion, err := strconv.Atoi(scoreData[18])
	if err != nil {
		return nil, err
	}

	mods, err := ParseModsData(scoreData[17])
	if err != nil {
		return nil, err
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

	arChange, err := strconv.Atoi(modData[0])
	if err != nil {
		return nil, err
	}

	odChange, err := strconv.Atoi(modData[1])
	if err != nil {
		return nil, err
	}

	csChange, err := strconv.Atoi(modData[2])
	if err != nil {
		return nil, err
	}

	hpChange, err := strconv.Atoi(modData[3])
	if err != nil {
		return nil, err
	}

	playSpeedMultiplier, err := strconv.Atoi(modData[4])
	if err != nil {
		return nil, err
	}

	playSpeed := 1 + (0.5 * float32(playSpeedMultiplier) / 10)
	hidden := modData[5] == "1"
	noFail := modData[6] == "1"
	auto := modData[7] == "1"

	return &Mods{
		ArChange:  arChange,
		OdChange:  odChange,
		CsChange:  csChange,
		HpChange:  hpChange,
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
	if ctx.Request.Method != "POST" {
		ctx.Response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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
