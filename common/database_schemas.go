package common

import (
	"time"

	"github.com/lib/pq"
)

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

type Beatmapset struct {
	Id                 int                 `gorm:"primaryKey;autoIncrement;not null"`
	Title              string              `gorm:"size:255;not null"`
	Artist             string              `gorm:"size:255;not null"`
	Source             string              `gorm:"size:255;not null"`
	Tags               pq.StringArray      `gorm:"type:text[];not null;default:'{}'"`
	CreatorId          int                 `gorm:"not null"`
	CreatedAt          time.Time           `gorm:"not null;default:now()"`
	LastUpdated        time.Time           `gorm:"not null;default:now()"`
	ApprovedAt         *time.Time          `gorm:"default:null"`
	ApprovedBy         *int                `gorm:"default:null"`
	Status             BeatmapStatus       `gorm:"not null;default:1"`
	Description        string              `gorm:"type:text;not null"`
	HasVideo           bool                `gorm:"not null;default:false"`
	AvailabilityStatus BeatmapAvailability `gorm:"not null;default:0"`
	AvailabilityInfo   string              `gorm:"type:text;not null;default:''"`

	Beatmaps []Beatmap `gorm:"foreignKey:SetId"`
	Creator  User      `gorm:"foreignKey:CreatorId"`
}

type Beatmap struct {
	Id            int           `gorm:"primaryKey;autoIncrement;not null"`
	SetId         int           `gorm:"not null"`
	Checksum      string        `gorm:"size:32;not null"`
	Version       string        `gorm:"size:255;not null"`
	Filename      string        `gorm:"size:512;not null"`
	CreatorId     int           `gorm:"not null"`
	CreatedAt     time.Time     `gorm:"not null;default:now()"`
	LastUpdated   time.Time     `gorm:"not null;default:now()"`
	Status        BeatmapStatus `gorm:"not null"`
	TotalLength   int           `gorm:"not null;default:0"`
	DrainLength   int           `gorm:"not null;default:0"`
	TotalCircles  int           `gorm:"not null;default:0"`
	TotalSliders  int           `gorm:"not null;default:0"`
	TotalSpinners int           `gorm:"not null;default:0"`
	TotalHolds    int           `gorm:"not null;default:0"`
	MaxCombo      int           `gorm:"not null;default:0"`
	MedianBpm     float64       `gorm:"not null;default:0"`
	HighestBpm    float64       `gorm:"not null;default:0"`
	LowestBpm     float64       `gorm:"not null;default:0"`
	CS            float64       `gorm:"not null;default:0"`
	HP            float64       `gorm:"not null;default:0"`
	OD            float64       `gorm:"not null;default:0"`
	AR            float64       `gorm:"not null;default:0"`
	SR            float64       `gorm:"not null;default:0"`

	Set     Beatmapset `gorm:"foreignKey:SetId"`
	Creator User       `gorm:"foreignKey:CreatorId"`
}

type Forum struct {
	Id          int       `gorm:"primaryKey;autoIncrement;not null"`
	ParentId    *int      `gorm:"default:null"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	Name        string    `gorm:"size:32;not null"`
	Description string    `gorm:"size:255;not null;default:''"`
	Hidden      bool      `gorm:"not null;default:false"`

	Parent *Forum `gorm:"foreignKey:ParentId"`
}

type ForumTopic struct {
	Id           int       `gorm:"primaryKey;autoIncrement;not null"`
	ForumId      int       `gorm:"not null"`
	CreatorId    int       `gorm:"not null"`
	Title        string    `gorm:"size:255;not null"`
	StatusText   *string   `gorm:"size:255;default:null"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	LastPostAt   time.Time `gorm:"not null;default:now()"`
	LockedAt     *time.Time
	Views        int  `gorm:"not null;default:0"`
	Announcement bool `gorm:"not null;default:false"`
	Hidden       bool `gorm:"not null;default:false"`
	Pinned       bool `gorm:"not null;default:false"`

	Forum   Forum `gorm:"foreignKey:ForumId"`
	Creator User  `gorm:"foreignKey:CreatorId"`
}

type ForumPost struct {
	Id         int       `gorm:"primaryKey;autoIncrement;not null"`
	TopicId    int       `gorm:"not null"`
	ForumId    int       `gorm:"not null"`
	UserId     int       `gorm:"not null"`
	Content    string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"not null;default:now()"`
	EditTime   time.Time `gorm:"not null;default:now()"`
	EditCount  int       `gorm:"not null;default:0"`
	EditLocked bool      `gorm:"not null;default:false"`
	Hidden     bool      `gorm:"not null;default:false"`
	Deleted    bool      `gorm:"not null;default:false"`

	Topic ForumTopic `gorm:"foreignKey:TopicId"`
	Forum Forum      `gorm:"foreignKey:ForumId"`
	User  User       `gorm:"foreignKey:UserId"`
}
