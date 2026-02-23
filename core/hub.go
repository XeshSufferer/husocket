package core

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Hub struct {
	clients  map[string]*Client
	handlers map[string]func(*Client, Message)
	rooms    map[string]map[string]*Client
	m        sync.RWMutex
	gm       sync.RWMutex

	onConnected    func(*Client, WSConnection)
	onDisconnected func(*Client, WSConnection)
	isCloseError   func(error, ...CloseCode) bool
}

type WSConnection interface {
	WriteMessage(mt int, data []byte) error
	ReadMessage() (int, []byte, error)
	Close() error
	WriteControl(mt int, data []byte, deadline time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

type Message struct {
	Method string          `json:"method"`
	Args   json.RawMessage `json:"args"`
}

func New() *Hub {
	h := &Hub{
		clients:        make(map[string]*Client),
		handlers:       make(map[string]func(*Client, Message)),
		rooms:          make(map[string]map[string]*Client),
		m:              sync.RWMutex{},
		gm:             sync.RWMutex{},
		onConnected:    func(*Client, WSConnection) {},
		onDisconnected: func(*Client, WSConnection) {},
	}

	log.SetPrefix("[husocket] ")

	h.RegisterHandler("close", func(client *Client, message Message) {
		log.Println("close client")
		client.Conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		client.Conn.Close()
	})
	return h
}

func (h *Hub) OnConnected(f func(*Client, WSConnection)) {
	h.onConnected = f
}

func (h *Hub) OnDisconnected(f func(*Client, WSConnection)) {
	h.onDisconnected = f
}

func (h *Hub) ServeWS(conn WSConnection) {
	client := &Client{
		Conn:           conn,
		Id:             uuid.New(),
		hub:            h,
		connectedRooms: make([]string, 0),
		Locals:         NewContext(),
	}

	h.m.Lock()
	h.clients[client.Id.String()] = client
	h.m.Unlock()

	h.onConnected(client, conn)
	defer h.disconnect_internal(client)

	for {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		mt, msg, readErr := conn.ReadMessage()

		if readErr != nil {
			if h.isCloseError != nil && !h.isCloseError(readErr, CloseGoingAway, CloseAbnormalClosure) {
				log.Printf("Client disconnected: %v", readErr)
			}
			return
		}

		if mt != int(MessageText) {
			continue
		}

		var message Message
		if parseErr := json.Unmarshal(msg, &message); parseErr != nil {
			log.Printf("JSON parse error: %v", parseErr)
			continue
		}

		log.Printf("Received: method=%s", message.Method)

		if handler, ok := h.handlers[message.Method]; ok {
			handler(client, message)
		}
	}
}

func (h *Hub) RegisterHandler(method string, handler func(*Client, Message)) {
	h.m.Lock()
	defer h.m.Unlock()
	h.handlers[method] = handler
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

func (h *Hub) SetCloseErrorChecker(f func(error, ...CloseCode) bool) {
	h.isCloseError = f
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

func (h *Hub) Close() {
	h.shutdown_internal("closed by server")
}

func (h *Hub) CloseWithReason(reason string) {
	h.shutdown_internal(reason)
}

func (h *Hub) shutdown_internal(reason string) {
	h.m.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for _, c := range h.clients {
		clients = append(clients, c)
	}
	h.m.RUnlock()

	for _, client := range clients {
		client.Send("close", reason)
	}

	time.Sleep(100 * time.Millisecond)

	for _, client := range clients {
		h.disconnect_internal(client)
	}
}

func (h *Hub) disconnect_internal(client *Client) {
	h.m.Lock()
	delete(h.clients, client.Id.String())
	h.m.Unlock()
	h.onDisconnected(client, client.Conn)

	client.m.Lock()
	for _, room := range client.connectedRooms {
		client.internalRemoveFromRoom(room)
	}
	client.m.Unlock()

	_ = client.Conn.WriteControl(
		int(CloseGoingAway),
		[]byte("server shutting down"),
		time.Now().Add(time.Second*5),
	)

	time.Sleep(100 * time.Millisecond)
	client.Conn.Close()
}
