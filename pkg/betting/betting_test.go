package betting

import (
	"context"
	"fmt"
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bombsimon/team-betting/pkg"
	"github.com/bombsimon/team-betting/pkg/database"
)

func setupService(t *testing.T) *Service {
	db := database.New(
		"betting:betting@tcp(127.0.0.1:3306)/betting?parseTime=true&charset=utf8mb4&collation=utf8mb4_bin",
	)

	db.DB.Exec("SET FOREIGN_KEY_CHECKS=0")
	defer db.DB.Exec("SET FOREIGN_KEY_CHECKS=1")

	for _, tbl := range []string{
		pkg.BetTable,
		pkg.BetterTable,
		pkg.CompetitionCompetitorTable,
		pkg.CompetitorTable,
		pkg.CompetitionTable,
	} {
		_, err := db.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tbl))
		require.NoError(t, err)
	}

	return &Service{
		DB: db,
	}
}

func TestService_AddCompetition(t *testing.T) {
	s := setupService(t)

	cases := []struct {
		description string
		competition *pkg.Competition
		errContains string
	}{
		{
			description: "all missing data",
			competition: &pkg.Competition{},
			errContains: "bad request: name: cannot be blank.",
		},
		{
			description: "successful create",
			competition: &pkg.Competition{
				Name:        "Unittest Challenge",
				Description: null.StringFrom("A test made for unit testing"),
			},
		},
		{
			description: "successful create with emojis",
			competition: &pkg.Competition{
				Name:        "Just for cats 😸",
				Description: null.StringFrom("Cat game only! 😻"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			r, err := s.AddCompetition(ctx, tc.competition)

			if tc.errContains != "" {
				require.Error(t, err)
				require.Nil(t, r)
				assert.Contains(t, err.Error(), tc.errContains)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, r)
		})
	}
}

func TestService_AddCompetitor(t *testing.T) {
	s := setupService(t)

	cases := []struct {
		description string
		competitor  *pkg.Competitor
		errContains string
	}{
		{
			description: "all missing data",
			competitor:  &pkg.Competitor{},
			errContains: "bad request: name: cannot be blank.",
		},
		{
			description: "successful create",
			competitor: &pkg.Competitor{
				Name:        "Unittest Competitor",
				Description: null.StringFrom("Someone who can compete!"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			r, err := s.AddCompetitor(ctx, tc.competitor, nil)

			if tc.errContains != "" {
				require.Error(t, err)
				require.Nil(t, r)
				assert.Contains(t, err.Error(), tc.errContains)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, r)
		})
	}
}
func TestService_AddCompetitorToCompetition(t *testing.T) {
	s := setupService(t)

	competition, err := s.AddCompetition(context.Background(), &pkg.Competition{
		Name: "Unittest competition",
	})

	require.NoError(t, err)

	competitor, err := s.AddCompetitor(context.Background(), &pkg.Competitor{
		Name: "Unittest competitor",
	}, nil)

	require.NoError(t, err)

	cases := []struct {
		description   string
		competitorID  int
		competitionID int
		errContains   string
	}{
		{
			description: "all missing data",
			errContains: "invalid competitor",
		},
		{
			description:   "competition/competitor does not exist",
			competitorID:  1000,
			competitionID: 1000,
			errContains:   "invalid competitor/competition combination: bad request",
		},
		{
			description:   "successful bind",
			competitorID:  competitor.ID,
			competitionID: competition.ID,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			err := s.AddCompetitorToCompetition(ctx, tc.competitorID, tc.competitionID)

			if tc.errContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestService_AddBetter(t *testing.T) {
	s := setupService(t)

	competition, err := s.AddCompetition(context.Background(), &pkg.Competition{
		Name: "Unittest competition",
	})

	require.NoError(t, err)

	_, err = s.AddCompetitor(context.Background(), &pkg.Competitor{
		Name: "Unittest competitor",
	}, &competition.ID)

	require.NoError(t, err)

	cases := []struct {
		description string
		better      *pkg.Better
		errContains string
	}{
		{
			description: "all missing data",
			better:      &pkg.Better{},
			errContains: "bad request: email: cannot be blank; name: cannot be blank.",
		},
		{
			description: "invalid email",
			better: &pkg.Better{
				Name:  "Unittest better",
				Email: "zzz",
			},
			errContains: "bad request: email: must be a valid email address.",
		},
		{
			description: "successful create",
			better: &pkg.Better{
				Name:  "Unittest better",
				Email: "unit@test.se",
			},
		},
		{
			description: "cannot add user with same email",
			better: &pkg.Better{
				Name:  "Unittest better",
				Email: "unit@test.se",
			},
			errContains: "a user with that email already exist!",
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			b, err := s.AddBetter(ctx, tc.better)

			if tc.errContains != "" {
				require.Error(t, err)
				require.Nil(t, b)
				assert.Contains(t, err.Error(), tc.errContains)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, b)
		})
	}
}

func TestService_AddBet(t *testing.T) {
	var (
		s             = setupService(t)
		competitorIDs []int
		betterIDs     []int
	)

	competition, err := s.AddCompetition(context.Background(), &pkg.Competition{
		Name: "Unittest competition",
	})

	require.NoError(t, err)

	for i := range make([]int, 3) {
		c, err := s.AddCompetitor(context.Background(), &pkg.Competitor{
			Name: fmt.Sprintf("Unittest competitor %d", i+1)},
			&competition.ID,
		)

		require.NoError(t, err)
		require.NotNil(t, c)

		competitorIDs = append(competitorIDs, c.ID)

		b, err := s.AddBetter(context.Background(), &pkg.Better{
			Name:  fmt.Sprintf("Unittest better %d", i+1),
			Email: fmt.Sprintf("user%d@test.se", i+1),
		})

		require.NoError(t, err)
		require.NotNil(t, b)

		betterIDs = append(betterIDs, b.ID)
	}

	cases := []struct {
		description string
		bet         *pkg.Bet
		errContains string
	}{
		{
			description: "no competition found",
			bet:         &pkg.Bet{},
			errContains: "no competition found",
		},
		{
			description: "all missing data",
			bet: &pkg.Bet{
				CompetitionID: competition.ID,
			},
			errContains: "bad request: id_better: cannot be blank; id_competitor: cannot be blank.",
		},
		{
			description: "successful bet",
			bet: &pkg.Bet{
				BetterID:      betterIDs[0],
				CompetitionID: competition.ID,
				CompetitorID:  competitorIDs[0],
				Placing:       null.IntFrom(3),
				Score:         null.IntFrom(8),
				Note:          null.StringFrom("Lots of lots and stuff"),
			},
		},
		{
			description: "successful update of exiting bet",
			bet: &pkg.Bet{
				BetterID:      betterIDs[0],
				CompetitionID: competition.ID,
				CompetitorID:  competitorIDs[0],
				Placing:       null.IntFrom(5),
				Note:          null.StringFrom("Some other stuff"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			err := s.AddBet(ctx, tc.bet)

			if tc.errContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)

				return
			}

			require.NoError(t, err)
		})
	}
}
