package main

import (
	"encoding/json"

	"github.com/XeshSufferer/husocket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	hub := husocket.New()

	hub.Register("app", app)

	hub.RegisterHandler("Join", func(client *husocket.Client, message husocket.Message) {
		var name string
		err := json.Unmarshal(message.Args, &name)
		if err != nil {
			return
		}

		// Broadcast
		client.Broadcast("OnJoin", message.Args)
		// Set name into context
		client.Context.Set("name", name)
	})

	hub.RegisterHandler("SendMessage", func(client *husocket.Client, message husocket.Message) {
		var msg string
		err := json.Unmarshal(message.Args, &msg)
		if err != nil {
			return
		}
		client.Broadcast("ReceiveMessage", map[string]string{
			"message": msg,
			"sender":  client.Context.Get("name").(string),
		})
	})

	hub.OnDisconnected(func(client *husocket.Client, conn *websocket.Conn) {
		name := client.Context.Get("name").(string)
		client.Broadcast("OnLeave", name)
	})

	app.Listen(":3000")
}
