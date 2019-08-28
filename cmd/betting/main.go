package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bombsimon/team-betting/pkg/betting"
	"github.com/bombsimon/team-betting/pkg/database"
	bhttp "github.com/bombsimon/team-betting/pkg/http"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	var (
		router    = gin.Default()
		wsManager = melody.New()
		logger    = log.New(os.Stdout, "TB: ", log.LstdFlags)

		bettingService = &betting.Service{
			DB: database.New(os.Getenv("DB_DSN")),
		}

		httpService = bhttp.Service{
			Betting: bettingService,
			WS:      wsManager,
		}
	)

	router.GET("/competition", httpService.GetCompetitions)
	router.POST("/competition", httpService.AddCompetition)
	router.GET("/competition/:id", httpService.GetCompetition)
	router.DELETE("/competition/:id", httpService.DeleteCompetition)

	router.GET("/competitor", notImplemented)
	router.POST("/competitor", notImplemented)
	router.GET("/competitor/:id", notImplemented)
	router.DELETE("/competitor/:id", notImplemented)

	router.GET("/better", notImplemented)
	router.POST("/better", notImplemented)
	router.GET("/better/:id", notImplemented)
	router.DELETE("/better/:id", notImplemented)

	router.GET("/bet", notImplemented)
	router.POST("/bet", httpService.AddBet)
	router.GET("/bet/:id", notImplemented)
	router.DELETE("/bet/:id", notImplemented)

	router.GET("/test", func(c *gin.Context) {
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
