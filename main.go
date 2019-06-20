package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
)

var server *socketio.Server

func socketHandler(c *gin.Context) {

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.ServeHTTP(c.Writer, c.Request)
}

func main() {
	router := gin.Default()
	var err error

	router.Static("/static", "./static")
	router.LoadHTMLGlob("view/*")

	server, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	router.GET("/socket.io/", socketHandler)
	router.POST("/socket.io/", socketHandler)
	router.Handle("WS", "/socket.io/", socketHandler)
	router.Handle("WSS", "/socket.io/", socketHandler)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	router.Run(":8080")
}
