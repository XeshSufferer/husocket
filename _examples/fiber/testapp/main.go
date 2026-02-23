package main

import (
	fiber_adapter "github.com/XeshSufferer/husocket/adapters/fiber"
	"github.com/XeshSufferer/husocket/core"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	hub := core.New()
	fiber_adapter.UseFiberWS("app", app, hub)
	hub.RegisterHandler("ping", func(client *core.Client, message core.Message) {
		client.Send("pong", "pong!")
	})
	app.Listen(":3000")
}
