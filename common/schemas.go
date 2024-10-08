package common

import "time"

type User struct {
	Id             int       `gorm:"primaryKey;autoIncrement;not null"`
	Name           string    `gorm:"size:32;not null"`
	Email          string    `gorm:"size:255;not null"`
	Password       string    `gorm:"size:60;not null"`
	Country        string    `gorm:"size:2;default:'XX';not null"`
	CreatedAt      time.Time `gorm:"not null;default:now()"`
	LatestActivity time.Time `gorm:"not null;default:now()"`
	Restricted     bool      `gorm:"not null;default:false"`
	Activated      bool      `gorm:"not null;default:false"`

	Stats Stats `gorm:"foreignKey:UserId"`
}

type Stats struct {
	UserId      int     `gorm:"primaryKey;not null"`
	Rank        int     `gorm:"not null;default:0"`
	TotalScore  int64   `gorm:"not null;default:0"`
	RankedScore int64   `gorm:"not null;default:0"`
	Playcount   int     `gorm:"not null;default:0"`
	Playtime    int     `gorm:"not null;default:0"`
	Accuracy    float64 `gorm:"not null;default:0.0000"`
	MaxCombo    int     `gorm:"not null;default:0"`
	TotalHits   int64   `gorm:"not null;default:0"`
	XHCount     int     `gorm:"not null;default:0"`
	XCount      int     `gorm:"not null;default:0"`
	SHCount     int     `gorm:"not null;default:0"`
	SCount      int     `gorm:"not null;default:0"`
	ACount      int     `gorm:"not null;default:0"`
	BCount      int     `gorm:"not null;default:0"`
	CCount      int     `gorm:"not null;default:0"`
	DCount      int     `gorm:"not null;default:0"`
}

type RelationshipStatus string

const (
	StatusFriend  RelationshipStatus = "friend"
	StatusBlocked RelationshipStatus = "blocked"
)

type Relationship struct {
	UserId   int                `gorm:"primaryKey;not null"`
	TargetId int                `gorm:"primaryKey;not null"`
	Status   RelationshipStatus `gorm:"type:relationship_status;not null"`

	User   User `gorm:"foreignKey:UserId"`
	Target User `gorm:"foreignKey:TargetId"`
}

func (user *User) EnsureStats(state *State) error {
	if user.Stats.UserId != 0 {
		return nil
	}

	// Create new stats object
	user.Stats = Stats{UserId: user.Id}
	return CreateStats(&user.Stats, state)
}
