package stdhttp

import (
	"net/http"
	"strings"

	"github.com/XeshSufferer/husocket/core"
	"github.com/gorilla/websocket"
)

var DefaultUpgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func UseStdHTTPWS(path string, mux *http.ServeMux, hub *core.Hub) {
	UseStdHTTPWSWithUpgrader(path, mux, hub, DefaultUpgrader)
}

func UseStdHTTPWSWithUpgrader(path string, mux *http.ServeMux, hub *core.Hub, upgrader *websocket.Upgrader) {

	path = "/" + strings.Trim(path, "/")
	if path == "/" {
		path = "/ws"
	}

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
			return
		}
		hub.ServeWS(conn)
	})
}
