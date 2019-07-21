package betting

import (
	"context"

	"github.com/doug-martin/goqu/v7"
	"github.com/pkg/errors"

	"github.com/bombsimon/team-betting/pkg"
	"github.com/bombsimon/team-betting/pkg/database"
)

// Service represents a service and implementation of the team betting
// interface.
type Service struct {
	DB *pkg.Database
}

// AddCompetition will add a new competition.
func (s *Service) AddCompetition(ctx context.Context, competition *pkg.Competition) (*pkg.Competition, error) {
	if err := competition.Validate(); err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	result, err := s.DB.Gq.From(pkg.CompetitionTable).
		Insert(
			goqu.Record{
				"image":       competition.Image,
				"name":        competition.Name,
				"description": competition.Description,
			},
		).Exec()

	if err != nil {
		return nil, errors.Wrap(err, "could not create competition")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "could not get competition ID")
	}

	return s.GetCompetition(ctx, int(id))
}

// AddCompetitor will add a new competitor that may be bound to a competition.
func (s *Service) AddCompetitor(ctx context.Context, competitor *pkg.Competitor, bindToCompetitionID *int) (*pkg.Competitor, error) {
	if err := competitor.Validate(); err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	result, err := s.DB.Gq.From(pkg.CompetitorTable).
		Insert(
			goqu.Record{
				"image":       competitor.Image,
				"name":        competitor.Name,
				"description": competitor.Description,
			},
		).Exec()

	if err != nil {
		return nil, errors.Wrap(err, "could not create competitor")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "could not get competitor ID")
	}

	if bindToCompetitionID != nil {
		if err := s.AddCompetitorToCompetition(ctx, int(id), *bindToCompetitionID); err != nil {
			return nil, errors.Wrap(err, "could not bind added competitor to competition")
		}
	}

	return s.GetCompetitor(ctx, int(id))
}

// AddBetter will add a new better that may place bets.
func (s *Service) AddBetter(ctx context.Context, better *pkg.Better) (*pkg.Better, error) {
	if err := better.Validate(); err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	result, err := s.DB.Gq.From(pkg.BetterTable).
		Insert(
			goqu.Record{
				"name":  better.Name,
				"email": better.Email,
				"image": better.Image,
			},
		).Exec()

	if err != nil {
		if database.ErrType(err) == database.ErrDuplicateKey {
			return nil, errors.New("a user with that email already exist!")
		}

		return nil, errors.Wrap(err, "could not create better")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "could not get better ID")
	}

	return s.GetBetter(ctx, int(id))
}

// AddBet will add a bet for a better to a competitor in a competition.
func (s *Service) AddBet(ctx context.Context, bet *pkg.Bet) error {
	_, err := s.DB.Gq.From(pkg.BetTable).
		Insert(
			goqu.Record{
				"id_better":                 bet.BetterID,
				"id_competition_competitor": bet.CompetitionCompetitorID,
				"placing":                   bet.Placing,
				"note":                      bet.Note,
			},
		).Exec()

	if err != nil {
		return errors.Wrap(err, "could not create bet")
	}

	return nil
}

// AddCompetitorToCompetition will add a competitor to a specific competition.
func (s *Service) AddCompetitorToCompetition(ctx context.Context, competitorID, competitionID int) error {
	if competitorID < 1 {
		return errors.New("invalid competitor")
	}

	if competitionID < 1 {
		return errors.New("invalid competition")
	}

	_, err := s.DB.Gq.From(pkg.CompetitionCompetitorTable).
		Insert(
			goqu.Record{
				"id_competitor":  competitorID,
				"id_competition": competitionID,
			},
		).Exec()

	if err != nil {
		if database.ErrType(err) == database.ErrForeignKeyConstraint {
			return errors.Wrap(pkg.ErrBadRequest, "invalid competitor/competition combination")
		}

		return errors.Wrap(err, "could not add competitor to competition")
	}

	return nil
}

// GetCompetition will return a competition based on a competition ID.
func (s *Service) GetCompetition(ctx context.Context, competitionID int) (*pkg.Competition, error) {
	c, err := s.GetCompetitions(ctx, []int{competitionID})
	if err != nil {
		return nil, err
	}

	if len(c) != 1 {
		return nil, errors.Wrap(pkg.ErrNotFound, "no competition found")
	}

	return c[0], nil
}

// GetCompetitions will return a list of competition based on competition IDs.
func (s *Service) GetCompetitions(ctx context.Context, competitionIDs []int) ([]*pkg.Competition, error) {
	var competitions []*pkg.Competition

	err := s.DB.Gq.From(pkg.CompetitionTable).
		Where(
			goqu.Ex{"id": competitionIDs},
		).ScanStructs(&competitions)

	if err != nil {
		return nil, errors.Wrap(err, "could not get competition")
	}

	return competitions, nil
}

// GetCompetitor will return a competitor based on a competitor ID.
func (s *Service) GetCompetitor(ctx context.Context, competitorID int) (*pkg.Competitor, error) {
	c, err := s.GetCompetitors(ctx, []int{competitorID})
	if err != nil {
		return nil, err
	}

	if len(c) != 1 {
		return nil, errors.Wrap(pkg.ErrNotFound, "no competitor found")
	}

	return c[0], nil
}

// GetCompetitors will return a list of competitor based on competitor IDs.
func (s *Service) GetCompetitors(ctx context.Context, competitorIDs []int) ([]*pkg.Competitor, error) {
	var competitors []*pkg.Competitor

	err := s.DB.Gq.From(pkg.CompetitorTable).
		Where(
			goqu.Ex{"id": competitorIDs},
		).ScanStructs(&competitors)

	if err != nil {
		return nil, errors.Wrap(err, "could not get competition")
	}

	return competitors, nil
}

// GetBetter will return a better based on a better ID.
func (s *Service) GetBetter(ctx context.Context, betterID int) (*pkg.Better, error) {
	b, err := s.GetBetters(ctx, []int{betterID})
	if err != nil {
		return nil, err
	}

	if len(b) != 1 {
		return nil, errors.Wrap(pkg.ErrNotFound, "no better found")
	}

	return b[0], nil
}

// GetBetters will return a list of better based on better IDs.
func (s *Service) GetBetters(ctx context.Context, betterIDs []int) ([]*pkg.Better, error) {
	var betters []*pkg.Better

	err := s.DB.Gq.From(pkg.BetterTable).
		Where(
			goqu.Ex{"id": betterIDs},
		).ScanStructs(&betters)

	if err != nil {
		return nil, errors.Wrap(err, "could not get betters")
	}

	return betters, nil
}

// GetCompetitorsForCompetition returns a slice with all competitors for a given
// competition.
func (s *Service) GetCompetitorsForCompetition(ctx context.Context, competitionID int) ([]*pkg.Competitor, error) {
	var competitors []*pkg.Competitor

	err := s.DB.Gq.From(pkg.CompetitorTable).
		Select("competitor.*").
		LeftJoin(
			goqu.I(pkg.CompetitionCompetitorTable),
			goqu.On(goqu.Ex{"id_competitor": "competitor.id"}),
		).
		Where(
			goqu.Ex{"id_competition": competitionID},
		).
		ScanStructs(&competitors)

	if err != nil {
		return nil, errors.Wrap(err, "could not get competitors")
	}

	return competitors, nil
}

func (s *Service) GetBetsForCompetition(ctx context.Context, competitionID int) ([]*pkg.Bet, error) {
	var (
		betRows []struct {
			*pkg.Bet
			*pkg.CompetitionCompetitor
		}
		bets               []*pkg.Bet
		uniqueBetters      map[int]struct{}
		uniqueCompetitions map[int]struct{}
		uniqueCompetitors  map[int]struct{}
		idToBetter         map[int]*pkg.Better
		idToCompetition    map[int]*pkg.Competition
		idToCompetitor     map[int]*pkg.Competitor
	)

	err := s.DB.Gq.From(pkg.BetTable).
		Join(
			goqu.I(pkg.CompetitionCompetitorTable),
			goqu.On(goqu.Ex{"bet.id_competition_competitor": "competition_competitor.id"}),
		).
		Where(
			goqu.Ex{"competition_competitor.id_competition": competitionID},
		).
		ScanStructs(&betRows)

	if err != nil {
		return nil, errors.Wrap(err, "could not get bets")
	}

	for _, bet := range betRows {
		uniqueBetters[bet.BetterID] = struct{}{}
		uniqueCompetitions[bet.CompetitionID] = struct{}{}
		uniqueCompetitors[bet.CompetitorID] = struct{}{}
	}

	mapToList := func(m map[int]struct{}) []int {
		var ids []int

		for id := range m {
			ids = append(ids, id)
		}

		return ids
	}

	betters, err := s.GetBetters(ctx, mapToList(uniqueBetters))
	if err != nil {
		return nil, err
	}

	competitions, err := s.GetCompetitions(ctx, mapToList(uniqueCompetitions))
	if err != nil {
		return nil, err
	}

	competitors, err := s.GetCompetitors(ctx, mapToList(uniqueCompetitors))
	if err != nil {
		return nil, err
	}

	for _, better := range betters {
		idToBetter[better.ID] = better
	}

	for _, competition := range competitions {
		idToCompetition[competition.ID] = competition
	}

	for _, competitor := range competitors {
		idToCompetitor[competitor.ID] = competitor
	}

	for _, b := range betRows {
		bet := b.Bet
		bet.Better = idToBetter[b.BetterID]
		bet.Competition = idToCompetition[b.CompetitionID]
		bet.Competitor = idToCompetitor[b.CompetitorID]

		bets = append(bets, bet)
	}

	return bets, nil
}
