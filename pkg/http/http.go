package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bombsimon/team-betting/pkg"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gopkg.in/olahol/melody.v1"
)

// Service represents the HTTP service serving the team betting.
type Service struct {
	Betting pkg.BettingService
	WS      *melody.Melody
}

// SignInEmail will sign in from mail.
func (s *Service) SignInEmail(c *gin.Context) {
	var (
		email = c.Query("email")
		hash  = c.Query("hash")
		data  map[string]string
	)

	jwtString, err := s.Betting.SignInFromEmail(context.Background(), email, hash)
	if err == nil {
		data = map[string]string{
			"jwt": jwtString,
		}
	}

	s.HandleResponse(c, nil, data, err)
}

// GetCompetitions returns all competitions.
func (s *Service) GetCompetitions(c *gin.Context) {
	data, err := s.Betting.GetCompetitions(context.Background(), []int{})

	s.HandleResponse(c, nil, data, err)
}

// GetCompetition returns a competition (if it exists).
func (s *Service) GetCompetition(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := s.Betting.GetCompetition(context.Background(), id)

	s.HandleResponse(c, nil, data, err)
}

// AddCompetition adds a competition.
func (s *Service) AddCompetition(c *gin.Context) {
	var competition pkg.Competition

	if err := c.ShouldBindJSON(&competition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := isCreator(c, competition.CreatedByID); err != nil {
		s.HandleResponse(c, nil, nil, err)
		return
	}

	data, err := s.Betting.AddCompetition(context.Background(), &competition)

	s.HandleResponse(c, nil, data, err)
}

// DeleteCompetition returns a competition (if it exists).
func (s *Service) DeleteCompetition(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := s.Betting.DeleteCompetition(context.Background(), id)

	s.HandleResponse(c, nil, nil, err)
}

// GetCompetitors returns all competitions.
func (s *Service) GetCompetitors(c *gin.Context) {
	data, err := s.Betting.GetCompetitors(context.Background(), []int{})

	s.HandleResponse(c, nil, data, err)
}

// GetCompetitor returns a competition (if it exists).
func (s *Service) GetCompetitor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := s.Betting.GetCompetitor(context.Background(), id)

	s.HandleResponse(c, nil, data, err)
}

// AddCompetitor adds a competitor.
func (s *Service) AddCompetitor(c *gin.Context) {
	var in struct {
		pkg.Competitor
		CompetitionID *int `json:"competition_id"`
	}

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := s.Betting.AddCompetitor(context.Background(), &in.Competitor, in.CompetitionID)

	s.HandleResponse(c, nil, data, err)
}

// DeleteCompetitor returns a competitor (if it exists).
func (s *Service) DeleteCompetitor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := s.Betting.DeleteCompetitor(context.Background(), id)

	s.HandleResponse(c, nil, nil, err)
}

// GetBetters returns all competitions.
func (s *Service) GetBetters(c *gin.Context) {
	data, err := s.Betting.GetBetters(context.Background(), []int{})

	s.HandleResponse(c, nil, data, err)
}

// GetBetter returns a competition (if it exists).
func (s *Service) GetBetter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := s.Betting.GetBetter(context.Background(), id)

	s.HandleResponse(c, nil, data, err)
}

// AddBetter adds a better.
func (s *Service) AddBetter(c *gin.Context) {
	var better pkg.Better

	if err := c.ShouldBindJSON(&better); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := s.Betting.AddBetter(context.Background(), &better)

	s.HandleResponse(c, nil, data, err)
}

// DeleteBetter returns a better (if it exists).
func (s *Service) DeleteBetter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := s.Betting.DeleteBetter(context.Background(), id)

	s.HandleResponse(c, nil, nil, err)
}

// GetBets returns all competitions.
func (s *Service) GetBets(c *gin.Context) {
	data, err := s.Betting.GetBets(context.Background(), []int{})

	s.HandleResponse(c, nil, data, err)
}

// GetBet returns a competition (if it exists).
func (s *Service) GetBet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := s.Betting.GetBet(context.Background(), id)

	s.HandleResponse(c, nil, data, err)
}

// AddBet adds a bet.
func (s *Service) AddBet(c *gin.Context) {
	var (
		bet pkg.Bet
		bc  []byte
	)

	if err := c.ShouldBindJSON(&bet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := s.Betting.AddBet(context.Background(), &bet)
	if data != nil {
		bc, _ = json.MarshalIndent(&data, "", "  ")
	}

	s.HandleResponse(c, bc, data, err)
}

// DeleteBet returns a bet (if it exists).
func (s *Service) DeleteBet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := s.Betting.DeleteBet(context.Background(), id)

	s.HandleResponse(c, nil, nil, err)
}

func isCreator(c *gin.Context, creatorID int) error {
	return nil

	b, ok := c.Get("better")
	if !ok {
		return errors.Wrap(pkg.ErrBadRequest, "no or invalid authorization header")
	}

	better, ok := b.(*pkg.Better)
	if !ok {
		return errors.Wrap(pkg.ErrBadRequest, "authorization isn't a better")
	}

	if better.ID != creatorID {
		return errors.Wrap(pkg.ErrBadRequest, "creator is not signed in user")
	}

	return nil
}

// HandleResponse will respond according to the object and error passed.
func (s *Service) HandleResponse(c *gin.Context, broadcast []byte, response interface{}, err error) {
	if err != nil {
		var httpStatus = http.StatusInternalServerError

		if errors.Cause(err) == pkg.ErrNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}

		if errors.Cause(err) == pkg.ErrBadRequest {
			httpStatus = http.StatusBadRequest
		}

		if _, ok := errors.Cause(err).(validation.Errors); ok {
			httpStatus = http.StatusBadRequest
		}

		c.JSON(httpStatus, gin.H{"error": err.Error()})

		return
	}

	if broadcast != nil {
		_ = s.WS.Broadcast(broadcast)
	}

	c.JSON(http.StatusOK, response)
}
