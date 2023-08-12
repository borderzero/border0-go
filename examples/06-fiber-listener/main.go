package main

import (
	"fmt"
	"log"
	"os"

	"github.com/borderzero/border0-go"
	"github.com/borderzero/border0-go/listen"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listener, err := border0.Listen(
		listen.WithSocketName("sdk-socket-http"),              // http socket name the listener will be bound to, socket will be created if not exists
		listen.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
	)
	if err != nil {
		log.Fatalln("failed to start listener:", err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		name := c.Get("X-Auth-Name") // border0 will set this header along with a few other identity related headers
		return c.SendString(fmt.Sprintf("Hello, %s! This is border0-go + fiber.", name))
	})

	log.Fatalln(app.Listener(listener))
}
