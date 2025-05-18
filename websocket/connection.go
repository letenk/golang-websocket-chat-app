package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketConnection represents a single connection and related data
// Send channel is used to send message safely in background
type WebSocketConnection struct {
	Conn     *websocket.Conn
	Username string
	Send     chan SocketResponse
}

// Store active connections with thread-safe access
var (
	connections = make(map[*WebSocketConnection]bool)
	mutex       sync.Mutex
)
