package hscore

import (
	"math"
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

type BeatmapSubmissionRequest struct {
	Username      string
	Password      string
	BeatmapIds    []int
	SetId         int
	ClientVersion int
}

func (req *BeatmapSubmissionRequest) String() string {
	return common.FormatStruct(req)
}

type BeatmapSubmissionResponse struct {
	StatusCode int
	SetId      int
	BeatmapIds []int
}

func (resp *BeatmapSubmissionResponse) String() string {
	return common.FormatStruct(resp)
}

func (resp *BeatmapSubmissionResponse) Write() string {
	beatmapIdsStrings := make([]string, 0, len(resp.BeatmapIds))

	for _, beatmapId := range resp.BeatmapIds {
		beatmapIdsStrings = append(beatmapIdsStrings, strconv.Itoa(beatmapId))
	}

	return strings.Join([]string{
		strconv.Itoa(resp.StatusCode),
		strconv.Itoa(resp.SetId),
		strings.Join(beatmapIdsStrings, ":"),
	}, ",")
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
	ClientBuildDate int
	ClientVersion   int
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
