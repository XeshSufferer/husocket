package husocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Hub struct {
	clients  map[string]*Client
	handlers map[string]func(*Client, Message)
	rooms    map[string]map[string]*Client
	m        sync.RWMutex
	gm       sync.RWMutex

	onConnected    func(*Client, *websocket.Conn)
	onDisconnected func(*Client, *websocket.Conn)
}

type Message struct {
	Method string          `json:"method"`
	Args   json.RawMessage `json:"args"`
}

func New() *Hub {
	return &Hub{
		clients:        make(map[string]*Client),
		handlers:       make(map[string]func(*Client, Message)),
		rooms:          make(map[string]map[string]*Client),
		m:              sync.RWMutex{},
		gm:             sync.RWMutex{},
		onConnected:    func(*Client, *websocket.Conn) {},
		onDisconnected: func(*Client, *websocket.Conn) {},
	}
}

func (h *Hub) OnConnected(f func(*Client, *websocket.Conn)) {
	h.onConnected = f
}

func (h *Hub) OnDisconnected(f func(*Client, *websocket.Conn)) {
	h.onDisconnected = f
}

func (h *Hub) RegisterHandler(method string, handler func(*Client, Message)) {
	h.m.Lock()
	defer h.m.Unlock()
	h.handlers[method] = handler
}

func (h *Hub) Register(path string, app *fiber.App) {
	app.Use(path, func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get(path, websocket.New(func(conn *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)

		client := &Client{Conn: conn, Id: uuid.New(), hub: h, connectedRooms: make([]string, 0), Context: NewContext()}

		h.m.Lock()
		h.clients[client.Id.String()] = client
		h.m.Unlock()

		h.onConnected(client, conn)

		defer func() {
			h.m.Lock()
			delete(h.clients, client.Id.String())
			h.m.Unlock()
			h.onDisconnected(client, conn)

			client.m.Lock()
			for _, room := range client.connectedRooms {
				client.internalRemoveFromRoom(room)
			}
			client.m.Unlock()

			client.Conn.Close()
		}()

		for {

			var message Message

			if mt, msg, err = conn.ReadMessage(); err != nil {
				log.Println("Ошибка чтения:", err)
				break
			}

			if mt != websocket.TextMessage {
				continue
			}

			parseErr := json.Unmarshal(msg, &message)
			if parseErr != nil {
				log.Println(parseErr)
			}

			if handler, ok := h.handlers[message.Method]; ok {
				handler(client, message)
			}
		}
	}))
}

func (h *Hub) RawBroadcast(mt int, msg []byte) {
	h.m.RLock()
	defer h.m.RUnlock()
	for _, client := range h.clients {
		err := client.SendRaw(mt, msg)
		if err != nil {
			log.Println("Error writing to client:", err)
		}
	}
}

func (h *Hub) RawBroadcastToRoom(roomName string, msgType int, msg []byte) {
	h.gm.RLock()
	clients := make([]*Client, 0)
	if r, ok := h.rooms[roomName]; ok {
		for _, c := range r {
			clients = append(clients, c)
		}
	}
	h.gm.RUnlock()

	for _, client := range clients {
		_ = client.SendRaw(msgType, msg)
	}
}

func (h *Hub) BroadcastToRoom(roomName string, method string, msg interface{}) {
	h.gm.RLock()
	clients := make([]*Client, 0)
	if r, ok := h.rooms[roomName]; ok {
		for _, c := range r {
			clients = append(clients, c)
		}
	}
	h.gm.RUnlock()

	for _, client := range clients {
		client.Send(method, msg)
	}
}

func (h *Hub) Broadcast(method string, msg interface{}) {
	h.m.RLock()
	defer h.m.RUnlock()
	for _, client := range h.clients {
		client.Send(method, msg)
	}
}
