package betting

import (
	"context"

	"github.com/bombsimon/team-betting/pkg"
	"github.com/doug-martin/goqu"
)

// Service represents a service and implementation of the team betting
// interface.
type Service struct {
	DB *pkg.Database
}

// GetCompetition will return a competition based on a competition ID.
func (s *Service) GetCompetition(ctx context.Context, competitionID int) (*pkg.Competition, error) {
	var competition *pkg.Competition

	found, err := s.DB.Gq.From(pkg.CompetitionTable).
		Where(
			goqu.Ex{"id": competitionID},
		).ScanStruct()

	if !found {
		return nil, errors.Wrap(pkg.ErrNotFound, "no competition found")
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not get competition")
	}

	return competition, nil
}

// AddCompetition will add a new competition.
func (s *Service) AddCompetition(ctx context.Context, competition *pkg.Competition) (*pkg.Competition, error) {
	if !competition.Valid() {
		return nil, errors.Wrap(err, "bad request")
	}

	result, err := s.DB.Gq.From(pkg.CompetitionTable).
		Insert(
			goqu.Record{
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

	return s.GetCompetition(id)
}

// GetCompetitor will return a competitor based on a competitor ID.
func (s *Service) GetCompetitor(ctx context.Context, competitorID int) (*pkg.COMpetitor, error) {
	return nil, nil
}
