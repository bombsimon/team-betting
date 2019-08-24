package pkg

import (
	"time"

	"github.com/guregu/null"
	"github.com/pkg/errors"
)

// Constants for table names in the data model.
const (
	BetTable                   = "bet"
	BetterTable                = "better"
	CompetitionCompetitorTable = "competition_competitor"
	CompetitionTable           = "competition"
	CompetitorTable            = "competitor"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
)

// CompetitionMetrics represents metrics that may be calculated for a
// competition.
type CompetitionMetrics struct {
	HighestAverageBetter *Better
	LowestAverageBetter  *Better
	MostTopScores        *Better
	MostBottomScores     *Better
	LongestNotes         *Better
	ShortestNotes        *Better
	NumberOfBottomScores int
	NumberOfTopScores    int
	GroupAverageScore    int
}

// Competition represents one competition, e.g. Eurovision Song Contest 2022.
type Competition struct {
	ID          int         `json:"id"           gorm:"primary_key"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   null.Time   `json:"deleted_at"`
	Name        string      `json:"name"         gorm:"type:varchar(100); not null"`
	Description null.String `json:"description"  gorm:"type:varchar(255)"`
	Image       null.String `json:"image"        gorm:"type:varchar(100)`
	MinScore    int         `json:"min_score"    gorm:"type:int; not null"`
	MaxScore    int         `json:"max_score"    gorm:"type:int; not null"`
	Locked      bool        `json:"locked"       gorm:"type:tinyint(1); default 0"`
}

// Competitor represents a team or player competing in a competition. A
// Competitor may compete in zero or several Competitions.
type Competitor struct {
	ID          int         `json:"id"          gorm:"primary_key"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   null.Time   `json:"updated_at"`
	DeletedAt   null.Time   `json:"deleted_at"`
	Name        string      `json:"name"        gorm:"type:varchar(100); not null"`
	Description null.String `json:"description" gorm:"type:varchar(255)"`
	Image       null.String `json:"image"       gorm:"type:varchar(100)`
}

// CompetitionCompetitor binds a competitor to a competition.
type CompetitionCompetitor struct {
	ID            int         `json:"id"          gorm:"primary_key"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     null.Time   `json:"updated_at"`
	CompetitionID int         `json:"-"           gorm:"column:id_competition; not null"`
	Competition   Competition `json:"competition"`
	CompetitorID  int         `json:"-"           gorm:"column:id_competitor; not null"`
	Competitor    Competitor  `json:"competitor"  gorm:"foreignkey:CompetitorID"`
}

// Better is someone who can make a Bet on a Competitor.
type Better struct {
	ID        int         `json:"id"         gorm:"primary_key"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt null.Time   `json:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at"`
	Confirmed bool        `json:"confirmed"  gorm:"type:tinyint(1); default 0"`
	Name      string      `json:"name"       gorm:"type:varchar(100); not null"`
	Email     string      `json:"email"      gorm:"type:varchar(100); not null; unique"`
	Image     null.String `json:"image"      gorm:"type:varchar(100)`
}

// Bet is a bet put on a Competitor in a certain Competition.
type Bet struct {
	ID                      int                   `json:"id"                     gorm:"primary_key"`
	CreatedAt               time.Time             `json:"created_at"`
	UpdatedAt               null.Time             `json:"updated_at"`
	Score                   null.Int              `json:"score"                  gorm:"type:int"`
	Placing                 null.Int              `json:"placing"                gorm:"type:int`
	Note                    null.String           `json:"note"                   gorm:"type:varchar(255)"`
	Better                  Better                `json:"better"`
	BetterID                int                   `json:"-"                      gorm:"column:id_better; not null"`
	CompetitionCompetitor   CompetitionCompetitor `json:"competition_competitor"`
	CompetitionCompetitorID int                   `json:"-"                      gorm:"column:id_competition_competitor; not null"`
}
