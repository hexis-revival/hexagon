package hscore

import (
	"archive/zip"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math"
	"strconv"
	"strings"

	"github.com/hexis-revival/hexagon/common"
)

type ScoreSubmissionRequest struct {
	ReplayFrames []*common.ReplayFrame
	ProcessList  []string
	ScoreData    *ScoreData
	Password     string
	ClientData   string
}

func (req *ScoreSubmissionRequest) String() string {
	return common.FormatStruct(req)
}

func (req *ScoreSubmissionRequest) ClientDataBase64() string {
	return base64.StdEncoding.EncodeToString([]byte(req.ClientData))
}

type ScoreSubmissionResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (resp *ScoreSubmissionResponse) String() string {
	return common.FormatStruct(resp)
}

type ReplayDownloadRequest struct {
	Username string
	Password string
	ScoreId  int
}

func (req *ReplayDownloadRequest) String() string {
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

type BeatmapUploadRequest struct {
	Username      string
	Password      string
	SetId         int
	ClientVersion int
	Package       *zip.Reader
}

func (req *BeatmapUploadRequest) String() string {
	return common.FormatStruct(req)
}

type BeatmapUploadResponse struct {
	Success bool
}

func (resp *BeatmapUploadResponse) String() string {
	return common.FormatStruct(resp)
}

func (resp *BeatmapUploadResponse) Write() string {
	uploadFailed := 0

	if !resp.Success {
		uploadFailed = 1
	}

	return strconv.Itoa(uploadFailed)
}

type BeatmapDescriptionRequest struct {
	Username      string
	Password      string
	SetId         int
	ClientVersion int
}

func (req *BeatmapDescriptionRequest) String() string {
	return common.FormatStruct(req)
}

type BeatmapDescriptionResponse struct {
	TopicId int
	Content string
}

func (resp *BeatmapDescriptionResponse) String() string {
	return common.FormatStruct(resp)
}

func (resp *BeatmapDescriptionResponse) Write() string {
	return strings.Join([]string{
		strconv.Itoa(resp.TopicId),
		resp.Content,
	}, ",")
}

type BeatmapPostRequest struct {
	Username      string
	Password      string
	SetId         int
	ClientVersion int
	Content       string
}

func (req *BeatmapPostRequest) String() string {
	return common.FormatStruct(req)
}

type BeatmapUpdateRequest struct {
	UserId    int
	SetId     int
	BeatmapId int
}

func (req *BeatmapUpdateRequest) String() string {
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
	ClientBuildDate int
	ClientVersion   int
	Mods            *Mods
}

func (scoreData *ScoreData) String() string {
	return common.FormatStruct(scoreData)
}

func (scoreData *ScoreData) CompareScoreChecksum(clientDataBase64 string) bool {
	expectedChecksum := scoreData.CreateScoreChecksum(clientDataBase64)
	if expectedChecksum == "" {
		return true
	}

	return strings.EqualFold(expectedChecksum, scoreData.ScoreChecksum)
}

func (scoreData *ScoreData) CreateScoreChecksum(clientDataBase64 string) string {
	totalScoreRounded := int(math.Round(float64(scoreData.TotalScore)))

	payload := strings.Join([]string{
		scoreData.BeatmapChecksum,
		"ngc",
		scoreData.Username,
		"dol",
		strconv.Itoa(scoreData.Count300 + scoreData.Count100),
		"w32",
		strconv.Itoa(scoreData.CountMiss),
		"ds",
		strconv.Itoa(scoreData.CountGeki),
		strconv.Itoa(scoreData.Count50),
		"x",
		strconv.Itoa(scoreData.MaxCombo),
		strconv.Itoa(boolToInt(scoreData.Perfect)),
		"rvl",
		strconv.Itoa(totalScoreRounded),
		strconv.Itoa(int(scoreData.Grade())),
		scoreData.Mods.ChecksumToken(),
		strconv.Itoa(boolToInt(scoreData.Passed)),
		"0",
		strconv.Itoa(scoreData.Time),
		strconv.Itoa(scoreData.ClientBuildDate),
		"snes",
		clientDataBase64,
	}, "")

	checksum := md5.Sum([]byte(payload))
	return hex.EncodeToString(checksum[:])
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func (scoreData *ScoreData) PassedObjects() int {
	return scoreData.Count300 + scoreData.Count100 + scoreData.Count50 + scoreData.CountMiss
}

func (scoreData *ScoreData) TotalHits() int {
	return scoreData.Count300 + scoreData.Count100 + scoreData.Count50
}

func (scoreData *ScoreData) Accuracy() float64 {
	return float64(scoreData.Count300*300+scoreData.Count100*100+scoreData.Count50*50) / float64(scoreData.TotalHits()*300)
}

func (scoreData *ScoreData) Grade() common.Grade {
	return common.CalculateGrade(
		scoreData.Passed,
		scoreData.Count300,
		scoreData.Count100,
		scoreData.Count50,
		scoreData.CountMiss,
		scoreData.Mods.Hidden,
	)
}

type Mods struct {
	ArOffset int
	OdOffset int
	CsOffset int
	HpOffset int
	PsOffset int
	Hidden   bool
	NoFail   bool
	Auto     bool
}

func (mods *Mods) String() string {
	return common.FormatStruct(mods)
}

func (mods *Mods) ChecksumToken() string {
	return common.CreateModsChecksumToken(
		mods.ArOffset,
		mods.OdOffset,
		mods.CsOffset,
		mods.HpOffset,
		mods.PsOffset,
		mods.Hidden,
		mods.NoFail,
		mods.Auto,
	)
}
