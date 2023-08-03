package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/borderzero/border0-go"
	"github.com/borderzero/border0-go/listen"
	"github.com/labstack/echo/v4"
)

func main() {
	listener, err := border0.Listen(
		listen.WithSocketName("sdk-socket-http"),              // http socket name the listener will be bound to, socket will be created if not exists
		listen.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
	)
	if err != nil {
		log.Fatalln("failed to start listener:", err)
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		name := c.Request().Header.Get("X-Auth-Name")
		return c.String(http.StatusOK, fmt.Sprintf("Hello, %s! This is border0-go + echo.", name))
	})

	e.Listener = listener
	e.Logger.Fatal(e.Start(""))
}
