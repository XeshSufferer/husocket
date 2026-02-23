package echo

import (
	"net/http"
	"strings"

	"github.com/XeshSufferer/husocket/core"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func UseEchoWS(path string, e *echo.Echo, hub *core.Hub) {

	path = "/" + strings.Trim(path, "/")
	if path == "/" {
		path = "/ws"
	}

	e.GET(path, func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		hub.ServeWS(conn)
		return nil
	})
}
