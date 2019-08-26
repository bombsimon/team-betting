package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/bombsimon/team-betting/pkg"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Service represents the HTTP service serving the team betting.
type Service struct {
	Betting pkg.BettingService
}

// GetCompetitions returns all competitions.
func (s *Service) GetCompetitions(c *gin.Context) {
	comp, err := s.Betting.GetCompetitions(context.Background(), []int{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, comp)
}

// GetCompetition returns a competition (if it exists).
func (s *Service) GetCompetition(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	comp, err := s.Betting.GetCompetition(context.Background(), id)
	if err != nil {
		if errors.Cause(err) == pkg.ErrNotFound {
			c.JSON(http.StatusNotFound, nil)

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, comp)
}

// AddCompetition adds a competition.
func (s *Service) AddCompetition(c *gin.Context) {
	var competition pkg.Competition

	if err := c.ShouldBindJSON(&competition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	comp, err := s.Betting.AddCompetition(context.Background(), &competition)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, comp)
}
