package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bombsimon/team-betting/pkg/betting"
	"github.com/bombsimon/team-betting/pkg/database"
	bhttp "github.com/bombsimon/team-betting/pkg/http"
	middleware "github.com/bombsimon/team-betting/pkg/http/middlewares"
	"github.com/bombsimon/team-betting/pkg/mail"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	var (
		router    = gin.Default()
		wsManager = melody.New()
		logger    = log.New(os.Stdout, "TB: ", log.LstdFlags)

		bettingService = &betting.Service{
			DB:          database.New(os.Getenv("DB_DSN")),
			MailService: mail.New(),
		}

		httpService = bhttp.Service{
			Betting: bettingService,
			WS:      wsManager,
			Logger:  logger,
		}
	)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))

	router.POST("/email/send", httpService.SendSignInEmail)
	router.POST("/email/verify", httpService.SignInEmail)
	router.POST("/better", httpService.AddBetter)

	authed := router.Group("/")

	if os.Getenv("DEVELOPMENT") == "" {
		authed.Use(middleware.AuthJWT(bettingService))
	}

	{
		authed.GET("/competition", httpService.GetCompetitions)
		authed.POST("/competition", httpService.AddCompetition)
		authed.GET("/competition/:id", httpService.GetCompetition)
		authed.DELETE("/competition/:id", httpService.DeleteCompetition)
		authed.POST("/competition/:id/lock", httpService.LockCompetition)
		authed.POST("/competition/:id/result", httpService.SetCompetitionResult)

		authed.GET("/competitor", httpService.GetCompetitor)
		authed.POST("/competitor", httpService.AddCompetitor)
		authed.GET("/competitor/:id", httpService.GetCompetitor)
		authed.DELETE("/competitor/:id", httpService.DeleteCompetitor)

		authed.GET("/better", httpService.GetBetters)
		authed.GET("/better/:id", httpService.GetBetter)
		authed.DELETE("/better/:id", httpService.DeleteBetter)

		authed.GET("/bet", httpService.GetBets)
		authed.PUT("/bet", httpService.AddBet)
		authed.GET("/bet/:id", httpService.GetBet)
		authed.DELETE("/bet/:id", httpService.DeleteBet)
	}

	router.Static("/assets", "./cmd/betting/assets")

	router.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./cmd/betting/index.html")
	})

	router.GET("/ws", func(c *gin.Context) {
		if err := wsManager.HandleRequest(c.Writer, c.Request); err != nil {
			logger.Printf("could not handle WS request: %s", err.Error())
		}
	})

	// Just broadcast all messages to everyone when a connection is upgrade to
	// websocket.
	wsManager.HandleMessage(func(s *melody.Session, msg []byte) {
		if err := wsManager.Broadcast(msg); err != nil {
			logger.Printf("could not broadcast ws message")
		}
	})

	if err := router.Run(":5000"); err != nil {
		panic(err)
	}
}

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not implemented"})
}
