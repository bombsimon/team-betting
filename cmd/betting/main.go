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
		r      = gin.Default()
		m      = melody.New()
		logger = log.New(os.Stdout, "TB: ", log.LstdFlags)

		bettingService = &betting.Service{
			DB: database.New(""),
		}

		httpService = bhttp.Service{
			Betting: bettingService,
			WS:      m,
		}
	)

	r.GET("/competition", httpService.GetCompetitions)
	r.POST("/competition", httpService.AddCompetition)
	r.GET("/competition/:id", httpService.GetCompetition)
	r.DELETE("/competition/:id", notImplemented)

	r.GET("/competitor", notImplemented)
	r.POST("/competitor", notImplemented)
	r.GET("/competitor/:id", notImplemented)
	r.DELETE("/competitor/:id", notImplemented)

	r.GET("/better", notImplemented)
	r.POST("/better", notImplemented)
	r.GET("/better/:id", notImplemented)
	r.DELETE("/better/:id", notImplemented)

	r.GET("/bet", notImplemented)
	r.POST("/bet", httpService.AddBet)
	r.GET("/bet/:id", notImplemented)
	r.DELETE("/bet/:id", notImplemented)

	r.GET("/test", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./cmd/betting/index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		if err := m.HandleRequest(c.Writer, c.Request); err != nil {
			logger.Printf("could not handle WS request: %s", err.Error())
		}
	})

	// Just broadcast all messages to everyone when a connection is upgrade to
	// websocket.
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		if err := m.Broadcast(msg); err != nil {
			logger.Printf("could not broadcast ws message")
		}
	})

	if err := r.Run(":5000"); err != nil {
		panic(err)
	}
}

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not implemented"})
}
