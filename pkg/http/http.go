package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
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
	Logger  *log.Logger
}

// SendSignInEmail will send sign in email.
func (s *Service) SendSignInEmail(c *gin.Context) {
	var d struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.Betting.SendSignInEmail(context.Background(), d.Email)

	s.HandleResponse(c, nil, nil, err)
}

// SignInEmail will sign in from mail.
func (s *Service) SignInEmail(c *gin.Context) {
	var (
		sd   pkg.SignInData
		data map[string]string
	)

	if err := c.ShouldBindJSON(&sd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b64Data, err := base64.StdEncoding.DecodeString(sd.Encoding)
	if err != nil {
		s.HandleResponse(c, nil, nil, err)
		return
	}

	if err := json.Unmarshal(b64Data, &sd); err != nil {
		s.HandleResponse(c, nil, nil, err)
		return
	}

	jwtString, err := s.Betting.SignInFromEmail(context.Background(), sd.Email, sd.LinkID)
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

	competition.CreatedByID = s.currentUserID(c)

	data, err := s.Betting.AddCompetition(context.Background(), &competition)

	s.HandleResponse(c, nil, data, err)
}

// DeleteCompetition returns a competition (if it exists).
func (s *Service) DeleteCompetition(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := s.Betting.DeleteCompetition(context.Background(), id)

	s.HandleResponse(c, nil, nil, err)
}

// LockCompetition will lock a competition.
func (s *Service) LockCompetition(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := s.Betting.LockCompetition(context.Background(), id)

	s.HandleResponse(c, nil, nil, err)
}

// SetCompetitionResult will set the result for a competition.
func (s *Service) SetCompetitionResult(c *gin.Context) {
	var result []*pkg.Result

	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := s.Betting.SetCompetitionResult(context.Background(), id, result)

	s.HandleResponse(c, nil, data, err)
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

	in.Competitor.CreatedByID = s.currentUserID(c)

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
	response := map[string]string{
		"jwt": data,
	}

	s.HandleResponse(c, nil, response, err)
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

	bet.BetterID = s.currentUserID(c)

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

func (s *Service) currentUserID(c *gin.Context) int {
	b, ok := c.Get("better")
	if !ok {
		s.Logger.Print("no or invalid authorization header")
		return -1
	}

	better, ok := b.(*pkg.Better)
	if !ok {
		s.Logger.Print("authorization isn't a better")
		return -1
	}

	return better.ID
}

// HandleResponse will respond according to the object and error passed.
func (s *Service) HandleResponse(c *gin.Context, broadcast []byte, response interface{}, err error) {
	if err != nil {
		var httpStatus = http.StatusInternalServerError

		switch errors.Cause(err) {
		case pkg.ErrNotFound:
			httpStatus = http.StatusNotFound
		case pkg.ErrBadRequest:
			httpStatus = http.StatusBadRequest
		default:
			if _, ok := errors.Cause(err).(validation.Errors); ok {
				httpStatus = http.StatusBadRequest
			}
		}

		c.JSON(httpStatus, gin.H{"error": err.Error()})

		return
	}

	if broadcast != nil {
		_ = s.WS.Broadcast(broadcast)
	}

	c.JSON(http.StatusOK, response)
}
