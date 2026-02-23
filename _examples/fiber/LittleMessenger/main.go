package main

import (
	"encoding/json"

	fiber_adapter "github.com/XeshSufferer/husocket/adapters/fiber"
	"github.com/XeshSufferer/husocket/core"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	hub := core.New()

	fiber_adapter.UseFiberWS("app", app, hub)
	hub.RegisterHandler("Join", func(client *core.Client, message core.Message) {
		var name string
		err := json.Unmarshal(message.Args, &name)
		if err != nil {
			return
		}

		// Broadcast
		client.Broadcast("OnJoin", message.Args)
		// Set name into locals
		client.Locals.Set("name", name)
	})

	hub.RegisterHandler("SendMessage", func(client *core.Client, message core.Message) {
		var msg string
		err := json.Unmarshal(message.Args, &msg)
		if err != nil {
			return
		}
		client.Broadcast("ReceiveMessage", map[string]string{
			"message": msg,
			"sender":  client.Locals.Get("name").(string),
		})
	})

	hub.OnDisconnected(func(client *core.Client, conn core.WSConnection) {
		name := client.Locals.Get("name").(string)
		client.Broadcast("OnLeave", name)
	})

	app.Listen(":3000")
}
