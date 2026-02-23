package main

import (
	"github.com/XeshSufferer/husocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	hub := husocket.New()
	hub.Register("app", app)
	hub.RegisterHandler("close_now", func(client *husocket.Client, message husocket.Message) {
		hub.Close()
	})
	defer hub.Close()

	husocket.ListenWithGracefulShutdownWithReason(app, ":3000", []*husocket.Hub{hub}, "Bye bye!")
}
