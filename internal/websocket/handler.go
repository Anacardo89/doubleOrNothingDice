package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// HandleMessage reads a message and dispatches it
func HandleMessage(conn *websocket.Conn, rawMsg []byte) {
	var msg Message
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		log.Println("Invalid message format:", err)
		sendError(conn, "invalid_message", "Could not parse message")
		return
	}

	switch msg.Type {
	case "wallet":
		handleWallet(conn, msg.Payload)

	case "play":
		handlePlay(conn, msg.Payload)

	case "end_play":
		handleEndPlay(conn, msg.Payload)

	default:
		sendError(conn, "unknown_type", "Unknown message type")
	}
}

func handleWallet(conn *websocket.Conn, payload json.RawMessage) {
	// Placeholder response
	res := WalletResponse{
		Balance: 100,
	}
	conn.WriteJSON(res)
}

func handlePlay(conn *websocket.Conn, payload json.RawMessage) {
	// Placeholder response
	res := PlayResponse{
		RolledNumber: 4,
		Result:       "win",
	}
	conn.WriteJSON(res)
}

func handleEndPlay(conn *websocket.Conn, payload json.RawMessage) {
	// Placeholder response
	res := EndPlayResponse{
		Winnings: 50,
	}
	conn.WriteJSON(res)
}

func sendError(conn *websocket.Conn, errType, message string) {
	conn.WriteJSON(map[string]interface{}{
		"type":    errType,
		"message": message,
	})
}
