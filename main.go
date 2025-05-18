package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type M map[string]interface{}

const MESSAGE_NEW_USER = "New User"
const MESSAGE_CHAT = "Chat"
const MESSAGE_LEAVE = "Leave"

var connections = make([]*WebSocketConnection, 0)

// SocketPayload as put payload send from frontend
type SocketPayload struct {
	Message string
}

// SocketResponse as
type SocketResponse struct {
	From    string
	Type    string
	Message string
}

// WebSocketConnection as
type WebSocketConnection struct {
	Conn     *websocket.Conn
	Username string
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(w, "Could not open request file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", content)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		currentGorillaConn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		}

		username := r.URL.Query().Get("username")
		currentConn := WebSocketConnection{
			Conn:     currentGorillaConn,
			Username: username,
		}

		connections = append(connections, &currentConn)

		go handleIO(&currentConn, connections)

	})

	port := "8080"
	fmt.Printf("Server starting at :%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}

func handleIO(currentConn *WebSocketConnection, connections []*WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	broadcastMessage(currentConn, MESSAGE_NEW_USER, "")

	for {
		payload := SocketPayload{}
		err := currentConn.Conn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				broadcastMessage(currentConn, MESSAGE_LEAVE, "")
				ejectConnection(currentConn)
				return

			}

			log.Println("ERROR", err.Error())
			continue

		}

		broadcastMessage(currentConn, MESSAGE_CHAT, payload.Message)
	}

}

func ejectConnection(currentConn *WebSocketConnection) {

	var newConnection []*WebSocketConnection
	for _, conn := range connections {
		if conn != currentConn {
			newConnection = append(newConnection, conn)
		}
	}

	connections = newConnection
}

func broadcastMessage(currentConn *WebSocketConnection, kind, message string) {
	for _, eachConn := range connections {
		if eachConn == currentConn {
			continue
		}

		eachConn.Conn.WriteJSON(SocketResponse{
			From:    currentConn.Username,
			Type:    kind,
			Message: message,
		})
	}
}
