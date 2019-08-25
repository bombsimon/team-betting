package betting

import (
	"context"

	"github.com/doug-martin/goqu"
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
		if err := s.AddCompetitorToCompetition(ctx, cleaned.ID, *bindToCompetitionID); err != nil {
			return nil, errors.Wrap(err, "could not bind added competitor to competition")
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
func (s *Service) AddBet(ctx context.Context, bet *pkg.Bet) error {
	competition, err := s.GetCompetition(ctx, bet.CompetitionID)
	if err != nil {
		return err
	}

	if err := bet.Validate(competition.MinScore, competition.MaxScore); err != nil {
		return errors.Wrap(err, "bad request")
	}

	var cc pkg.CompetitionCompetitor
	result := s.DB.Gorm.
		Where("id_competition = ? AND id_competitor = ?", bet.CompetitionID, bet.CompetitorID).
		First(&cc)

	if result.RecordNotFound() {
		return errors.Wrap(pkg.ErrBadRequest, "invalid competition/competitor combination")
	}

	where := pkg.Bet{
		BetterID:                bet.BetterID,
		CompetitionCompetitorID: cc.ID,
	}

	cleaned := pkg.Bet{
		Placing: bet.Placing,
		Score:   bet.Score,
		Note:    bet.Note,
	}

	result = s.DB.Gorm.Where(where).
		Assign(pkg.Bet{
			Placing: cleaned.Placing,
			Score:   cleaned.Score,
			Note:    cleaned.Note,
		}).
		FirstOrCreate(&cleaned)

	if result.Error != nil {
		return errors.Wrap(result.Error, "could not create or update bet")
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

	result := s.DB.Gorm.Create(&pkg.CompetitionCompetitor{
		CompetitionID: competitionID,
		CompetitorID:  competitorID,
	})

	if err := result.Error; err != nil {
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

	if err := s.DB.Gorm.Where(competitionIDs).Find(&competitions).Error; err != nil {
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

	if err := s.DB.Gorm.Where(competitorIDs).Find(&competitors).Error; err != nil {
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

	if err := s.DB.Gorm.Where(betterIDs).Find(&betters).Error; err != nil {
		return nil, errors.Wrap(err, "could not get betters")
	}

	return betters, nil
}

// GetCompetitorsForCompetition returns a slice with all competitors for a given
// competition.
func (s *Service) GetCompetitorsForCompetition(ctx context.Context, competitionID int) ([]*pkg.Competitor, error) {
	var (
		competition pkg.Competition
	)

	if s.DB.Gorm.Preload("Competitors").Where("id = ?", competitionID).First(&competition).RecordNotFound() {
		return nil, errors.Wrap(pkg.ErrNotFound, "competition not found")
	}

	return competition.Competitors, nil
}

func (s *Service) GetBetsForCompetition(ctx context.Context, competitionID int) ([]*pkg.Bet, error) {
	var (
		bets               []*pkg.Bet
		uniqueBetters      = map[int]struct{}{}
		uniqueCompetitions = map[int]struct{}{}
		uniqueCompetitors  = map[int]struct{}{}
		idToBetter         = map[int]*pkg.Better{}
		idToCompetition    = map[int]*pkg.Competition{}
		idToCompetitor     = map[int]*pkg.Competitor{}
	)

	err := s.DB.Gq.From(pkg.BetTable).
		Select(
			goqu.I("bet.id"),
			goqu.I("bet.created_at"),
			goqu.I("bet.updated_at"),
			goqu.I("bet.id_better"),
			goqu.I("competition_competitor.id_competition"),
			goqu.I("competition_competitor.id_competitor"),
			goqu.I("bet.score"),
			goqu.I("bet.placing"),
			goqu.I("bet.note"),
		).
		Join(
			goqu.I(pkg.CompetitionCompetitorTable),
			goqu.On(goqu.Ex{"bet.id_competition_competitor": goqu.I("competition_competitor.id")}),
		).
		Where(
			goqu.Ex{"competition_competitor.id_competition": competitionID},
		).
		ScanStructs(&bets)

	if err != nil {
		return nil, errors.Wrap(err, "could not get bets")
	}

	if len(bets) < 1 {
		return bets, nil
	}

	for _, bet := range bets {
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

	for i := range bets {
		bets[i].Better = *idToBetter[bets[i].BetterID]
		bets[i].Competition = *idToCompetition[bets[i].CompetitionID]
		bets[i].Competitor = *idToCompetitor[bets[i].CompetitorID]
	}

	return bets, nil
}
