package fiber

import (
	"strings"

	"github.com/XeshSufferer/husocket/core"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func UseFiberWS(path string, app *fiber.App, hub *core.Hub) {

	path = "/" + strings.Trim(path, "/")
	if path == "/" {
		path = "/ws"
	}

	app.Use(path, func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get(path, websocket.New(func(conn *websocket.Conn) {
		hub.SetCloseErrorChecker(func(err error, codes ...core.CloseCode) bool {
			return websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure)
		})
		hub.ServeWS(conn)
	}))
}
