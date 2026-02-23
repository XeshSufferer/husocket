package main

import (
	"net/http"

	"github.com/XeshSufferer/husocket/adapters/stdhttp"

	"github.com/XeshSufferer/husocket/core"
)

func main() {
	mux := http.NewServeMux()
	hub := core.New()
	stdhttp.UseStdHTTPWS("/app", mux, hub)
	hub.RegisterHandler("ping", func(client *core.Client, message core.Message) {
		client.Send("pong", "pong!")
	})
	http.ListenAndServe(":3000", mux)
}
