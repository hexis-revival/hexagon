package common

import (
	"database/sql/driver"
	"fmt"
)

type RelationshipStatus string

const (
	StatusFriend  RelationshipStatus = "friend"
	StatusBlocked RelationshipStatus = "blocked"
)

type BeatmapStatus int

const (
	BeatmapStatusUnknown      BeatmapStatus = iota
	BeatmapStatusNotSubmitted BeatmapStatus = iota
	BeatmapStatusPending      BeatmapStatus = iota
	BeatmapStatusRanked       BeatmapStatus = iota
	BeatmapStatusApproved     BeatmapStatus = iota
)

type BeatmapAvailability int

const (
	BeatmapHasDownload             BeatmapAvailability = iota
	BeatmapHasDMCA                 BeatmapAvailability = iota
	BeatmapHasInappropriateContent BeatmapAvailability = iota
)

type ScoreStatus int

const (
	ScoreStatusUnranked  ScoreStatus = iota
	ScoreStatusFailed    ScoreStatus = iota
	ScoreStatusSubmitted ScoreStatus = iota
	ScoreStatusPB        ScoreStatus = iota
)

type Grade int

const (
	GradeF Grade = iota + 1
	GradeD
	GradeC
	GradeB
	GradeA
	GradeS
	GradeSH
	GradeX
	GradeXH
)

var gradeStrings = map[Grade]string{
	GradeF:  "F",
	GradeD:  "D",
	GradeC:  "C",
	GradeB:  "B",
	GradeA:  "A",
	GradeS:  "S",
	GradeSH: "SH",
	GradeX:  "X",
	GradeXH: "XH",
}

func (g Grade) String() string {
	return gradeStrings[g]
}

func (g *Grade) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan Grade: %v", value)
	}

	// Convert the string to the corresponding Grade
	for grade, gradeString := range gradeStrings {
		if gradeString == str {
			*g = grade
			return nil
		}
	}

	return fmt.Errorf("invalid Grade: %s", str)
}

func (g Grade) Value() (driver.Value, error) {
	gradeStr, ok := gradeStrings[g]
	if !ok {
		return nil, fmt.Errorf("invalid Grade: %d", g)
	}
	return gradeStr, nil
}
