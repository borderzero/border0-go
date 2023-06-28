package main

import (
	"log"
	"net/http"
	"os"

	"github.com/borderzero/border0-go"
	"github.com/borderzero/border0-go/listen"
	"github.com/gin-gonic/gin"
)

func main() {
	listener, err := border0.Listen(
		listen.WithSocketName("sdk-socket-http"),              // http socket name the listener will be bound to, socket will be created if not exists
		listen.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
	)
	if err != nil {
		log.Fatalln("failed to start listener:", err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		name := c.Request.Header.Get("X-Auth-Name")
		c.String(http.StatusOK, "Hello, %s! This is border0-go + gin.", name)
	})

	log.Fatalln(r.RunListener(listener))
}
