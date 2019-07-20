package pkg

// Constants for table names in the data model.
const (
	CompetitionTable           = "competition"
	CompetitorTable            = "competitor"
	CompetitionCompetitorTable = "competition_competitor"
	BetterTable                = "better"
	BetTable                   = "bet"
)

// Competition represents one competition, e.g. Eurovision Song Contest 2022.
type Competition struct {
	ID          int           `db:"id"          json:"id"`
	Name        string        `db:"name"        json:"name"`
	Description string        `db:"description" json:"description"`
	Image       string        `db:"image"       json:"image"`
	Competitors []*Competitor `                 json:"competitors"`
}

// Competitor represents a team or player competing in a competition. A
// Competitor may compete in zero or several Competitions.
type Competitor struct {
	ID            int          `db:"id"             json:"id"`
	CompetitionID int          `db:"competition_id" json:"-"`
	Name          string       `db:"name"           json:"name"`
	Description   string       `db:"description"    json:"description"`
	Image         string       `db:"image"          json:"image"`
	Competition   *Competition `                    json:"competition"`
}

// Better is someone who can make a Bet on a Competitor.
type Better struct {
	id    int    `db:"id"    json:"id"`
	Image string `db:"image" json:"image"`
	Name  string `db:"name"  json:"name"`
}

// Bet is a bet put on a Competitor in a certain Competition.
type Bet struct {
	ID                      int          `db:"id"                        json:"id"`
	BetterID                int          `db:"id_better"                 json:"-"`
	CompetitionCompetitorID int          `db:"id_competition_competitor" json:"-"`
	Placing                 int          `db:"placing"                   json:"placing"`
	Note                    string       `db:"note"                      json:"note"`
	Better                  *Better      `                               json:"better"`
	Competition             *Competition `                               json:"competition"`
	Competitor              *Competitor  `                               json:"competitor"`
}
