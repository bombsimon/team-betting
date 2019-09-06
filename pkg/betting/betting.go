package betting

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/bombsimon/team-betting/pkg"
	"github.com/bombsimon/team-betting/pkg/database"
)

// Service represents a service and implementation of the team betting
// interface.
type Service struct {
	DB          *pkg.Database
	MailService pkg.MailService
}

// AddCompetition will add a new competition.
func (s *Service) AddCompetition(ctx context.Context, competition *pkg.Competition) (*pkg.Competition, error) {
	if err := competition.Validate(); err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	cleaned := pkg.Competition{
		CreatedByID: competition.CreatedByID,
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
		CreatedByID: competitor.CreatedByID,
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
func (s *Service) AddBetter(ctx context.Context, better *pkg.Better) (string, error) {
	if err := better.Validate(); err != nil {
		return "", errors.Wrap(err, "bad request")
	}

	cleaned := pkg.Better{
		Name:  better.Name,
		Email: better.Email,
		Image: better.Image,
	}

	if err := s.DB.Gorm.Save(&cleaned).Error; err != nil {
		if database.ErrType(err) == database.ErrDuplicateKey {
			return "", errors.New("a user with that email already exist")
		}

		return "", errors.Wrap(err, "could not create competitor")
	}

	signedToken, err := s.JWTForBetter(ctx, &cleaned)
	if err != nil {
		return "", errors.Wrap(err, "could not create JWT for new user")
	}

	return signedToken, nil
}

// AddBet will add a bet for a better to a competitor in a competition.
func (s *Service) AddBet(ctx context.Context, bet *pkg.Bet) (*pkg.Bet, error) {
	if err := bet.ValidateInit(); err != nil {
		return nil, errors.Wrap(err, "bad request")
	}

	competition, err := s.GetCompetition(ctx, bet.CompetitionID)
	if err != nil {
		return nil, errors.Wrap(err, "could not find competition to add bet to")
	}

	if err := bet.Validate(competition.MinScore, competition.MaxScore, len(competition.Competitors)); err != nil {
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

	// Ensure fields are non-nil when inflating.
	cleaned.Better = &pkg.Better{}
	cleaned.Competitor = &pkg.Competitor{}

	s.DB.Gorm.Model(&cleaned).
		Related(&cleaned.Better, "BetterID").
		Related(&cleaned.Competitor, "CompetitorID")

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

	err := q.
		Preload("CreatedBy").
		Preload("Competitors").
		Preload("Bets.Better").
		Preload("Bets.Competitor").
		Find(&competitions).
		Error

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

// DeleteCompetition will delete a competition
func (s *Service) DeleteCompetition(ctx context.Context, id int) error {
	c, err := s.GetCompetition(ctx, id)
	if err != nil {
		return err
	}

	if err := s.DB.Gorm.Delete(c).Error; err != nil {
		return errors.Wrap(err, "could not delete competition")
	}

	return nil
}

// DeleteCompetitor will delete a competitor.
func (s *Service) DeleteCompetitor(ctx context.Context, id int) error {
	c, err := s.GetCompetitor(ctx, id)
	if err != nil {
		return err
	}

	if err := s.DB.Gorm.Delete(c).Error; err != nil {
		return errors.Wrap(err, "could not delete competitor")
	}

	return nil
}

// DeleteBetter will delete a better.
func (s *Service) DeleteBetter(ctx context.Context, id int) error {
	b, err := s.GetBetter(ctx, id)
	if err != nil {
		return err
	}

	if err := s.DB.Gorm.Delete(b).Error; err != nil {
		return errors.Wrap(err, "could not delete better")
	}

	return nil
}

// DeleteBet will delete a bet.
func (s *Service) DeleteBet(ctx context.Context, id int) error {
	b, err := s.GetBet(ctx, id)
	if err != nil {
		return err
	}

	if err := s.DB.Gorm.Delete(b).Error; err != nil {
		return errors.Wrap(err, "could not delete bet")
	}

	return nil
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

// GetCreatedObjectsForBetter returns created competitions, competitors and
// bets created by a given user.
func (s *Service) GetCreatedObjectsForBetter(ctx context.Context, id int) ([]*pkg.Competition, []*pkg.Competitor, []*pkg.Bet, error) {
	var (
		competitions []*pkg.Competition
		competitors  []*pkg.Competitor
		bets         []*pkg.Bet
	)

	if err := s.DB.Gorm.Where("created_by_id = ?", id).Find(&competitions).Error; err != nil {
		return nil, nil, nil, errors.Wrap(err, "could not get competitions for better")
	}

	if err := s.DB.Gorm.Where("created_by_id = ?", id).Find(&competitors).Error; err != nil {
		return nil, nil, nil, errors.Wrap(err, "could not get competitors for better")
	}

	if err := s.DB.Gorm.Where("better_id = ?", id).Find(&bets).Error; err != nil {
		return nil, nil, nil, errors.Wrap(err, "could not get bets for better")
	}

	return competitions, competitors, bets, nil
}

// LockCompetition takes the final result and locks a competition.
func (s *Service) LockCompetition(ctx context.Context, id int, result []*pkg.Result) (*pkg.CompetitionMetrics, error) {
	c, err := s.GetCompetition(ctx, id)
	if err != nil {
		return nil, err
	}

	if c.Locked {
		return nil, errors.Wrap(pkg.ErrBadRequest, "competition already locked")
	}

	competitorIDsInCompetition := map[int]struct{}{}
	for _, v := range c.Competitors {
		competitorIDsInCompetition[v.ID] = struct{}{}
	}

	tx := s.DB.Gorm.Begin()

	c.Locked = true
	if err := tx.Save(c).Error; err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "could not lock competition")
	}

	for _, r := range result {
		r.CompetitionID = id

		if _, ok := competitorIDsInCompetition[r.CompetitorID]; !ok {
			tx.Rollback()
			return nil, errors.Wrap(pkg.ErrBadRequest, "competitor does not compete in competition")
		}

		if err := tx.Save(r).Error; err != nil {
			tx.Rollback()

			if strings.Contains(err.Error(), pkg.ResultPlacingKey) {
				return nil, errors.Wrap(pkg.ErrBadRequest, "only cone competitor can be placed at each position")
			}

			return nil, errors.Wrap(err, "could not save result")
		}
	}

	tx.Commit()

	return s.GetCompetitionMetrics(ctx, id)
}
