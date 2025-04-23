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
	payloadBytes, err := json.Marshal(res)
	if err != nil {
		log.Println("Error marshaling WalletResponse:", err)
		return
	}
	msg := Message{
		Type:    "wallet_response",
		Payload: payloadBytes,
	}
	if err := conn.WriteJSON(msg); err != nil {
		log.Println("Error sending message:", err)
	}
}

func handlePlay(conn *websocket.Conn, payload json.RawMessage) {
	// Placeholder response
	res := PlayResponse{
		RolledNumber: 4,
		NextBet:      40,
		Outcome:      "win",
	}
	payloadBytes, err := json.Marshal(res)
	if err != nil {
		log.Println("Error marshaling PlayResponse:", err)
		return
	}
	msg := Message{
		Type:    "play_response",
		Payload: payloadBytes,
	}
	if err := conn.WriteJSON(msg); err != nil {
		log.Println("Error sending message:", err)
	}
}

func handleEndPlay(conn *websocket.Conn, payload json.RawMessage) {
	// Placeholder response
	res := EndPlayResponse{
		Winnings: 40,
		Balance:  120,
	}
	payloadBytes, err := json.Marshal(res)
	if err != nil {
		log.Println("Error marshaling EndPlayResponse:", err)
		return
	}
	msg := Message{
		Type:    "end_play_response",
		Payload: payloadBytes,
	}
	if err := conn.WriteJSON(msg); err != nil {
		log.Println("Error sending message:", err)
	}
}

func sendError(conn *websocket.Conn, errType, message string) {
	conn.WriteJSON(map[string]interface{}{
		"type":    errType,
		"message": message,
	})
}
