package common

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
)

func CalculateGrade(passed bool, count300, count100, count50, countMiss int, hidden bool) Grade {
	if !passed {
		return GradeF
	}

	totalHits := count300 + count100 + count50 + countMiss
	if totalHits == 0 {
		return GradeF
	}

	accuracyRatio := float64(count300) / float64(totalHits)
	count50Ratio := float64(count50) / float64(totalHits)

	switch {
	case accuracyRatio == 1.0:
		if hidden {
			return GradeXH
		}
		return GradeX
	case accuracyRatio > 0.9 && count50Ratio <= 0.01 && countMiss == 0:
		if hidden {
			return GradeSH
		}
		return GradeS
	case (accuracyRatio > 0.8 && countMiss == 0) || accuracyRatio > 0.9:
		return GradeA
	case (accuracyRatio <= 0.7 || countMiss != 0) && accuracyRatio <= 0.8:
		if accuracyRatio > 0.6 {
			return GradeC
		}
		return GradeD
	default:
		return GradeB
	}
}

func CreateReplayChecksum(
	playerName string,
	beatmapChecksum string,
	passed bool,
	count300 int,
	count100 int,
	count50 int,
	countGeki int,
	countGood int,
	countMiss int,
	maxCombo int,
	fullCombo bool,
	totalScore int,
	grade Grade,
	modsToken string,
) string {
	payload := strings.Join(
		[]string{
			playerName,
			beatmapChecksum,
			strconv.Itoa(int(BooleanToInteger(passed))),
			strconv.Itoa(count300 + count100),
			strconv.Itoa(count50),
			strconv.Itoa(countGeki),
			strconv.Itoa(countGood),
			strconv.Itoa(countMiss),
			strconv.Itoa(maxCombo),
			strconv.Itoa(int(BooleanToInteger(fullCombo))),
			strconv.Itoa(totalScore),
			strconv.Itoa(int(grade)),
			modsToken,
		},
		"|",
	)

	hash := md5.Sum([]byte(payload))
	return hex.EncodeToString(hash[:])
}

func CreateModsChecksumToken(ar, od, cs, hp, ps int, hidden, noFail, autoplay bool) string {
	hiddenDigit := "3"
	if hidden {
		hiddenDigit = "2"
	}

	noFailDigit := "0"
	if noFail {
		noFailDigit = "5"
	}

	autoplayDigit := "1"
	if autoplay {
		autoplayDigit = "6"
	}

	return strings.Join([]string{
		strconv.Itoa(2 * (23 - ar)),
		"1",
		strconv.Itoa(5 * (61 - od)),
		"3",
		strconv.Itoa(3 * (47 - cs)),
		"5",
		strconv.Itoa(17 - hp),
		"9",
		strconv.Itoa(7 * (100 - ps)),
		"7",
		hiddenDigit,
		noFailDigit,
		autoplayDigit,
	}, "")
}
