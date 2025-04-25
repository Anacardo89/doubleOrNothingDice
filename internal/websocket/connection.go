package websocket

import (
	"log"
	"net/http"

	"github.com/Anacardo89/doubleOrNothingDice/internal/auth"
	"github.com/Anacardo89/doubleOrNothingDice/internal/user"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader       websocket.Upgrader
	sessionManager *user.SessionManager
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins (adjust in production)
		},
		sessionManager: user.NewSessionManager(),
	}
}

func (s *Server) UpgradeConnToWS(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	claims, err := auth.ParseToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	userID := claims.ClientID
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	log.Println("New WebSocket connection established")
	go s.handleConnection(conn, userID)
}

func (s *Server) handleConnection(conn *websocket.Conn, userID string) {
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		HandleMessage(conn, userID, msg, s)
	}
}
