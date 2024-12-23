package common

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
