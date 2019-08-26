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
		r              = gin.Default()
		m              = melody.New()
		logger         = log.New(os.Stdout, "TB: ", log.LstdFlags)
		bettingService = &betting.Service{
			DB: database.New(""),
		}
		httpService = bhttp.Service{
			Betting: bettingService,
		}
	)

	r.GET("/competition", httpService.GetCompetitions)
	r.POST("/competition", httpService.AddCompetition)
	r.GET("/competition/:id", httpService.GetCompetition)

	r.POST("/competition/:id/competitor", notImplemented)
	r.DELETE("/competition/:id/competitor", notImplemented)

	r.POST("/better", notImplemented)

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
