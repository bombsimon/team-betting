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

	s := &Service{
		DB: db,
	}

	// Ensure there's always at least one better.
	_, err := s.AddBetter(context.Background(), &pkg.Better{
		Name:  "Unittest better",
		Email: "user@iamveryunique.se",
	})

	require.NoError(t, err)

	return s
}

func (s *Service) anyBetter() *pkg.Better {
	var b pkg.Better

	if err := s.DB.Gorm.First(&b).Error; err != nil {
		panic(err)
	}

	return &b
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
				CreatedByID: s.anyBetter().ID,
				Name:        "Unittest Challenge",
				Description: null.StringFrom("A test made for unit testing"),
			},
		},
		{
			description: "successful create with emojis",
			competition: &pkg.Competition{
				CreatedByID: s.anyBetter().ID,
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
				CreatedByID: s.anyBetter().ID,
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

func TestService_AddBetter(t *testing.T) {
	s := setupService(t)

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
			errContains: "a user with that email already exist",
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
		CreatedByID: s.anyBetter().ID,
		Name:        "Unittest competition",
	})

	require.NoError(t, err)

	competitionWithoutCompetitors, err := s.AddCompetition(context.Background(), &pkg.Competition{
		CreatedByID: s.anyBetter().ID,
		Name:        "Sad unittest competition",
	})

	require.NoError(t, err)

	for i := range make([]int, 3) {
		c, err := s.AddCompetitor(context.Background(), &pkg.Competitor{
			CreatedByID: s.anyBetter().ID,
			Name:        fmt.Sprintf("Unittest competitor %d", i+1),
		}, &competition.ID)

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
			description: "competitor does not compete under competition",
			bet: &pkg.Bet{
				BetterID:      betterIDs[0],
				CompetitionID: competitionWithoutCompetitors.ID,
				CompetitorID:  competitorIDs[0],
				Note:          null.StringFrom("This shouldn't work"),
			},
			errContains: "invalid competition/competitor combination",
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
		{
			description: "successful second bet",
			bet: &pkg.Bet{
				BetterID:      betterIDs[1],
				CompetitionID: competition.ID,
				CompetitorID:  competitorIDs[0],
				Placing:       null.IntFrom(3),
				Score:         null.IntFrom(6),
				Note:          null.StringFrom("Want more bets"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			_, err := s.AddBet(ctx, tc.bet)

			if tc.errContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)

				return
			}

			// TODO: Check the data is correct
			require.NoError(t, err)
		})
	}

	bets, err := s.GetBetsForCompetition(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, bets)

	assert.Len(t, bets, 2)
}

func TestPreloadMany2Many(t *testing.T) {
	s := setupService(t)

	competition, err := s.AddCompetition(context.Background(), &pkg.Competition{
		CreatedByID: s.anyBetter().ID,
		Name:        "Unittest competition",
	})

	require.NoError(t, err)

	for i := range make([]int, 3) {
		_, err := s.AddCompetitor(context.Background(), &pkg.Competitor{
			CreatedByID: s.anyBetter().ID,
			Name:        fmt.Sprintf("Unittest competitor %d", i+1),
		}, &competition.ID)

		require.NoError(t, err)
	}

	r, err := s.GetCompetitorsForCompetition(context.Background(), 1)

	require.Nil(t, err)
	assert.Len(t, r, 3)
}

func TestGetCompetitionMetrics(t *testing.T) {
	var (
		s             = setupService(t)
		competitorIDs []int
		betterIDs     []int
	)

	competition, err := s.AddCompetition(context.Background(), &pkg.Competition{
		Name:        "Unittest competition",
		CreatedByID: s.anyBetter().ID,
	})

	require.NoError(t, err)

	for i := range make([]int, 3) {
		c, err := s.AddCompetitor(context.Background(), &pkg.Competitor{
			CreatedByID: s.anyBetter().ID,
			Name:        fmt.Sprintf("Unittest competitor %d", i+1),
		}, &competition.ID)

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

	for i := range make([]int, 3) {
		_, err := s.AddBet(context.Background(), &pkg.Bet{
			BetterID:      betterIDs[0],
			CompetitionID: competition.ID,
			CompetitorID:  competitorIDs[i],
			Placing:       null.IntFrom(2),
			Score:         null.IntFrom(int64(0 + i)),
			Note:          null.StringFrom("Want more bets"),
		})

		require.NoError(t, err)

		_, err = s.AddBet(context.Background(), &pkg.Bet{
			BetterID:      betterIDs[1],
			CompetitionID: competition.ID,
			CompetitorID:  competitorIDs[i],
			Placing:       null.IntFrom(int64(i * i)),
			Score:         null.IntFrom(int64(10 - i)),
			Note:          null.StringFrom("Another one with some notes here"),
		})

		require.NoError(t, err)
	}

	m, err := s.GetCompetitionMetrics(context.Background(), competition.ID)

	require.NoError(t, err)

	assert.Equal(t, "Unittest better 2", m.HighestAverageBetter.Who.Name)
	assert.Equal(t, m.HighestAverageBetter.Value, float64(9))
	assert.Equal(t, m.LowestAverageBetter.Who.Name, "Unittest better 1")
	assert.Equal(t, m.LowestAverageBetter.Value, float64(1))
	assert.Equal(t, m.MostTopScores.Who.Name, "Unittest better 2")
	assert.Equal(t, m.MostTopScores.Value, 1)
	assert.Equal(t, m.MostBottomScores.Who.Name, "Unittest better 1")
	assert.Equal(t, m.MostBottomScores.Value, 1)
	assert.Equal(t, m.LongestNote.Who.Name, "Unittest better 2")
	assert.Equal(t, m.LongestNote.Value, "Another one with some notes here")
	assert.Equal(t, m.ShortestNote.Who.Name, "Unittest better 1")
	assert.Equal(t, m.ShortestNote.Value, "Want more bets")
	assert.Equal(t, m.NumberOfBottomScores, 1)
	assert.Equal(t, m.NumberOfTopScores, 1)
	assert.Equal(t, m.GroupAverageScore, float64(5))
}
