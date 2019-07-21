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

// Competition represents one competition, e.g. Eurovision Song Contest 2022.
type Competition struct {
	ID          int           `db:"id"          json:"id"`
	CreatedAt   time.Time     `db:"created_at"  json:"created_at"`
	UpdatedAt   null.Time     `db:"updated_at"  json:"updated_at"`
	Name        string        `db:"name"        json:"name"`
	Description null.String   `db:"description" json:"description"`
	Image       null.String   `db:"image"       json:"image"`
	Locked      bool          `db:"locked"      json:"locked"`
	Competitors []*Competitor `db:"-"           json:"competitors"`
}

// Competitor represents a team or player competing in a competition. A
// Competitor may compete in zero or several Competitions.
type Competitor struct {
	ID          int          `db:"id"             json:"id"`
	CreatedAt   time.Time    `db:"created_at"     json:"created_at"`
	UpdatedAt   null.Time    `db:"updated_at"     json:"updated_at"`
	Name        string       `db:"name"           json:"name"`
	Description null.String  `db:"description"    json:"description"`
	Image       null.String  `db:"image"          json:"image"`
	Competition *Competition `db:"-"              json:"competition"`
}

// CompetitionCompetitor represents the linking between a competition and a
// competitor.
type CompetitionCompetitor struct {
	ID            int `db:"id"`
	CompetitionID int `db:"id_competition"`
	CompetitorID  int `db:"id_competitor"`
}

// Better is someone who can make a Bet on a Competitor.
type Better struct {
	ID        int         `db:"id"         json:"id"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt null.Time   `db:"updated_at" json:"updated_at"`
	Name      string      `db:"name"       json:"name"`
	Email     string      `db:"email"      json:"email"`
	Confirmed bool        `db:"confirmed"  json:"confirmed"`
	Image     null.String `db:"image"      json:"image"`
}

// Bet is a bet put on a Competitor in a certain Competition.
type Bet struct {
	ID                      int          `db:"id"                        json:"id"`
	CreatedAt               time.Time    `db:"created_at"                json:"created_at"`
	UpdatedAt               null.Time    `db:"updated_at"                json:"updated_at"`
	BetterID                int          `db:"id_better"                 json:"-"`
	CompetitionCompetitorID int          `db:"id_competition_competitor" json:"-"`
	Placing                 null.Int     `db:"placing"                   json:"placing"`
	Note                    null.String  `db:"note"                      json:"note"`
	Better                  *Better      `db:"-"                         json:"better"`
	Competition             *Competition `db:"-"                         json:"competition"`
	Competitor              *Competitor  `db:"-"                         json:"competitor"`
}
