package main

import (
	fasthttp_adapter "github.com/XeshSufferer/husocket/adapters/fasthttp"
	"github.com/XeshSufferer/husocket/core"
	"github.com/valyala/fasthttp"
)

func main() {
	app := &fasthttp.Server{}
	hub := core.New()
	fasthttp_adapter.UseFastHTTPWS("/ws", app, hub)
	hub.RegisterHandler("ping", func(client *core.Client, message core.Message) {
		client.Send("pong", "pong!")
	})
	app.ListenAndServe(":3000")
}
