package pkg

import (
	"context"
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

// Common errors returned throughout the service.
var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
)

// BettingService represents the service implementing how to bet on teams.
type BettingService interface {
	AddCompetition(ctx context.Context, competition *Competition) (*Competition, error)
	AddCompetitor(ctx context.Context, competitor *Competitor, bindToCompetitionID *int) (*Competitor, error)
	AddBetter(ctx context.Context, better *Better) (*Better, error)
	AddBet(ctx context.Context, bet *Bet) (*Bet, error)

	GetCompetition(ctx context.Context, id int) (*Competition, error)
	GetCompetitions(ctx context.Context, ids []int) ([]*Competition, error)
	GetCompetitor(ctx context.Context, id int) (*Competitor, error)
	GetCompetitors(ctx context.Context, ids []int) ([]*Competitor, error)
	GetBetter(ctx context.Context, id int) (*Better, error)
	GetBetters(ctx context.Context, ids []int) ([]*Better, error)
	GetBet(ctx context.Context, id int) (*Bet, error)
	GetBets(ctx context.Context, ids []int) ([]*Bet, error)

	DeleteCompetition(ctx context.Context, id int) error
	DeleteCompetitor(ctx context.Context, id int) error
	DeleteBetter(ctx context.Context, id int) error
	DeleteBet(ctx context.Context, id int) error

	GetCompetitionMetrics(ctx context.Context, id int) (*CompetitionMetrics, error)
	GetCompetitorsForCompetition(ctx context.Context, id int) ([]*Competitor, error)
	GetBetsForCompetition(ctx context.Context, id int) ([]*Bet, error)
	GetCreatedObjectsForBetter(ctx context.Context, id int) ([]*Competition, []*Competitor, []*Bet, error)
}

// MetricValue represents who has what value, e.g. who has the lowest average
// and what the value is.
type MetricValue struct {
	Who   *Better
	Value interface{}
}

// CompetitionMetrics represents metrics that may be calculated for a
// competition.
type CompetitionMetrics struct {
	HighestAverageBetter MetricValue
	LowestAverageBetter  MetricValue
	MostTopScores        MetricValue
	MostBottomScores     MetricValue
	LongestNote          MetricValue
	ShortestNote         MetricValue
	NumberOfBottomScores int
	NumberOfTopScores    int
	GroupAverageScore    float64
}

// Competition represents one competition, e.g. Eurovision Song Contest 2022.
type Competition struct {
	ID          int                 `db:"id"          json:"id"            gorm:"primary_key"`
	CreatedAt   time.Time           `db:"created_at"  json:"created_at"`
	UpdatedAt   time.Time           `db:"updated_at"  json:"updated_at"`
	DeletedAt   null.Time           `db:"deleted_at"  json:"deleted_at"`
	CreatedBy   *Better             `db:"-"           json:"created_by"    gorm:"foreignkey:CreatedByID"`
	CreatedByID int                 `db:"created_by"  json:"created_by_id" gorm:"not null"`
	Name        string              `db:"name"        json:"name"          gorm:"type:varchar(100); not null"`
	Description null.String         `db:"description" json:"description"   gorm:"type:varchar(255)"`
	Code        null.String         `db:"code"        json:"code"          gorm:"code:varchar(10)"`
	Image       null.String         `db:"image"       json:"image"         gorm:"type:varchar(100)"`
	MinScore    int                 `db:"min_score"   json:"min_score"     gorm:"type:int; not null"`
	MaxScore    int                 `db:"max_score"   json:"max_score"     gorm:"type:int; not null"`
	Locked      bool                `db:"locked"      json:"locked"        gorm:"type:tinyint(1); default 0"`
	Metrics     *CompetitionMetrics `db:"-"           json:"metrics"       gorm:"-"`
	Competitors []*Competitor       `db:"-"           json:"competitors"   gorm:"many2many:competition_competitor"`
	Bets        []*Bet              `db:"-"           json:"bets"`
}

// Competitor represents a team or player competing in a competition. A
// Competitor may compete in zero or several Competitions.
type Competitor struct {
	ID           int            `db:"id"          json:"id"            gorm:"primary_key"`
	CreatedAt    time.Time      `db:"created_at"  json:"created_at"`
	UpdatedAt    null.Time      `db:"updated_at"  json:"updated_at"`
	DeletedAt    null.Time      `db:"deleted_at"  json:"deleted_at"`
	CreatedBy    *Better        `db:"-"           json:"created_by"    gorm:"foreignkey:CreatedByID"`
	CreatedByID  int            `db:"created_by"  json:"created_by_id" gorm:"not null"`
	Name         string         `db:"name"        json:"name"          gorm:"type:varchar(100); not null"`
	Description  null.String    `db:"description" json:"description"   gorm:"type:varchar(255)"`
	Image        null.String    `db:"image"       json:"image"         gorm:"type:varchar(100)"`
	Competitions []*Competition `db:"-"           json:"competitions"  gorm:"many2many:competition_competitor"`
}

// Better is someone who can make a Bet on a Competitor.
type Better struct {
	ID        int         `db:"id"         json:"id"         gorm:"primary_key"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt null.Time   `db:"updated_at" json:"updated_at"`
	DeletedAt null.Time   `db:"deleted_at" json:"deleted_at"`
	Confirmed bool        `db:"confirmed"  json:"confirmed"  gorm:"type:tinyint(1); default 0"`
	Name      string      `db:"name"       json:"name"       gorm:"type:varchar(100); not null"`
	Email     string      `db:"email"      json:"email"      gorm:"type:varchar(100); not null; unique"`
	Image     null.String `db:"image"      json:"image"      gorm:"type:varchar(100)"`
}

// Bet is a bet put on a Competitor in a certain Competition.
type Bet struct {
	ID            int          `db:"id"                        json:"id"                     gorm:"primary_key"`
	CreatedAt     time.Time    `db:"created_at"                json:"created_at"`
	UpdatedAt     null.Time    `db:"updated_at"                json:"updated_at"`
	Score         null.Int     `db:"score"                     json:"score"                  gorm:"type:int"`
	Placing       null.Int     `db:"placing"                   json:"placing"                gorm:"type:int"`
	Note          null.String  `db:"note"                      json:"note"                   gorm:"type:varchar(255)"`
	Better        *Better      `db:"-"                         json:"better"`
	BetterID      int          `db:"better_id"                 json:"better_id"              gorm:"unique_index:idx_better_id_competition_id_competitor_id; not null"`
	Competition   *Competition `db:"-"                         json:"competition"`
	CompetitionID int          `db:"competition_id"            json:"competition_id"         gorm:"unique_index:idx_better_id_competition_id_competitor_id; not null"`
	Competitor    *Competitor  `db:"-"                         json:"competitor"`
	CompetitorID  int          `db:"competitor_id"             json:"competitor_id"          gorm:"unique_index:idx_better_id_competition_id_competitor_id; not null"`
}
