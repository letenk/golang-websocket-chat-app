package websocket

const (
	MESSAGE_NEW_USER = "New User"
	MESSAGE_CHAT     = "Chat"
	MESSAGE_LEAVE    = "Leave"
)

type SocketPayload struct {
	Message string
}

type SocketResponse struct {
	From    string
	Type    string
	Message string
}
