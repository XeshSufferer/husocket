package husocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type Client struct {
	Conn           *websocket.Conn
	Id             uuid.UUID
	Context        *Context
	hub            *Hub
	connectedRooms []string
	m              sync.Mutex
}

func (c *Client) Send(method string, message interface{}) {
	args, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	msg := Message{method, args}
	json, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	c.SendRaw(websocket.TextMessage, json)
}

func (c *Client) SendRaw(msgType int, msg []byte) error {
	c.m.Lock()
	defer c.m.Unlock()
	writeErr := c.Conn.WriteMessage(msgType, msg)

	if writeErr != nil {
		return writeErr
	}
	return nil
}

func (c *Client) AddToRoom(roomName string) {
	c.hub.gm.Lock()
	if r, ok := c.hub.rooms[roomName]; ok {
		r[c.Id.String()] = c
	} else {
		c.hub.rooms[roomName] = map[string]*Client{c.Id.String(): c}
	}
	c.hub.gm.Unlock()

	c.m.Lock()
	c.connectedRooms = append(c.connectedRooms, roomName)
	c.m.Unlock()
}

func (c *Client) RemoveFromRoom(roomName string) {
	c.internalRemoveFromRoom(roomName)

	c.m.Lock()
	defer c.m.Unlock()

	newRooms := make([]string, 0, len(c.connectedRooms))
	for _, room := range c.connectedRooms {
		if room != roomName {
			newRooms = append(newRooms, room)
		}
	}
	c.connectedRooms = newRooms
}

func (c *Client) internalRemoveFromRoom(roomName string) {
	c.hub.gm.Lock()
	defer c.hub.gm.Unlock()
	if r, ok := c.hub.rooms[roomName]; ok {
		delete(r, c.Id.String())
	}
}

func (c *Client) GetRooms() []string {
	c.m.Lock()
	defer c.m.Unlock()

	rooms := make([]string, len(c.connectedRooms))
	copy(rooms, c.connectedRooms)
	return rooms
}

func (c *Client) RawBroadcast(mt int, msg []byte) {
	c.hub.m.RLock()
	defer c.hub.m.RUnlock()
	for _, client := range c.hub.clients {
		err := client.SendRaw(mt, msg)
		if err != nil {
			log.Println("Error writing to client:", err)
		}
	}
}

func (c *Client) RawBroadcastToRoom(roomName string, msgType int, msg []byte) {
	c.hub.gm.RLock()
	clients := make([]*Client, 0)
	if r, ok := c.hub.rooms[roomName]; ok {
		for _, c := range r {
			clients = append(clients, c)
		}
	}
	c.hub.gm.RUnlock()

	for _, client := range clients {
		_ = client.SendRaw(msgType, msg)
	}
}

func (c *Client) BroadcastToRoom(roomName string, method string, msg interface{}) {
	c.hub.gm.RLock()
	clients := make([]*Client, 0)
	if r, ok := c.hub.rooms[roomName]; ok {
		for _, c := range r {
			clients = append(clients, c)
		}
	}
	c.hub.gm.RUnlock()

	for _, client := range clients {
		client.Send(method, msg)
	}
}

func (c *Client) Broadcast(method string, msg interface{}) {
	c.hub.m.RLock()
	defer c.hub.m.RUnlock()
	for _, client := range c.hub.clients {
		client.Send(method, msg)
	}
}
