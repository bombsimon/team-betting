package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	r := gin.Default()
	m := melody.New()

	r.GET("/", mockFunc)

	r.POST("/competition", mockFunc)
	r.GET("/competition/:id", mockFunc)

	r.POST("/competition/:id/competitor", mockFunc)
	r.DELETE("/competition/:id/competitor", mockFunc)

	r.POST("/better", mockFunc)

	r.GET("/test", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./cmd/betting/index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	// Just broadcast all messages to everyone when a connection is upgrade to
	// websocket.
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	r.Run(":5000")
}

func mockFunc(c *gin.Context) {
}
