package main

import (
	"fmt"
	"os"

	"github.com/bombsimon/team-betting/pkg"
	"github.com/bombsimon/team-betting/pkg/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/guregu/null"
	"github.com/jinzhu/gorm"
)

func main() {
	db := database.New("").Gorm

	defer db.Close()

	// Don't use pluralis names for table names.
	db.SingularTable(true)

	// Add all auto migrations.
	db.AutoMigrate(&pkg.Competition{})
	db.AutoMigrate(&pkg.Competitor{})
	db.AutoMigrate(&pkg.Better{})

	// Add foreign keys keys, not done bu auto migrations
	// https://github.com/jinzhu/gorm/issues/450
	db.AutoMigrate(&pkg.Bet{}).
		AddForeignKey("better_id", "better(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("competition_id", "competition(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("competitor_id", "competitor(id)", "RESTRICT", "RESTRICT")

	if os.Getenv("ADD_DATA") != "" {
		testAddData(db)
	}

	if os.Getenv("GET_DATA") != "" {
		testGetData(db)
	}
}

func testAddData(db *gorm.DB) {
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
		c.Competitions = []*pkg.Competition{competition}

		db.Create(c)
	}

	for _, b := range betters {
		db.Create(b)
	}

	db.Create(&pkg.Bet{
		Better:      betters[0],
		Competitor:  competitors[0],
		Competition: competition,
		Score:       null.IntFrom(6),
	})
}

func testGetData(db *gorm.DB) {
	// Fetch the first bet
	var bet pkg.Bet
	db.Preload("Competition").Preload("Competitor").Preload("Better").First(&bet)

	fmt.Printf("%-20s %s\n", "Created", bet.CreatedAt)
	fmt.Printf("%-20s %s\n", "Better", bet.Better.Name)
	fmt.Printf("%-20s %s\n", "Competition", bet.Competition.Name)
	fmt.Printf("%-20s %s\n", "Competitor", bet.Competitor.Name)
	fmt.Printf("%-20s %d\n", "Score", bet.Score.ValueOrZero())
	fmt.Printf("%-20s %d\n", "Placing", bet.Placing.ValueOrZero())
	fmt.Printf("%-20s %s\n", "Note", bet.Note.ValueOrZero())
}
