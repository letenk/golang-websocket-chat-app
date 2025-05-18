package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/letenk/golang-websocket-chat-app/websocket"
)

func main() {

	// Route handler for the main chat page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Could not open request file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", content)
	})

	// Handle WebSocket connections
	http.HandleFunc("/ws", websocket.HandleWebSocket)

	// Start the HTTP server
	port := "8080"
	fmt.Printf("Server starting at :%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
