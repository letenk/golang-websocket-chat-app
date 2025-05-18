package websocket

import "log"

// broadcastMessage sends a message to all connected users except the sender
func broadcastMessage(senderConn *WebSocketConnection, msgType string, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn := range connections {
		if conn != senderConn {
			conn.Send <- SocketResponse{
				From:    senderConn.Username,
				Type:    msgType,
				Message: message,
			}
		}
	}
}

// ejectConnection as removes connections that have left the chat from the connection list and cleans up resource
func ejectConnection(currentConn *WebSocketConnection) {

	mutex.Lock()
	delete(connections, currentConn)
	mutex.Unlock()

	// close channel to stop sender goroutine
	close(currentConn.Send)
	currentConn.Conn.Close()

	log.Printf("User %s disconnected\n", currentConn.Username)

}
