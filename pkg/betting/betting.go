package betting

import (
	"context"

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

	cleaned := pkg.Competition{
		Name:        competition.Name,
		Description: competition.Description,
		Image:       competition.Image,
		MinScore:    competition.MinScore,
		MaxScore:    competition.MaxScore,
	}

	if cleaned.MaxScore == 0 {
		cleaned.MaxScore = 10
	}

	if err := s.DB.Gorm.Save(&cleaned).Error; err != nil {
		return nil, errors.Wrap(err, "could not create competition")
	}

	return &cleaned, nil
}

// AddCompetitor will add a new competitor that may be bound to a competition.
func (s *Service) AddCompetitor(ctx context.Context, competitor *pkg.Competitor, bindToCompetitionID *int) (*pkg.Competitor, error) {
	if err := competitor.Validate(); err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	cleaned := pkg.Competitor{
		Name:        competitor.Name,
		Description: competitor.Description,
		Image:       competitor.Image,
	}

	if err := s.DB.Gorm.Save(&cleaned).Error; err != nil {
		return nil, errors.Wrap(err, "could not create competitor")
	}

	if bindToCompetitionID != nil {
		competition, err := s.GetCompetition(ctx, *bindToCompetitionID)
		if err != nil {
			return nil, errors.Wrap(err, "could not find competition to bind to competitor to")
		}

		if err := s.DB.Gorm.Model(competition).Association("Competitors").Append(&cleaned).Error; err != nil {
			return nil, errors.Wrap(err, "could not link competitor to competition")
		}
	}

	return &cleaned, nil
}

// AddBetter will add a new better that may place bets.
func (s *Service) AddBetter(ctx context.Context, better *pkg.Better) (*pkg.Better, error) {
	if err := better.Validate(); err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	cleaned := pkg.Better{
		Name:  better.Name,
		Email: better.Email,
		Image: better.Image,
	}

	if err := s.DB.Gorm.Save(&cleaned).Error; err != nil {
		if database.ErrType(err) == database.ErrDuplicateKey {
			return nil, errors.New("a user with that email already exist")
		}

		return nil, errors.Wrap(err, "could not create competitor")
	}

	return &cleaned, nil
}

// AddBet will add a bet for a better to a competitor in a competition.
func (s *Service) AddBet(ctx context.Context, bet *pkg.Bet) (*pkg.Bet, error) {
	competition, err := s.GetCompetition(ctx, bet.CompetitionID)
	if err != nil {
		return nil, err
	}

	if err := bet.Validate(competition.MinScore, competition.MaxScore); err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	// Ensure the competitor actually competes in the competition.
	r := s.DB.Gorm.
		Model(&pkg.Competition{ID: bet.CompetitionID}).
		Association("Competitors").
		Find(&pkg.Competitor{ID: bet.CompetitorID})

	if r.Error != nil {
		return nil, errors.Wrap(r.Error, "invalid competition/competitor combination")
	}

	where := pkg.Bet{
		BetterID:      bet.BetterID,
		CompetitionID: bet.CompetitionID,
		CompetitorID:  bet.CompetitorID,
	}

	cleaned := pkg.Bet{
		Placing: bet.Placing,
		Score:   bet.Score,
		Note:    bet.Note,
	}

	result := s.DB.Gorm.Where(where).
		Assign(pkg.Bet{
			Placing: cleaned.Placing,
			Score:   cleaned.Score,
			Note:    cleaned.Note,
		}).
		FirstOrCreate(&cleaned)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "could not create or update bet")
	}

	return &cleaned, nil
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

	q := s.DB.Gorm
	if len(competitionIDs) > 0 {
		q = q.Where(competitionIDs)
	}

	if err := q.Find(&competitions).Error; err != nil {
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

	q := s.DB.Gorm
	if len(competitorIDs) > 0 {
		q = q.Where(competitorIDs)
	}

	if err := q.Find(&competitors).Error; err != nil {
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

	q := s.DB.Gorm
	if len(betterIDs) > 0 {
		q = q.Where(betterIDs)
	}

	if err := q.Find(&betters).Error; err != nil {
		return nil, errors.Wrap(err, "could not get betters")
	}

	return betters, nil
}

// GetBet will return a bet based on a bet ID.
func (s *Service) GetBet(ctx context.Context, betID int) (*pkg.Bet, error) {
	b, err := s.GetBets(ctx, []int{betID})
	if err != nil {
		return nil, err
	}

	if len(b) != 1 {
		return nil, errors.Wrap(pkg.ErrNotFound, "no bet found")
	}

	return b[0], nil
}

// GetBets will return a list of bets based on bet IDs.
func (s *Service) GetBets(ctx context.Context, betIDs []int) ([]*pkg.Bet, error) {
	var bets []*pkg.Bet

	q := s.DB.Gorm
	if len(betIDs) > 0 {
		q = q.Where(betIDs)
	}

	if err := q.Find(&bets).Error; err != nil {
		return nil, errors.Wrap(err, "could not get bets")
	}

	return bets, nil
}

// GetCompetitorsForCompetition returns a slice with all competitors for a given
// competition.
func (s *Service) GetCompetitorsForCompetition(ctx context.Context, competitionID int) ([]*pkg.Competitor, error) {
	var competition pkg.Competition

	if s.DB.Gorm.Preload("Competitors").Where("id = ?", competitionID).First(&competition).RecordNotFound() {
		return nil, errors.Wrap(pkg.ErrNotFound, "competition not found")
	}

	return competition.Competitors, nil
}

// GetBetsForCompetition returns all bets for a given competition.
func (s *Service) GetBetsForCompetition(ctx context.Context, competitionID int) ([]*pkg.Bet, error) {
	var bets []*pkg.Bet

	r := s.DB.Gorm.Preload("Competition").Preload("Competitor").
		Where("competition_id = ?", competitionID).
		Find(&bets)

	if err := r.Error; err != nil {
		return nil, errors.Wrap(err, "could not get bets")
	}

	return bets, nil
}
