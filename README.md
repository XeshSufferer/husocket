# husocket

A Go WebSocket server for handling events and managing [Go](https://github.com/XeshSufferer/huclient) and [TypeScript](https://github.com/XeshSufferer/huclient.ts) clients via rooms.

## Description

`husocket` is a lightweight library for building WebSocket servers on top of [Fiber](https://gofiber.io). Features include:

- Message handler registration by method name
- Connected client management
- Client grouping into rooms
- Message broadcasting: to all clients, to a room, or to a specific client
- Context (Locals) for storing per-client data
- Connection/disconnection event handling

## Installation

```bash
go get github.com/XeshSufferer/husocket
```

## Dependencies

- `github.com/gofiber/fiber/v2` — web framework
- `github.com/gofiber/contrib/websocket` — WebSocket middleware for Fiber
- `github.com/google/uuid` — client ID generation

## Quick Start

```go
package main

import (
    "encoding/json"
    "github.com/XeshSufferer/husocket"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/contrib/websocket"
)

func main() {
    app := fiber.New()
    hub := husocket.New()

    // Register WebSocket endpoint
    hub.Register("/ws", app)

    // Handler for "echo" message
    hub.RegisterHandler("echo", func(client *husocket.Client, msg husocket.Message) {
        client.Send("echo_response", msg.Args)
    })

    // Connection event
    hub.OnConnected(func(client *husocket.Client, conn *websocket.Conn) {
        client.Locals.Set("connected_at", time.Now())
    })

    app.Listen(":3000")
}
```

## API

### Hub

| Method | Description |
|--------|-------------|
| `New() *Hub` | Create a new hub instance |
| `Register(path string, app *fiber.App)` | Register WebSocket endpoint in Fiber app |
| `RegisterHandler(method string, handler func(*Client, Message))` | Register a message handler |
| `OnConnected(f func(*Client, *websocket.Conn))` | Callback triggered on client connection |
| `OnDisconnected(f func(*Client, *websocket.Conn))` | Callback triggered on client disconnection |
| `Broadcast(method string, msg interface{})` | Send message to all connected clients |
| `BroadcastToRoom(roomName, method string, msg interface{})` | Send message to a specific room |
| `Close()` / `CloseWithReason(reason string)` | Close all active connections |

### Client

| Method | Description |
|--------|-------------|
| `Send(method string, message interface{})` | Send a message to the client |
| `SendRaw(msgType int, msg []byte)` | Send raw WebSocket data |
| `AddToRoom(roomName string)` / `RemoveFromRoom(roomName string)` | Manage room membership |
| `GetRooms() []string` | Get list of rooms the client belongs to |
| `Broadcast(...)` / `BroadcastToRoom(...)` | Broadcast messages on behalf of the client |
| `Locals *Locals` | Context for storing client-specific data |

### Locals

Thread-safe key-value storage:

```go
client.Locals.Set("user_id", 123)
name := client.Locals.Get("name").(string)
if client.Locals.Exists("token") { ... }
```

## Message Format

Clients and server exchange JSON objects in the following format:

```json
{
  "method": "event_name",
  "args": { ... }
}
```

## Examples

See `_examples/LittleMessenger/` for a working chat application demo.

## License

MIT
