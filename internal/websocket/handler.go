package websocket

import (
	"encoding/json"
	"log"

	"github.com/Anacardo89/doubleOrNothingDice/internal/game"
	"github.com/Anacardo89/doubleOrNothingDice/internal/user"
	"github.com/gorilla/websocket"
)

// HandleMessage reads a message and dispatches it
func HandleMessage(conn *websocket.Conn, rawMsg []byte, s *Server) {
	var msg Message
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		log.Println("Invalid message format:", err)
		sendError(conn, "invalid_message", "Could not parse message")
		return
	}

	switch msg.Type {
	case "wallet":
		handleWallet(conn, msg.Payload, s.sessionManager)

	case "play":
		handlePlay(conn, msg.Payload, s.sessionManager)

	case "end_play":
		handleEndPlay(conn, msg.Payload, s.sessionManager)

	default:
		sendError(conn, "unknown_type", "Unknown message type")
	}
}

func handleWallet(conn *websocket.Conn, payload json.RawMessage, sm *user.SessionManager) {
	var req WalletRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(conn, "invalid_wallet_request", "Could not unmarshall wallet request payload")
		return
	}
	session := getSession(sm, req)
	res := WalletResponse{
		Balance: session.Balance,
	}
	sendMessage(conn, "wallet_response", res)
}

func handlePlay(conn *websocket.Conn, payload json.RawMessage, sm *user.SessionManager) {
	var req PlayRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(conn, "invalid_play_request", "Could not unmarshall play request payload")
		return
	}
	session := getSession(sm, req)
	if session.Game == nil || !session.Game.IsActive {
		session.Game = game.NewGame(req.ClientID, req.BetAmount)
		session.Balance -= req.BetAmount
	}
	playResult, err := session.Game.Play(req.BetType)
	if err != nil {
		sendError(conn, "play_error", err.Error())
		return
	}
	sendMessage(conn, "play_response", playResult)
}

func handleEndPlay(conn *websocket.Conn, payload json.RawMessage, sm *user.SessionManager) {
	var req EndPlayRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(conn, "invalid_end_play_request", "Could not unmarshall end play request payload")
		return
	}
	session := getSession(sm, req)
	if session.Game == nil || !session.Game.IsActive {
		sendError(conn, "no_active_game", "No active game to end")
		return
	}
	session.Game.EndGame()
	session.Balance += session.Game.CurrentBet
	res := EndPlayResponse{
		Winnings: session.Game.CurrentBet,
		Balance:  session.Balance,
	}
	session.Game = nil
	sendMessage(conn, "end_play_response", res)
}

func getSession(sm *user.SessionManager, req ReqWithClient) *user.Session {
	session, exists := sm.Get(req.GetClientID())
	if !exists {
		session = sm.Create(req.GetClientID())
	}
	return session
}

func sendMessage(conn *websocket.Conn, msgType string, payload interface{}) {
	conn.WriteJSON(map[string]any{
		"type":    msgType,
		"payload": payload,
	})
}

func sendError(conn *websocket.Conn, errType, message string) {
	conn.WriteJSON(map[string]any{
		"type":    errType,
		"message": message,
	})
}
