package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins (adjust in production)
		},
	}
}

func (s *Server) UpgradeConnToWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	log.Println("New WebSocket connection established")
	go s.handleConnection(conn)
}

func (s *Server) handleConnection(conn *websocket.Conn) {
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		HandleMessage(conn, msg)
	}
}
