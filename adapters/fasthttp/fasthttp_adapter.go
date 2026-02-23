package fasthttp

import (
	"log"
	"strings"

	"github.com/XeshSufferer/husocket/core"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"time"
)

type wsConn struct {
	*websocket.Conn
}

func (c *wsConn) WriteMessage(mt int, data []byte) error {
	return c.Conn.WriteMessage(mt, data)
}

func (c *wsConn) ReadMessage() (int, []byte, error) {
	return c.Conn.ReadMessage()
}

func (c *wsConn) Close() error {
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	_ = c.Conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))
	time.Sleep(50 * time.Millisecond)
	return c.Conn.Close()
}

func (c *wsConn) WriteControl(mt int, data []byte, deadline time.Time) error {
	if mt == websocket.CloseMessage {
		formatted := websocket.FormatCloseMessage(websocket.CloseNormalClosure, string(data))
		return c.Conn.WriteControl(websocket.CloseMessage, formatted, deadline)
	}
	return c.Conn.WriteControl(mt, data, deadline)
}

func (c *wsConn) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *wsConn) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}

func UseFastHTTPWS(path string, srv *fasthttp.Server, hub *core.Hub) {
	UseFastHTTPWSWithUpgrader(path, srv, hub, websocket.FastHTTPUpgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
}

func UseFastHTTPWSWithUpgrader(path string, srv *fasthttp.Server, hub *core.Hub, upgrader websocket.FastHTTPUpgrader) {
	oldHandler := srv.Handler

	registerPath := "/" + strings.Trim(path, "/")
	if registerPath == "/" {
		registerPath = "/ws" // дефолт, если пусто
	}

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		requestPath := string(ctx.Path())
		if requestPath == "" {
			requestPath = "/"
		}
		requestPath = strings.TrimRight(requestPath, "/")

		if requestPath != registerPath {
			if oldHandler != nil {
				oldHandler(ctx)
			} else {
				ctx.Error("Not found", fasthttp.StatusNotFound)
			}
			return
		}

		if !websocket.FastHTTPIsWebSocketUpgrade(ctx) {
			ctx.Error("WebSocket upgrade required", fasthttp.StatusBadRequest)
			return
		}

		err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
			defer func() {
				_ = recover()
			}()
			hub.ServeWS(&wsConn{conn})
		})
		if err != nil {
			log.Printf("[fasthttp] WebSocket upgrade error: %v", err)
		}
	}
}
