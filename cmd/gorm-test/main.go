package main

import (
	"github.com/bombsimon/team-betting/pkg"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/guregu/null"
	"github.com/jinzhu/gorm"
)

func main() {
	db, err := gorm.Open("mysql", "betting:betting@/betting?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	// Add all auto migrations.
	db.AutoMigrate(&pkg.Competition{})
	db.AutoMigrate(&pkg.Competitor{})
	db.AutoMigrate(&pkg.Better{})

	// Add foreign keys keys, not done bu auto migrations
	// https://github.com/jinzhu/gorm/issues/450
	db.AutoMigrate(&pkg.CompetitionCompetitor{}).
		AddForeignKey("id_competition", "competitions(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("id_competitor", "competitors(id)", "RESTRICT", "RESTRICT")

	db.AutoMigrate(&pkg.Bet{}).
		AddForeignKey("id_better", "betters(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("id_competition_competitor", "competition_competitors(id)", "RESTRICT", "RESTRICT")

	// Create test data
	competition := &pkg.Competition{
		Name:        "Eurovision Song Contest 2020",
		Description: null.StringFrom("The one that started it all"),
	}

	competitors := []*pkg.Competitor{
		{Name: "Sweden - Swedish song"},
		{Name: "Norway - Norwegian song"},
	}

	betters := []*pkg.Better{
		{Name: "Testy Testsson", Email: "testy@testsson.se", Confirmed: true},
		{Name: "Another Tester", Email: "testy@anotherone.se"},
	}

	db.Create(competition)

	for _, c := range competitors {
		db.Create(c)

		db.Create(&pkg.CompetitionCompetitor{
			Competition: *competition,
			Competitor:  *c,
		})
	}

	for _, b := range betters {
		db.Create(b)
	}

	var cc pkg.CompetitionCompetitor
	db.First(&cc)

	db.Create(&pkg.Bet{
		Better:                *betters[0],
		CompetitionCompetitor: cc,
		Score:                 null.IntFrom(6),
	})

	// Fetch the first bet
	var bet pkg.Bet
	db.First(&bet)

	// Find related - must specify field
	// https://github.com/jinzhu/gorm/issues/2615
	db.Model(&bet).
		Related(&bet.CompetitionCompetitor, "CompetitionCompetitor").
		Related(&bet.Better, "Better")

	db.Model(&bet.CompetitionCompetitor).
		Related(&bet.CompetitionCompetitor.Competition, "Competition").
		Related(&bet.CompetitionCompetitor.Competitor, "Competitor")

	spew.Dump(bet)
}
