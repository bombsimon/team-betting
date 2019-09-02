package middleware

import (
	"context"
	"strings"

	"github.com/bombsimon/team-betting/pkg"
	"github.com/gin-gonic/gin"
)

// AuthJWT will ensure routes which requires it authenticated for.
func AuthJWT(s pkg.BettingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.Next()
			return
		}

		authorizationParts := strings.Split(auth, " ")
		if len(authorizationParts) != 2 {
			c.Next()
			return
		}

		better, err := s.BetterFromJWT(context.Background(), authorizationParts[1])
		if err != nil {
			c.Next()
			return
		}

		c.Set("better", better)

		c.Next()
	}
}
