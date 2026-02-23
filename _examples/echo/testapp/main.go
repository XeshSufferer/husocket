package main

import (
	echo_adapter "github.com/XeshSufferer/husocket/adapters/echo"
	core "github.com/XeshSufferer/husocket/core"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	hub := core.New()
	echo_adapter.UseEchoWS("app", app, hub)
	hub.RegisterHandler("ping", func(client *core.Client, message core.Message) {
		client.Send("pong", "pong!")
	})
	app.Start(":3000")
}
