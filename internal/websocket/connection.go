package websocket

import (
	"log"
	"net/http"

	"github.com/Anacardo89/doubleOrNothingDice/internal/auth"
	"github.com/Anacardo89/doubleOrNothingDice/internal/db"
	"github.com/Anacardo89/doubleOrNothingDice/internal/redis"
	"github.com/Anacardo89/doubleOrNothingDice/internal/user"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader          websocket.Upgrader
	connectionManager *ConnectionManager
	sessionManager    *user.SessionManager
}

func NewServer(dbManager *db.Manager, redisManager *redis.Manager) *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		connectionManager: NewConnectionManager(5),
		sessionManager:    user.NewSessionManager(dbManager, redisManager),
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
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	log.Printf("New WebSocket connection established for user %s\n", userID)
	s.connectionManager.Add(userID, conn)
	go s.handleConnection(conn, userID)
}

func (s *Server) handleConnection(conn *websocket.Conn, userID string) {
	defer func() {
		s.connectionManager.Remove(userID)
		conn.Close()
	}()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		s.connectionManager.UpdateActivity(userID)
		HandleMessage(conn, userID, msg, s)
	}
}
