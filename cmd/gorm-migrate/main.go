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
	db := database.New(os.Getenv("DB_DSN")).Gorm

	defer db.Close()

	// Don't use pluralis names for table names.
	db.SingularTable(true)

	// Add all auto migrations.
	db.AutoMigrate(&pkg.Better{})

	// Add foreign keys keys, not done bu auto migrations
	// https://github.com/jinzhu/gorm/issues/450
	db.AutoMigrate(&pkg.Competitor{}).
		AddForeignKey("created_by_id", "better(id)", "CASCADE", "CASCADE")

	db.AutoMigrate(&pkg.Competition{}).
		AddForeignKey("created_by_id", "better(id)", "CASCADE", "CASCADE")

	db.AutoMigrate(&pkg.Result{}).
		AddForeignKey("competition_id", "competition(id)", "CASCADE", "CASCADE").
		AddForeignKey("competitor_id", "competitor(id)", "CASCADE", "CASCADE")

	db.AutoMigrate(&pkg.Bet{}).
		AddForeignKey("better_id", "better(id)", "CASCADE", "CASCADE").
		AddForeignKey("competition_id", "competition(id)", "CASCADE", "CASCADE").
		AddForeignKey("competitor_id", "competitor(id)", "CASCADE", "CASCADE")

	if os.Getenv("ADD_DATA") != "" {
		testAddData(db)
	}

	if os.Getenv("GET_DATA") != "" {
		testGetData(db)
	}
}

func testAddData(db *gorm.DB) {
	betters := []*pkg.Better{
		{Name: "Testy Testsson", Email: "testy@testsson.se", Confirmed: true},
		{Name: "Another Tester", Email: "testy@anotherone.se"},
	}

	for _, b := range betters {
		db.Create(b)
	}

	competition := &pkg.Competition{
		Name:        "Eurovision Song Contest 2020",
		Description: null.StringFrom("The one that started it all"),
		CreatedBy:   betters[0],
		MinScore:    0,
		MaxScore:    10,
	}

	competitors := []*pkg.Competitor{
		{
			Name:        "Sweden",
			Description: null.StringFrom("Take me to your heaven"),
			CreatedBy:   betters[0],
		},
		{
			Name:        "Norway",
			Description: null.StringFrom("Norwegian song"),
			CreatedBy:   betters[0],
		},
	}

	db.Create(competition)

	for _, c := range competitors {
		c.Competitions = []*pkg.Competition{competition}

		db.Create(c)
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
