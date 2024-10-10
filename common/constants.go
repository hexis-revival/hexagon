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
	GradeSS
	GradeXH
)

type RelationshipStatus string

const (
	StatusFriend  RelationshipStatus = "friend"
	StatusBlocked RelationshipStatus = "blocked"
)
