package hscore

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hexis-revival/hexagon/common"
)

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

	clientVersion := scoreData[2]
	clientVersionCheck := scoreData[19]

	if clientVersion != clientVersionCheck {
		return nil, fmt.Errorf("client version mismatch: %s != %s", clientVersion, clientVersionCheck)
	}

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

	clientBuildDate, err := strconv.Atoi(scoreData[18])
	collection.Add(err)

	clientVersionInt, err := strconv.Atoi(clientVersion)
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
		ClientBuildDate: clientBuildDate,
		ClientVersion:   clientVersionInt,
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
