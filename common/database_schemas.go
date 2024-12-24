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

func (stats *Stats) Clears() int {
	return stats.XHCount + stats.XCount + stats.SHCount + stats.SCount + stats.ACount + stats.BCount + stats.CCount + stats.DCount
}

func (user *User) EnsureStats(state *State) error {
	if user.Stats.UserId != 0 {
		return nil
	}

	// Create new stats object
	user.Stats = Stats{UserId: user.Id}
	return CreateStats(&user.Stats, state)
}

type Relationship struct {
	UserId   int                `gorm:"primaryKey;not null"`
	TargetId int                `gorm:"primaryKey;not null"`
	Status   RelationshipStatus `gorm:"type:relationship_status;not null"`

	User   User `gorm:"foreignKey:UserId"`
	Target User `gorm:"foreignKey:TargetId"`
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
	TopicId            *int                `gorm:"default:null"`

	Beatmaps []Beatmap  `gorm:"foreignKey:SetId"`
	Topic    ForumTopic `gorm:"foreignKey:TopicId"`
	Creator  User       `gorm:"foreignKey:CreatorId"`
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
	Id           int        `gorm:"primaryKey;autoIncrement;not null"`
	ForumId      int        `gorm:"not null"`
	CreatorId    int        `gorm:"not null"`
	Title        string     `gorm:"size:255;not null"`
	StatusText   *string    `gorm:"size:255;default:null"`
	CreatedAt    time.Time  `gorm:"not null;default:now()"`
	LastPostAt   time.Time  `gorm:"not null;default:now()"`
	LockedAt     *time.Time `gorm:"default:null"`
	Views        int        `gorm:"not null;default:0"`
	Announcement bool       `gorm:"not null;default:false"`
	Hidden       bool       `gorm:"not null;default:false"`
	Pinned       bool       `gorm:"not null;default:false"`

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

type Score struct {
	Id            int         `gorm:"primaryKey;autoIncrement;not null"`
	BeatmapId     int         `gorm:"not null"`
	UserId        int         `gorm:"not null"`
	Checksum      string      `gorm:"size:32;not null"`
	Status        ScoreStatus `gorm:"not null"`
	CreatedAt     time.Time   `gorm:"not null;default:now()"`
	ClientVersion int         `gorm:"not null"`
	TotalScore    int64       `gorm:"not null"`
	MaxCombo      int         `gorm:"not null"`
	Accuracy      float64     `gorm:"not null"`
	FullCombo     bool        `gorm:"not null"`
	Passed        bool        `gorm:"not null"`
	Grade         Grade       `gorm:"type:score_grade;not null"`
	Count300      int         `gorm:"not null;column:count_300"`
	Count100      int         `gorm:"not null;column:count_100"`
	Count50       int         `gorm:"not null;column:count_50"`
	CountGeki     int         `gorm:"not null;column:count_geki"`
	CountKatu     int         `gorm:"not null;column:count_katu"`
	CountGood     int         `gorm:"not null;column:count_good"`
	CountMiss     int         `gorm:"not null;column:count_miss"`
	AROffset      int         `gorm:"not null"`
	ODOffset      int         `gorm:"not null"`
	CSOffset      int         `gorm:"not null"`
	HPOffset      int         `gorm:"not null"`
	PSOffset      int         `gorm:"not null"`
	ModHidden     bool        `gorm:"not null"`
	ModNoFail     bool        `gorm:"not null;column:mod_nofail"`
	Visible       bool        `gorm:"not null;default:true"`
	Pinned        bool        `gorm:"not null;default:false"`

	Beatmap Beatmap `gorm:"foreignKey:BeatmapId"`
	User    User    `gorm:"foreignKey:UserId"`
}
