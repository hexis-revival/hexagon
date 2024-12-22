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
	StatusUnknown      BeatmapStatus = iota
	StatusNotSubmitted BeatmapStatus = iota
	StatusPending      BeatmapStatus = iota
	StatusRanked       BeatmapStatus = iota
	StatusApproved     BeatmapStatus = iota
)

type BeatmapAvailability int

const (
	BeatmapHasDownload             BeatmapAvailability = iota
	BeatmapHasDMCA                 BeatmapAvailability = iota
	BeatmapHasInappropriateContent BeatmapAvailability = iota
)
