package websocket

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to a WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	currentGorillaConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	// Get the username from the query parameter
	username := r.URL.Query().Get("username")
	if username == "" {
		username = fmt.Sprintf("user-anonym-%d", len(connections)+1)
	}

	// Create a new WebSocketConnection instance from coming connection
	wsConn := &WebSocketConnection{
		Conn:     currentGorillaConn,
		Username: username,
		Send:     make(chan SocketResponse),
	}

	// Save connection safely
	mutex.Lock()
	connections[wsConn] = true
	mutex.Unlock()

	// Notify others that user has joined
	broadcastMessage(wsConn, MESSAGE_NEW_USER, "")

	// Start goroutines to handle read and write
	go handleReceive(wsConn)
	go handleSend(wsConn)
}

// handleIO handles incoming messages from a specific client connectionn
func handleReceive(currentConn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("RECOVER ERROR:", r)
		}
		// Ensure cleanup
		ejectConnection(currentConn)
		currentConn.Conn.Close()
	}()

	// Continuously read messages from the current connectiont
	for {
		payload := SocketPayload{}

		// Read incoming JSON message
		err := currentConn.Conn.ReadJSON(&payload)
		if err != nil {

			// If the connection is closed by the client
			if strings.Contains(err.Error(), "websocket: close") {
				// Notify others that the user has left
				broadcastMessage(currentConn, MESSAGE_LEAVE, "")
				return
			}

			log.Println("ERROR", err.Error())
			continue
		}

		// Broadcast the received message to other users
		broadcastMessage(currentConn, MESSAGE_CHAT, payload.Message)
	}

}

// handleSend sends messages from the Send channel to the WebSocket connection
func handleSend(currentConn *WebSocketConnection) {
	for msg := range currentConn.Send {
		err := currentConn.Conn.WriteJSON(msg)
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}
