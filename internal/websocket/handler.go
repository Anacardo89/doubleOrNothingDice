package websocket

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Anacardo89/doubleOrNothingDice/internal/user"
	"github.com/gorilla/websocket"
)

// HandleMessage reads a message and dispatches it
func HandleMessage(conn *websocket.Conn, userID string, rawMsg []byte, s *Server) {
	var msg Message
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		log.Println("Invalid message format:", err)
		sendError(conn, "invalid_message", "Could not parse message")
		return
	}

	switch msg.Type {
	case "wallet":
		handleWallet(conn, userID, msg.Payload, s.sessionManager)
	case "play":
		handlePlay(conn, userID, msg.Payload, s.sessionManager)
	case "end_play":
		handleEndPlay(conn, userID, msg.Payload, s.sessionManager)
	case "deposit":
		handleDeposit(conn, userID, msg.Payload, s.sessionManager)
	default:
		sendError(conn, "unknown_type", "Unknown message type")
	}
}

func handleWallet(conn *websocket.Conn, userID string, payload json.RawMessage, sm *user.SessionManager) {
	var req WalletRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(conn, "invalid_wallet_request", "Could not unmarshall wallet request payload")
		return
	}
	session, err := getSession(sm, userID, req)
	if err != nil {
		sendError(conn, "unauthorized_action", err.Error())
		return
	}
	res := WalletResponse{
		Balance: session.Balance,
	}
	sendMessage(conn, "wallet_response", res)
}

func handlePlay(conn *websocket.Conn, userID string, payload json.RawMessage, sm *user.SessionManager) {
	var req PlayRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(conn, "invalid_play_request", "Could not unmarshall play request payload")
		return
	}
	session, err := getSession(sm, userID, req)
	if err != nil {
		sendError(conn, "unauthorized_action", err.Error())
		return
	}
	if session.Game == nil || !session.Game.IsActive {
		_, err := sm.StartGame(userID, req.BetAmount)
		if err != nil {
			sendError(conn, "start_game_error", err.Error())
			return
		}
	}
	playResult, err := sm.PlayRound(userID, req.BetType)
	if err != nil {
		sendError(conn, "play_error", err.Error())
		return
	}
	sendMessage(conn, "play_response", playResult)
}

func handleEndPlay(conn *websocket.Conn, userID string, payload json.RawMessage, sm *user.SessionManager) {
	var req EndPlayRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(conn, "invalid_end_play_request", "Could not unmarshall end play request payload")
		return
	}
	session, err := getSession(sm, userID, req)
	if err != nil {
		sendError(conn, "unauthorized_action", err.Error())
		return
	}
	winnings := 0
	if session.Game != nil {
		winnings = session.Game.CurrentBet
	}
	if err := sm.EndGame(userID); err != nil {
		sendError(conn, "end_game_error", err.Error())
		return
	}
	res := EndPlayResponse{
		Winnings: winnings,
		Balance:  session.Balance,
	}
	sendMessage(conn, "end_play_response", res)
}

func handleDeposit(conn *websocket.Conn, userID string, payload json.RawMessage, sm *user.SessionManager) {
	var req DepositRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		sendError(conn, "invalid_deposit_request", "Could not unmarshall wallet request payload")
		return
	}
	session, err := getSession(sm, userID, req)
	if err != nil {
		sendError(conn, "unauthorized_action", err.Error())
		return
	}
	session.Balance += req.Deposit
	res := DepositResponse{
		Balance: session.Balance,
	}
	sendMessage(conn, "deposit_response", res)
}

func getSession(sm *user.SessionManager, userID string, req ReqWithClient) (*user.Session, error) {
	clientID := req.GetClientID()
	if clientID != userID {
		return nil, errors.New("ClientID in request does not match authenticated user.")
	}
	session, exists := sm.Get(req.GetClientID())
	if !exists {
		session, err := sm.Create(req.GetClientID())
		if err != nil {
			return nil, err
		}
		return session, nil
	}
	return session, nil
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
